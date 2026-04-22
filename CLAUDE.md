# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository

Altair is a Go-based API gateway (entry point `altair.go`, root package `main`). It depends on the sibling `github.com/kodefluence/monorepo` library for `kontext`, `exception`, `jsonapi`, `db`, and `memorystore`. Go 1.19, Gin, zerolog, Cobra, MySQL driver.

## Commands

All commands are driven by the `Makefile`; prefer them over ad-hoc `go` invocations because `make test` excludes packages that don't contribute coverage (`altair$`, `core`, `mock`, `interfaces`, `testhelper`).

- `make test` — full test suite with race detector and `cover.out`.
- Run a single test: `go test -run TestName ./module/router/usecase/...` (pass `-race` to match CI).
- `make mock_all` — regenerate every mock (`mock_metric`, `mock_plugin`, `mock_loader`, `mock_routing` plus service/formatter/model/validator targets). Individual mocks use `mockgen -source core/<file>.go -destination mock/mock_<file>.go -package mock`; sub-package mocks are driven by `//go:generate mockgen ...` directives co-located with the interface they mock.
- Build: `make build_linux` / `build_darwin` / `build_windows` (produce UPX-packed binaries in `./build_output/<os>/altair`). `make build` runs all three.
- Docker: `make build_docker` / `build_docker_latest` / `push_docker`. The `Dockerfile` expects `build_linux` to have run first — it COPYs `./build_output/linux/altair` into the image.
- Local stack: copy `env.sample` to `.env`, then `make docker-compose-up` (MySQL 5.7 on 127.0.0.1:3306). Pair it with `altair run` in another terminal.
- Linting matches CI via `golangci-lint run` (CI pins `v1.50.1`).

## Runtime entry

`altair.go` is a Cobra root with three subcommands:
- `altair run` — loads `config/app.yml`, `config/database.yml`, `config/plugin/*`, fabricates DB connections, mounts healthcheck + plugins, serves on `appConfig.Port()` (default 1304). Graceful shutdown listens for SIGINT/SIGTERM.
- `altair config [app|db|all]` — dumps loaded config.
- `altair plugin ...` — runs plugin-provided Cobra subcommands (migrations, application creation, etc.). `DisableFlagParsing: true` so each plugin owns its flag surface.

All three share `loadConfig()` at startup, so missing `config/` files silently yield `nil` configs that each subcommand checks before executing.

## Architecture

The repo is a strict five-layer, inward-pointing design. Dependencies flow only toward `core` and `entity`; nothing in `core` or `entity` imports outward.

```
core/      — interfaces only (AppConfig, Controller, Metric, DownStreamPlugin, RouteCompiler, …)
entity/    — pure data types + constructor option structs (AppConfig, RouteObject, DBConfig)
adapter/   — wraps external concretes into core.* interfaces (entity.AppConfig -> core.AppConfig)
cfg/       — YAML loaders + "bearer" state holders; env-interpolated via text/template with {{ env "FOO" }}
module/    — feature assemblies; every feature has provider.go + usecase/ (+ controller/ when HTTP-facing)
plugin/    — versioned extensions (oauth, metric) with their own mini-layering inside
```

### Module pattern

Every module exposes a single entry: `Provide(...)` (returns a `module.*` interface) or `Load(appModule)` (mutates the controller registry). Internally:
- `module/<name>/provider.go` — thin DI: constructs usecase and returns it typed as the `module.*` interface.
- `module/<name>/usecase/` — business logic; one file per struct, each with its own `_test.go`. Usecases consume consumer-owned interfaces defined in the same file (see `module/oauth/.../application_manager.go` — `OauthApplicationRepository` and `Formatter` are declared there, not in the repo package).
- `module/<name>/controller/` — HTTP/command/downstream adapters (only for features with transport surfaces, e.g. `healthcheck`, `projectgenerator`, or inside plugins).

`module/interface.go` is the shared contract: `App`, `Controller`, `HttpController`, `CommandController`, `DownstreamController`, `MetricController`, `ApiError`, `RouterPath`.

### Request lifecycle

`module/controller/usecase/http.go` `InjectHTTP` wraps every registered `HttpController`:
1. Generates `request_id` (uuid), sets `start_time`, fabricates a `kontext.Context` wrapping the gin context.
2. Buffers request body so downstream handlers can re-read it.
3. `defer ctrl.httpRecoverFunc(...)` catches panics, emits an Internal Server Error JSON:API response, logs with structured tags.
4. Dispatches to `httpController.Control(ktx, c)`.
5. Branches on `c.Writer.Status()` to log info vs error with elapsed time.
6. Increments `controller_hits` counter + `controller_elapsed_time_seconds` histogram against every registered `MetricController`.

Controllers themselves never touch loggers or metrics — they just call their manager and `c.JSON(status, jsonapi.BuildResponse(...))`.

### Routing / proxy

