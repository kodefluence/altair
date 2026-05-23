//go:build e2e

package harness

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLDSN is the DSN the harness uses to seed the oauth_applications
// table directly. Matches docker-compose.yaml locally and the `services:`
// MySQL block in CI: `root` / `root` on `127.0.0.1:3306`, db
// `altair_development`. The harness DSN deliberately mismatches
// env.sample's `rootpw` because env.sample is broken (see the smoke spec
// at docs/superpowers/specs/2026-04-23-altair-smoke-test-design.md).
const MySQLDSN = "root:root@tcp(127.0.0.1:3306)/altair_development?parseTime=true&multiStatements=true"

// OauthApp identifies a seeded application along with the bearer it just
// exchanged. Tests use Token in the Authorization header; ClientUID/Secret
// are kept around so individual subtests can run their own grant exchanges.
type OauthApp struct {
	ID           int64
	ClientUID    string
	ClientSecret string
	Scopes       string
	Token        string // populated by ExchangeToken
}

// OpenMySQL opens a connection to the smoke-test database and pings it.
// Callers are expected to defer Close. Test failure on connect is the
// right signal — Phase B is gated on MySQL being available.
func OpenMySQL(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("mysql", MySQLDSN)
	if err != nil {
		t.Fatalf("open mysql: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		t.Fatalf("mysql unreachable at 127.0.0.1:3306: %v\n"+
			"local: run `docker compose --env-file .env up -d`\n"+
			"CI:    job needs `services: mysql:` block", err)
	}
	return db
}

// SeedApplication INSERTs a randomly-named oauth_application with the given
// scopes (space-separated, per the schema). Returns the row including the
// generated client_uid + client_secret so the caller can exchange them.
func SeedApplication(t *testing.T, db *sql.DB, scopes string) *OauthApp {
	t.Helper()

	clientUID := "smoke-uid-" + randomHex(8)
	clientSecret := "smoke-secret-" + randomHex(16)

	res, err := db.Exec(`
		INSERT INTO oauth_applications
			(owner_type, scopes, client_uid, client_secret, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`, "User", scopes, clientUID, clientSecret, "altair smoke test")
	if err != nil {
		t.Fatalf("seed oauth_application: %v", err)
	}
	id, _ := res.LastInsertId()

	return &OauthApp{
		ID:           id,
		ClientUID:    clientUID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
	}
}

// ResetOauthTables drops migration state so the next altair --auto-migrate
// boot re-runs the migrations cleanly. Idempotent; safe across repeated
// smoke runs and across CI's clean-room invocations.
func ResetOauthTables(t *testing.T, db *sql.DB) {
	t.Helper()
	tables := []string{
		// Drop in FK-safe order: refresh_tokens -> access_tokens ->
		// access_grants -> applications. Then the migrate state table.
		"oauth_refresh_tokens",
		"oauth_access_tokens",
		"oauth_access_grants",
		"oauth_applications",
		// oauth-plugin migration state must drop too — without this,
		// migrate sees the version table populated and skips re-applying
		// migrations on a DB whose data tables we just blew away.
		"oauth_plugin_db_versions",
	}
	for _, t := range tables {
		_, _ = db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", t))
	}
}

// BasicAuthUser / BasicAuthPass match what harness.Start injects via
// BASIC_AUTH_PASSWORD. The plugin endpoints (everything under /_plugins/)
// are gated by gin.BasicAuth in altair.go — the plugin admin surface, not
// the public proxy surface.
const (
	BasicAuthUser = "altair"
	BasicAuthPass = "smoke"
)

// ExchangeToken hits POST /_plugins/oauth/authorizations/token with
// grant_type=client_credentials. The endpoint requires HTTP Basic Auth
// (the plugin admin auth from app.yml `authorization:` block) on top of
// the client_uid / client_secret in the JSON body — Basic Auth gates
// access to the *issuance* surface, the body fields identify *which*
// application is asking. Returns the bearer the gateway expects on the
// public proxy routes' Authorization header.
func ExchangeToken(t *testing.T, gatewayBase string, app *OauthApp, scope string) string {
	t.Helper()

	body, _ := json.Marshal(map[string]string{
		"grant_type":    "client_credentials",
		"client_uid":    app.ClientUID,
		"client_secret": app.ClientSecret,
		"scope":         scope,
	})

	url := gatewayBase + "/_plugins/oauth/authorizations/token"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("build token request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(BasicAuthUser, BasicAuthPass)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("POST %s: %v", url, err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("token exchange: want 200, got %d. body=%s", resp.StatusCode, raw)
	}

	var envelope struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		t.Fatalf("decode token envelope: %v\nbody=%s", err, raw)
	}
	if envelope.Data.Token == "" {
		t.Fatalf("token exchange: empty token in response. body=%s", raw)
	}
	return envelope.Data.Token
}

func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
