# Altair end-to-end smoke test — design

**Status:** design approved, pending implementation plan.
**Date:** 2026-04-23.

## Context

Altair is a language-agnostic API gateway where the binary is the single
control point: operators run `altair new <dir>` to scaffold a project and
`altair run` to serve traffic. Today no automated check exercises that
full cycle. Unit tests cover individual packages but not the real
binary, not the generated config, not the oauth plugin's token + scope
flow, and not the actual proxy forwarding. A regression in the
generator, the plugin loader, the proxy logic, or the oauth
authorization code would slip past `make test` and surface as a broken
Docker image or a failed deploy.

This spec defines a smoke test that drives `altair new` → `altair run` →
HTTP requests through the full auth matrix against a local mock
upstream, end-to-end.

## Goals

- Catch regressions in: project generation, plugin registry loading,
  config parsing, route compilation, proxy forwarding (body + headers),
  oauth token issuance, scope enforcement.
- Runs in both CI (`general.yml`) and locally (`make smoke`). Same code
  path in both.
- Completes in ≤3 minutes wall time.
- Easy to extend when new routes, plugins, or auth cases land.

## Non-goals

- Performance / load testing (no k6, no vegeta).
- Bad-gateway / upstream-failure handling (follow-up iteration).
- Concurrent-request contention (not a smoke concern).
- Multi-service, multi-route-file topologies (current sample has one;
  expand when altair grows it).
- Rate-limit or other plugin-chain interactions (no such plugins exist
  today).

## High-level approach

Go `*_test.go` harness in a new `e2e/` package gated by
`//go:build e2e`. MySQL comes from the existing `docker-compose.yaml`
locally, from GitHub Actions' `services:` block in CI. The mock upstream
is an in-process `httptest.Server` that echoes the incoming request as
JSON. No new Go dependencies.

Rejected alternatives:
- **testcontainers-go**: adds ~20 transitive deps for functionality the
  repo's existing `docker-compose.yaml` already provides.
- **docker-compose orchestrating altair itself + test runner**: two
  orchestration layers, slower dev loop, container rebuild cycle on
  every code change.
- **Hurl / venom DSL**: adds a runtime dependency; weak typing for
  assertions beyond status codes; can't drive SQL seeding.

## Architecture

```
e2e/
  harness/
    harness.go       # build binary, write config, spawn altair,
                     # wait for healthcheck, orchestrate teardown
    upstream.go      # httptest.Server that echoes {method, path,
                     # headers, body} as JSON
    oauth.go         # seed oauth_application via SQL; exchange
                     # client_credentials for a bearer token
    ports.go         # allocate a free TCP port (net.Listen ":0")
  smoke_test.go      # TestMain + subtest matrix
  testdata/
    app.yml.tmpl     # reserved for overrides beyond `altair new`
  README.md          # run instructions + failure-mode decision tree

scripts/
  wait-for-mysql.sh  # 10-line `nc -z` loop for local dev
```

Build tag `//go:build e2e` on every file in `e2e/` keeps the default
test run fast. Smoke runs only via `go test -tags=e2e ./e2e/...`.

### Component responsibilities

- **`harness.Harness`** — one struct, one lifecycle (`Start` / `Stop`).
  Owns the subprocess, the upstream server, the tmpdir, and the DB
  handle for seeding. Subtests receive a `*Harness` and nothing else.
- **`upstream.Server`** — wraps `httptest.NewServer` with a handler that
  JSON-encodes the incoming request into a response with this shape:
  `{"method":"GET","path":"/users/me","headers":{"X-Request-Id":["…"],…},"body":"…"}`.
  Body is captured as a string (not base64) because the sample routes
  only move text through. Exposes `.URL()` and `.Host()` so we can feed
  `EXAMPLE_USERS_SERVICE_HOST` into the generated config.
- **`oauth.Seed`** — opens a `*sql.DB` to the same MySQL the gateway
  talks to, `INSERT`s an oauth_application with deterministic
  `client_uid`/`client_secret` and a `"users"` scope. Then calls
  `POST /oauth/token` with `grant_type=client_credentials`, parses the
  bearer. Returns `{Token, ClientUID, ClientSecret, ApplicationID}`.
- **`ports.FreeTCP`** — `net.Listen("tcp", ":0")`, capture port, close,
  return. Two callers: altair port + mock upstream port. Microsecond
  race acceptable for smoke; escalate to retry loop only if flakes
  appear.
- **`smoke_test.go TestMain`** — one harness shared across subtests.
  Subtests are isolated by HTTP semantics, not infra.