`module/router/usecase/compiler.go` walks `./routes/*.{yml,yaml}` with `text/template` (`{{ env ... }}` helper) and unmarshals each into an `entity.RouteObject` containing per-path `Auth`/`Scope`. `generator.go` registers one Gin `engine.Any(urlPath, ...)` per path; at request time it:
1. Builds a `proxyReq` copy of the client request.
2. Runs every registered `DownstreamController.Intervene` (this is where the oauth plugin validates tokens and rejects early).
3. Calls the upstream, copies headers/body back.
4. Observes metrics `routes_downstream_hits`, `routes_downstream_latency_seconds`, and per-plugin `routes_downstream_plugin_latency_seconds`.

### Plugins

`plugin/<name>/loader.go` exports `Load(...)` and sometimes `LoadCommand(...)`. Each delegates to a version function selected via `switch pluginBearer.PluginVersion(name)` — see `plugin/oauth/loader.go` dispatching to `version_1_0(...)` in `plugin/oauth/version_1.0.go`. Adding a new plugin version means adding `version_1_1.go` and a case in the switch; old versions remain intact.

Inside a plugin, the structure mirrors the top-level:
- `plugin/<name>/entity/` — data + plugin-specific config.
- `plugin/<name>/module/<feature>/` — own `loader.go` + `usecase/` + `controller/{http,command,downstream}/`.
- `plugin/<name>/repository/mysql/` — one file per table; every method takes `(ktx kontext.Context, ..., tx db.TX)` and returns `(..., exception.Exception)`.

### Config loading

`cfg/` exposes `App()`, `Database()`, `Plugin()` factories each implementing a `core.*Loader` interface. Each `Compile(path)`:
1. Reads the file, runs it through `compileTemplate` which supports `{{ env "VAR" }}` with empty-string fallback.
2. Unmarshals into a versioned base struct (`Version string`).
3. Switches on version, validates, applies defaults (e.g. port 1304, proxy host `www.local.host`), and constructs the entity via `entity.New*` then wraps with `adapter.*Config`.

`cfg/app_bearer.go` and `cfg/database_bearer.go` are mutable registries handed to plugin code so plugins can look up DB instances and inject `DownStreamPlugin`s.

## Go conventions in this repo

- **Interface ownership.** Every interface lives where it is *consumed*, not where it is implemented. `core/` holds cross-cutting contracts; usecase files declare the narrow interfaces they need (`OauthApplicationRepository`, `Formatter`, `ApplicationManager`). Do not export "one big interface" from a repository or adapter package.
- **Constructors.** `Fabricate*` creates or reuses a `sync.Map`-cached singleton bound to an external resource (DB, memcache). `New*` is a plain struct constructor with no side effects. `Adapt*` wraps a concrete external type into one of our interfaces. `Provide*` is the DI entry into a `module/` package. Don't mix.
- **Functional options.** Configurable constructors always take `(required..., opts ...Option)` with `Option func(*Config)`. Each knob gets a `WithXxx` helper. Do not add positional parameters — add an option.
- **Errors.** Below the HTTP controller layer, return `exception.Exception` (from `github.com/kodefluence/monorepo/exception`) — never `error`. Set `WithType(exception.NotFound|BadInput|Unauthorized|Forbidden|Duplicated|Unexpected)` so controllers can branch. At the controller, convert via `cr.apiError.XxxError(...)` which produces a `jsonapi.Option`, then `c.JSON(status, jsonapi.BuildResponse(opt))`. Panics are recovered centrally — do not recover locally.
- **Context.** Thread `kontext.Context` (not `context.Context`) through every usecase → repo call. `kontext.Fabricate(kontext.WithDefaultContext(ctxWithTimeout))` wraps a derived stdlib context when you need a per-query timeout. `request_id` and `start_time` are standard keys set by the HTTP middleware.
- **DB access.** Every repo method takes `(ktx kontext.Context, …, tx db.TX)` and returns `exception.Exception`. Wrap the operation in `context.WithTimeout(ktx.Ctx(), 10*time.Second)` (or appropriate) and rebuild the kontext before calling `tx.QueryContext`/`ExecContext`. Transactions use `sqldb.Transaction(ktx, "key", func(tx db.TX) exception.Exception { ... })` — rollback on non-nil return is automatic.
- **Logging.** Use `zerolog` with tag arrays: `log.Error().Err(err).Stack().Array("tags", zerolog.Arr().Str("altair").Str("controller").Str(httpController.Path())).Msg("...")`. Tags are the primary search dimension; include the component name, the subcomponent, and the action.
- **Plugin versioning.** Never modify `version_1_0.go` after release — add `version_1_1.go`. The `loader.go` switch is the single source of truth for which version runs.
- **Generics.** Constrained generics are fine where they remove repetition (see `util.ValueToPointer[V Value]`); keep the constraint explicit (`type Value interface { int | string | time.Time }`), don't reach for `any`.
- **Mocks.** Co-locate a `//go:generate mockgen -source <file>.go -destination ./mock/mock.go -package mock` directive at the top of every file that declares a consumer interface. Generated mocks land in a sibling `mock/` package. Keep them committed.

## Pull requests & contribution

`CONTRIBUTING.md` is authoritative: discuss via issue first, follow Chris Beams' commit-message conventions, keep `README.md` in sync with feature additions, one reviewer sign-off required before merge.
