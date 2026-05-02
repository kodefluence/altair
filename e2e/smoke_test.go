//go:build e2e

package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kodefluence/altair/e2e/harness"
)

// binPath is built once in TestMain and reused across every subtest.
var binPath string

func TestMain(m *testing.M) {
	tmp, err := os.MkdirTemp("", "altair-e2e-bin-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "mktmp for binary: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmp)

	binPath = filepath.Join(tmp, "altair-e2e")

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "getwd: %v\n", err)
		os.Exit(1)
	}
	repoRoot := filepath.Dir(cwd) // e2e/ -> repo root

	build := exec.Command("go", "build", "-o", binPath, ".")
	build.Dir = repoRoot
	if out, err := build.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "go build failed: %v\n%s\n", err, out)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

// startGateway starts a fresh altair process for one test. Each test gets
// its own subprocess so config-level differences (timeout, body cap) don't
// require a single bloated harness.
func startGateway(t *testing.T, cfg harness.Config) *harness.Harness {
	t.Helper()
	h := harness.Start(t, binPath, cfg)
	t.Cleanup(h.Stop)
	return h
}

// T1 — forwarding_unauthed: a plain GET round-trips through the gateway
// to the mock upstream, with the gateway-injected request id propagated.
func TestSmoke_ForwardingUnauthed(t *testing.T) {
	h := startGateway(t, harness.Config{
		ProxyHost:          "altair.test",
		UpstreamTimeout:    "10s",
		MaxRequestBodySize: "1MB",
	})

	resp, err := http.Get(h.URL("/users/profiles/123"))
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	body := harness.DrainBody(t, resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: want 200, got %d. body=%s", resp.StatusCode, body)
	}

	var echo harness.Echo
	if err := json.Unmarshal(body, &echo); err != nil {
		t.Fatalf("decode echo: %v\nbody=%s", err, body)
	}
	if echo.Method != "GET" {
		t.Errorf("upstream method: want GET, got %q", echo.Method)
	}
	if echo.Path != "/users/profiles/123" {
		t.Errorf("upstream path: want /users/profiles/123, got %q", echo.Path)
	}
	if got := echo.Headers["X-Request-Id"]; len(got) == 0 || got[0] == "" {
		t.Errorf("upstream missing X-Request-Id header; got: %v", echo.Headers)
	}
}

// T2 — proxy_host: the configured proxy.host arrives at the upstream as
// the Host header, regardless of the client's Host. Captured once at
// gateway boot so this is stable for the life of the process.
func TestSmoke_ProxyHostInjected(t *testing.T) {
	h := startGateway(t, harness.Config{
		ProxyHost:          "captured.example.com",
		UpstreamTimeout:    "10s",
		MaxRequestBodySize: "1MB",
	})

	resp, err := http.Get(h.URL("/users/profiles/abc"))
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	body := harness.DrainBody(t, resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: want 200, got %d. body=%s", resp.StatusCode, body)
	}

	var echo harness.Echo
	if err := json.Unmarshal(body, &echo); err != nil {
		t.Fatalf("decode echo: %v", err)
	}
	if echo.Host != "captured.example.com" {
		t.Errorf("upstream Host: want captured.example.com, got %q", echo.Host)
	}
}

// T3 — body_size_cap: a POST whose body exceeds proxy.max_request_body_size
// is rejected with 413 before the upstream is dialed.
func TestSmoke_BodySizeCapRejects(t *testing.T) {
	h := startGateway(t, harness.Config{
		ProxyHost:          "altair.test",
		UpstreamTimeout:    "10s",
		MaxRequestBodySize: "16B", // 16 bytes
	})

	// 64 bytes — well over the 16-byte cap.
	body := bytes.Repeat([]byte("a"), 64)
	resp, err := http.Post(h.URL("/users/echo"), "application/octet-stream", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("status: want 413, got %d", resp.StatusCode)
	}
}

// T4 — upstream_timeout: with a tight timeout against a slow upstream the
// gateway short-circuits with 502 well before the upstream finishes.
func TestSmoke_UpstreamTimeoutFires(t *testing.T) {
	h := startGateway(t, harness.Config{
		ProxyHost:          "altair.test",
		UpstreamTimeout:    "300ms",
		MaxRequestBodySize: "1MB",
	})

	start := time.Now()
	resp, err := http.Get(h.URL("/users/slow-3s/abc"))
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	elapsed := time.Since(start)
	resp.Body.Close()

	if resp.StatusCode != http.StatusBadGateway {
		t.Fatalf("status: want 502, got %d", resp.StatusCode)
	}
	if elapsed > 2*time.Second {
		t.Fatalf("expected fast 502 (<2s), but took %s — timeout not honored", elapsed)
	}
}

// Sanity: when the body cap is unset (zero) and a generous timeout is
// configured, oversized POSTs flow through. Pins the "default unlimited"
// promise from cfg/app.go.
func TestSmoke_NoBodyCapDefaultsUnlimited(t *testing.T) {
	h := startGateway(t, harness.Config{
		ProxyHost:          "altair.test",
		UpstreamTimeout:    "10s",
		MaxRequestBodySize: "0", // 0 = unlimited
	})

	body := strings.Repeat("a", 4096)
	resp, err := http.Post(h.URL("/users/echo"), "text/plain", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	got := harness.DrainBody(t, resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: want 200, got %d. body=%s", resp.StatusCode, got)
	}

	var echo harness.Echo
	if err := json.Unmarshal(got, &echo); err != nil {
		t.Fatalf("decode echo: %v", err)
	}
	if len(echo.Body) != len(body) {
		t.Errorf("echo body length: want %d, got %d", len(body), len(echo.Body))
	}
}
