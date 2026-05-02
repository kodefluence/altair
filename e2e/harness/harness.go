//go:build e2e

package harness

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

// Harness owns one altair subprocess driven by config materialized in a
// per-test tmpdir. Each subtest receives the same harness; isolation comes
// from request semantics (different paths, headers), not infra rebuild.
type Harness struct {
	t        *testing.T
	BinPath  string
	WorkDir  string
	Port     int
	Upstream *Upstream
	cmd      *exec.Cmd
}

// Config tweaks the materialized app.yml so subtests can stress-test the
// proxy options that landed in this branch.
type Config struct {
	ProxyHost          string // proxy.host header sent upstream
	UpstreamTimeout    string // proxy.upstream_timeout (Go duration string)
	MaxRequestBodySize string // proxy.max_request_body_size ("100B", "10KB", ...)

	// EnableOauth keeps the oauth plugin active and runs altair with
	// --auto-migrate so migrations apply on boot. Requires MySQL on
	// 127.0.0.1:3306 with credentials matching MySQLDSN.
	EnableOauth bool
}

// Start scaffolds via the binary's `new` subcommand, patches the generated
// config to point at the in-process mock upstream + Phase-A-only feature
// matrix, and waits until /healthcheck answers 200.
func Start(t *testing.T, binPath string, cfg Config) *Harness {
	t.Helper()

	port, err := FreeTCP()
	if err != nil {
		t.Fatalf("allocate altair port: %v", err)
	}

	workDir := t.TempDir()
	upstream := NewUpstream()

	// `altair new .` produces a config tree, .env, routes/, etc. Run it in
	// the tmpdir.
	scaffold := exec.Command(binPath, "new", ".")
	scaffold.Dir = workDir
	if out, err := scaffold.CombinedOutput(); err != nil {
		upstream.Close()
		t.Fatalf("altair new failed: %v\n%s", err, out)
	}

	// Phase A doesn't run oauth (no MySQL). Drop oauth from the activation
	// list and rewrite the route file to keep only the auth=none path.
	if err := patchAppYAML(workDir, port, cfg); err != nil {
		upstream.Close()
		t.Fatalf("patch app.yml: %v", err)
	}
	if err := patchRoutes(workDir, upstream.HostPort(), cfg.EnableOauth); err != nil {
		upstream.Close()
		t.Fatalf("patch routes: %v", err)
	}

	// Inject env vars the templated config consumes. We don't rely on
	// env.sample because the smoke test must run from a clean checkout
	// without manual env setup. Database vars are always set so even the
	// gateway-only path can pass them through harmlessly; the oauth path
	// actually uses them via --auto-migrate.
	env := append(os.Environ(),
		"BASIC_AUTH_PASSWORD=smoke",
		"PROXY_HOST="+cfg.ProxyHost,
		"EXAMPLE_USERS_SERVICE_HOST="+upstream.HostPort(),
		"DATABASE_HOST=127.0.0.1",
		"DATABASE_PORT=3306",
		"DATABASE_NAME=altair_development",
		"DATABASE_USERNAME=root",
		"DATABASE_PASSWORD=root",
	)

	args := []string{"run"}
	if cfg.EnableOauth {
		args = append(args, "--auto-migrate")
	}
	cmd := exec.Command(binPath, args...)
	cmd.Dir = workDir
	cmd.Env = env
	cmd.Stdout = &testWriter{t: t, prefix: "altair stdout"}
	cmd.Stderr = &testWriter{t: t, prefix: "altair stderr"}
	// Put altair in its own process group so we can SIGTERM the whole tree
	// (Cobra spawns nothing today, but defending against future changes).
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		upstream.Close()
		t.Fatalf("start altair: %v", err)
	}

	h := &Harness{
		t:        t,
		BinPath:  binPath,
		WorkDir:  workDir,
		Port:     port,
		Upstream: upstream,
		cmd:      cmd,
	}

	if err := h.waitHealthy(30 * time.Second); err != nil {
		h.Stop()
		t.Fatalf("altair never healthy: %v\nworkdir: %s", err, workDir)
	}
	return h
}

