//go:build e2e

package e2e_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/kodefluence/altair/e2e/harness"
)

// startOauthGateway brings up an altair instance with the oauth plugin
// active and --auto-migrate. Resets oauth_* tables before boot so the
// migration applies cleanly regardless of prior runs. The returned db
// handle stays open for the duration of the test so subtests can seed
// their own application rows.
func startOauthGateway(t *testing.T) (*harness.Harness, *sql.DB) {
	t.Helper()

	db := harness.OpenMySQL(t)
	harness.ResetOauthTables(t, db)

	h := harness.Start(t, binPath, harness.Config{
		ProxyHost:          "altair.test",
		UpstreamTimeout:    "10s",
		MaxRequestBodySize: "1MB",
		EnableOauth:        true,
	})
	t.Cleanup(h.Stop)
	t.Cleanup(func() {
		harness.ResetOauthTables(t, db)
		_ = db.Close()
	})
	return h, db
}

// T2 — oauth_happy_path: a freshly-issued bearer with the right scope
// flows through the gateway and reaches the upstream.
func TestSmoke_OauthHappyPath(t *testing.T) {
	h, db := startOauthGateway(t)

	app := harness.SeedApplication(t, db, "users")
	app.Token = harness.ExchangeToken(t, h.URL(""), app, "users")

	req, _ := http.NewRequest(http.MethodGet, h.URL("/users/me"), nil)
	req.Header.Set("Authorization", "Bearer "+app.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("GET /users/me: %v", err)
	}
	body := harness.DrainBody(t, resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: want 200, got %d. body=%s", resp.StatusCode, body)
	}
	var echo harness.Echo
	if err := json.Unmarshal(body, &echo); err != nil {
		t.Fatalf("decode echo: %v", err)
	}
	if echo.Path != "/users/me" {
		t.Errorf("upstream path: want /users/me, got %q", echo.Path)
	}
}

// T3 — oauth_missing_token: a request to a protected route without an
// Authorization header is rejected before reaching the upstream.
func TestSmoke_OauthMissingToken(t *testing.T) {
	h, _ := startOauthGateway(t)

	resp, err := http.Get(h.URL("/users/me"))
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	resp.Body.Close()

	// Anything in the 4xx range is a pass — exact code (401 vs 403)
	// depends on the downstream plugin's choice; what matters is that
	// the request was rejected and didn't fall through to a 200/echo.
	if resp.StatusCode < 400 || resp.StatusCode >= 500 {
		t.Fatalf("status: want 4xx, got %d", resp.StatusCode)
	}
}

// T4 — oauth_invalid_token: a syntactically present but unknown bearer
// is rejected. Same shape as T3 but exercises the "token not found" path.
func TestSmoke_OauthInvalidToken(t *testing.T) {
	h, _ := startOauthGateway(t)

	req, _ := http.NewRequest(http.MethodGet, h.URL("/users/me"), nil)
	req.Header.Set("Authorization", "Bearer this-is-not-a-real-token")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode < 400 || resp.StatusCode >= 500 {
		t.Fatalf("status: want 4xx, got %d", resp.StatusCode)
	}
}

// T5 — oauth_wrong_scope: token issued with one scope is rejected on a
// route that requires a different scope. Pins the scope-check downstream
// plugin behaviour — a regression in `plugin/oauth/module/authorization/
// controller/downstream/` would surface here.
func TestSmoke_OauthWrongScope(t *testing.T) {
	h, db := startOauthGateway(t)

	app := harness.SeedApplication(t, db, "wrong-scope")
	app.Token = harness.ExchangeToken(t, h.URL(""), app, "wrong-scope")

	req, _ := http.NewRequest(http.MethodGet, h.URL("/users/me"), nil)
	req.Header.Set("Authorization", "Bearer "+app.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode < 400 || resp.StatusCode >= 500 {
		t.Fatalf("status: want 4xx (scope mismatch), got %d", resp.StatusCode)
	}
}

// T6 — oauth_body_and_headers: an authenticated POST round-trips both a
// custom header and the request body verbatim to the upstream.
func TestSmoke_OauthBodyAndHeaders(t *testing.T) {
	h, db := startOauthGateway(t)

	app := harness.SeedApplication(t, db, "users")
	app.Token = harness.ExchangeToken(t, h.URL(""), app, "users")

	payload := `{"hello":"world","n":42}`
	req, _ := http.NewRequest(http.MethodPost, h.URL("/users/me"), strings.NewReader(payload))
	req.Header.Set("Authorization", "Bearer "+app.Token)
	req.Header.Set("X-Trace", "smoke-trace-abc")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	body := harness.DrainBody(t, resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: want 200, got %d. body=%s", resp.StatusCode, body)
	}
	var echo harness.Echo
	if err := json.Unmarshal(body, &echo); err != nil {
		t.Fatalf("decode echo: %v", err)
	}
	if echo.Body != payload {
		t.Errorf("upstream body: want %q, got %q", payload, echo.Body)
	}
	if got := echo.Headers["X-Trace"]; len(got) == 0 || got[0] != "smoke-trace-abc" {
		t.Errorf("upstream missing X-Trace header; got: %v", echo.Headers)
	}
}