## Test flow

### `TestMain` lifecycle (runs once)

1. **Sanity**: dial `127.0.0.1:3306` with 2s timeout. If dead,
   `t.Skip("docker/mysql unavailable: run `make docker-compose-up`
   first")` when not in CI; hard-fail in CI (MySQL is guaranteed there
   via `services:`).
2. **Clean slate**: `DROP TABLE IF EXISTS` on `oauth_applications`,
   `oauth_access_tokens`, `oauth_access_grants`,
   `oauth_refresh_tokens`, `oauth_plugin_db_versions`. Guarantees
   repeatability across runs.
3. **Build binary**: `go build -o <tmpdir>/altair-e2e .`.
4. **tmpdir**: `t.TempDir()`, chdir in. Retained on failure for
   debugging.
5. **Mock upstream**: start `httptest.Server` on a free port; capture
   `host:port`.
6. **Scaffold**: run `<bin> new .` in tmpdir.
7. **Patch three files post-scaffold**:
   - `config/app.yml` — set `port:` to the free port we allocated.
   - `.env` — set `EXAMPLE_USERS_SERVICE_HOST=<upstream>`, `DATABASE_*`
     to match whatever MySQL environment we're targeting (local
     `docker-compose.yaml` exposes root with password `root`; CI
     `services:` block matches). Note: the committed `env.sample` has
     a latent inconsistency (password `rootpw` doesn't match compose's
     `root`); the smoke test overrides rather than relies on
     `env.sample`. Also set a deterministic `BASIC_AUTH_PASSWORD`.
   - `config/plugin/oauth.yml` — unchanged, the default `main_database`
     is what we want.
8. **Run gateway**: `<bin> run --auto-migrate` as child process with
   stdout/stderr piped through `t.Log`.
9. **Wait healthy**: poll `GET http://127.0.0.1:<port>/healthcheck`
   every 200ms, 30s timeout. Fail loudly if never up.
10. **Seed oauth**: direct `INSERT INTO oauth_applications` with known
    uid/secret + scope `"users"`.
11. **Exchange token**: `POST /oauth/token` with `grant_type=
    client_credentials`. Cache the bearer on the harness.
12. **Run subtests.**
13. **Teardown**: SIGTERM altair → wait up to 5s → SIGKILL fallback.
    Close upstream. `DROP TABLE` cleanup mirror step 2.

### Subtest matrix

| # | Name | Request | Expected |
|---|------|---------|----------|
| T1 | `forwarding_unauthed` | `GET /users/profiles/123` (the path overrides `auth: none`) | 200; echo body contains `method=GET`, `path=/users/profiles/123`; upstream saw altair-injected `X-Request-Id`, `X-Real-Ip-Address`, `X-Forwarded-For` |
| T2 | `forwarding_oauth_happy` | `GET /users/me` + `Bearer <valid>` | 200; echo body matches; `X-Request-Id` present |
| T3 | `oauth_missing_token` | `GET /users/me` (no Authorization header) | 401; JSON:API response with `code=ERR0401` |
| T4 | `oauth_invalid_token` | `GET /users/me` + `Bearer garbage` | 401; JSON:API shape |
| T5 | `oauth_wrong_scope` | Seed a second oauth_application with scope `"wrong"`, get its token, hit `/users/me` | 403 with `code=ERR0403` |
| T6 | `forwarding_body_and_headers` | `POST /users/profiles/123` with header `X-Trace: abc` and body `{"hi":1}` | 200; echo confirms both body and header arrived at upstream |

**Ordering invariant**: T1 runs first because it's the fastest sanity
check (no oauth state needed). If T1 fails, oauth-dependent failures
are noise — put T1 first so diagnostics stay clear.

### Assertion helper

```go
func (h *Harness) Hit(t *testing.T, req *http.Request) (*http.Response, upstream.Echo)
```

Returns the response and a decoded echo struct (empty for non-forwarded
errors like 401). Every subtest uses this so "did the gateway forward
correctly" is uniformly expressed.

## Integration with existing infrastructure

### Makefile

```makefile
smoke: docker-compose-up
	@./scripts/wait-for-mysql.sh 127.0.0.1 3306 30
	go test -tags=e2e -count=1 -v -timeout=5m ./e2e/...

smoke-fast: docker-compose-up
	@./scripts/wait-for-mysql.sh 127.0.0.1 3306 30
	go test -tags=e2e -count=1 -v -run TestSmoke/forwarding_unauthed ./e2e/...
```