// Stop is idempotent: tears down the subprocess and the mock upstream.
// Tests should `defer h.Stop()` immediately after Start.
func (h *Harness) Stop() {
	if h == nil {
		return
	}
	if h.cmd != nil && h.cmd.Process != nil {
		_ = syscall.Kill(-h.cmd.Process.Pid, syscall.SIGTERM)
		done := make(chan error, 1)
		go func() { done <- h.cmd.Wait() }()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			_ = syscall.Kill(-h.cmd.Process.Pid, syscall.SIGKILL)
			<-done
		}
	}
	if h.Upstream != nil {
		h.Upstream.Close()
	}
}

// URL builds an absolute URL against the running gateway.
func (h *Harness) URL(path string) string {
	return fmt.Sprintf("http://127.0.0.1:%d%s", h.Port, path)
}

func (h *Harness) waitHealthy(deadline time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	url := h.URL("/health")
	for {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}
		if ctx.Err() != nil {
			return fmt.Errorf("timed out polling %s", url)
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func patchAppYAML(workDir string, port int, cfg Config) error {
	path := filepath.Join(workDir, "config", "app.yml")
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	src := string(raw)

	// Replace the templated port with our allocated one. Keep the surrounding
	// formatting; sample uses the literal `port: 1304`.
	src = strings.Replace(src, "port: 1304", fmt.Sprintf("port: %d", port), 1)

	// Override the proxy block with the test-specific settings. The
	// generated file uses `host: {{ env "PROXY_HOST" }}` so we substitute
	// the literal value and override the remaining knobs.
	proxyBlock := fmt.Sprintf(
		"proxy:\n  host: %s\n  upstream_timeout: %s\n  max_request_body_size: %s\n",
		cfg.ProxyHost, cfg.UpstreamTimeout, cfg.MaxRequestBodySize,
	)
	src = replaceBlock(src, "proxy:\n", "authorization:\n", proxyBlock+"authorization:\n")

	// Phase A drops oauth from the active plugin list (the AND-gate in
	// plugin/runner.go means the oauth.yml file alone won't activate it).
	// Phase B leaves it in so the oauth flow loads.
	if !cfg.EnableOauth {
		src = strings.Replace(src, "  - oauth\n", "", 1)
	}

	return os.WriteFile(path, []byte(src), 0o644)
}

// replaceBlock replaces a contiguous segment whose start matches `start` and
// continues up to (but not including) the next `boundary` line. Used to
// swap the multi-line `proxy:` block atomically without a brittle line-by-
// line edit.
func replaceBlock(src, start, boundary, replacement string) string {
	startIdx := strings.Index(src, start)
	if startIdx < 0 {
		return src
	}
	endIdx := strings.Index(src[startIdx:], boundary)
	if endIdx < 0 {
		return src
	}
	return src[:startIdx] + replacement + src[startIdx+endIdx+len(boundary):]
}

func patchRoutes(workDir, upstreamHostPort string, enableOauth bool) error {
	// The default route gets a mix of auth=none paths (for the gateway-
	// only smoke matrix) and, when oauth is enabled, a single
	// auth=oauth path with scope=users so subtests can validate the
	// downstream auth-and-scope check.
	content := fmt.Sprintf(`name: users
auth: none
prefix: /users
host: %s
path:
  /profiles/:id:
    auth: "none"
  /echo:
    auth: "none"
  /slow-3s/:id:
    auth: "none"
`, upstreamHostPort)
	if enableOauth {
		content += `  /me:
    auth: "oauth"
    scope: "users"
`
	}
	return os.WriteFile(filepath.Join(workDir, "routes", "service-a.yml"), []byte(content), 0o644)
}

// testWriter funnels child-process output into t.Log so CI failures surface
// the altair logs alongside the assertion that tripped.
type testWriter struct {
	t      *testing.T
	prefix string
}

func (w *testWriter) Write(p []byte) (int, error) {
	w.t.Logf("[%s] %s", w.prefix, strings.TrimRight(string(p), "\n"))
	return len(p), nil
}

// DrainBody is a tiny helper for subtests: read + close, return body bytes.
func DrainBody(t *testing.T, r io.ReadCloser) []byte {
	t.Helper()
	defer r.Close()
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return b
}