`smoke-fast` runs only T1 — the tightest dev-loop feedback.

### CI (new job in `.github/workflows/general.yml`)

```yaml
smoke:
  name: Smoke (e2e)
  runs-on: ubuntu-latest
  needs: verify
  services:
    mysql:
      image: mysql:5.7
      env:
        MYSQL_ROOT_PASSWORD: root
        MYSQL_DATABASE: altair_development
      ports: ['3306:3306']
      options: >-
        --health-cmd="mysqladmin ping -uroot -proot"
        --health-interval=10s --health-timeout=5s --health-retries=10
  steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v5
      with: { go-version: stable, check-latest: true, cache: true }
    - run: cp env.sample .env
    - run: go test -tags=e2e -count=1 -v -timeout=5m ./e2e/...
```

CI uses the native `services:` block instead of docker-compose. Two
reasons: (a) built-in health checks, (b) no docker-in-docker. Local dev
uses `make docker-compose-up`; CI uses `services:`. Both converge on
"MySQL on 127.0.0.1:3306 with altair_development database ready".

`needs: verify` gates on fmt/vet/tidy/generate before spinning up
MySQL — no CI minutes on broken PRs.

## Operational details

1. **Clean slate per run.** `TestMain` drops oauth_* tables and
   `oauth_plugin_db_versions` before and after. Repeatable regardless
   of prior run state.
2. **Subprocess log piping.** Altair's zerolog output reaches `t.Log`
   so CI failures are diagnosable without shell access.
3. **Port allocation.** `net.Listen(":0")` → close → reuse. Acceptable
   for smoke; retry loop if flakes ever appear.
4. **Subprocess lifecycle.** `cmd.Start()` not `cmd.Run()`; capture
   `*os.Process`; SIGTERM → 5s wait → SIGKILL.
5. **Debugging locally.** Test logs print the tmpdir path. Failed runs
   retain the dir (Go test default) — `cd` in and `cat config/app.yml`
   to inspect generated state.
6. **Flake safety.** `-count=1` defeats Go test cache; `-timeout=5m`
   caps worst case.
7. **Skippability.** Missing MySQL → `t.Skip` locally, hard-fail in CI.
   Laptops without docker can still run unit tests with no friction.

### Failure-mode decision tree (in `e2e/README.md`)

| Symptom | Likely cause | Fix |
|---|---|---|
| T1 fails with `dial tcp 127.0.0.1:<port>: connection refused` | altair didn't start / crashed during migrate | Read the piped altair logs above the failure |
| T2 fails with 401 | oauth_application seed didn't commit | Check MySQL connection; check `Seed()` error in harness logs |
| T5 fails with 200 instead of 403 | scope-check in downstream plugin regressed | Real regression — see `plugin/oauth/module/authorization/` |
| `make smoke` fails immediately | MySQL not up | `make docker-compose-up` first |
| Healthcheck timeout | Migration failed; port collision; env var missing | Check altair stdout in test log; re-run with `go test -v` |

## Success criteria

- `make smoke` on a fresh checkout with docker running completes all 6
  subtests in ≤3 minutes.
- CI `smoke` job green on main.
- A deliberate regression — e.g. comment out the scope check in
  `plugin/oauth/module/authorization/controller/downstream/` — causes
  T5 to fail with a readable diagnostic pointing at the code path.
- Future contributors can add T7 (bad-gateway), T8 (new plugin) by
  adding one `t.Run` block and one helper on the harness, without
  touching the setup code.

## Open questions deferred to implementation

These require reading code that's cheaper to inspect when writing the
harness than when designing it:

- Exact stdout format of `altair plugin oauth-application-create` —
  implementation will prefer direct SQL INSERT regardless because it
  decouples the test from CLI output formatting. Call the CLI only if
  there's a strong reason, e.g. to exercise the generator itself.
- Whether the `POST /oauth/token` endpoint is at `/oauth/token` (public)
  or under `/_plugins/oauth/token` (basic-auth-gated) — read
  `plugin/oauth/module/authorization/controller/http/` to confirm.
- Whether `wait-for-mysql.sh` needs a schema-level ping (not just port
  open) — add only if early CI runs flake on connect-before-ready.
- Whether the oauth plugin supports `grant_type=client_credentials` out
  of the box or whether the test needs to use `password` / `refresh`
  grant types — the entity supports refresh and implicit grants per
  `OauthPlugin.Config.RefreshToken` and `ImplicitGrant`; confirm the
  issuance endpoint accepts client_credentials when reading the
  authorization usecase.
