package module

import (
	"io/fs"

	"github.com/kodefluence/monorepo/db"
	"github.com/rs/zerolog"

	"github.com/kodefluence/altair/core"
)

// Plugin is the contract every built-in Altair plugin implements. The registry
// in package plugin iterates []Plugin for Load, LoadCommand, migrations, and
// `altair new` sample generation. Plugin lives in module/ rather than core/
// because PluginContext references module.App and module.ApiError.
type Plugin interface {
	// Name is the stable identifier used in app.yml `plugins:` and as the
	// filename of the plugin's config under config/plugin/<name>.yml.
	Name() string

	// DependsOn lists plugin names (by Name()) that must load before this one
	// when both are active. Semantics are soft: a named dependency that is not
	// active is skipped without error.
	DependsOn() []string

	// Migrations returns the migration sets owned by this plugin. An empty
	// slice signals "no schema". DatabaseInstance is resolved from the parsed
	// plugin config, which is why this is a method rather than a field.
	Migrations(ctx PluginContext) []MigrationSet

	// Load registers HTTP, downstream, and metric controllers for `altair run`.
	Load(ctx PluginContext) error

	// LoadCommand registers CLI subcommands under `altair plugin ...`.
	// The generic migration commands are injected once by module/migration,
	// not per-plugin here.
	LoadCommand(ctx PluginContext) error

	// SampleConfig returns the default YAML bytes that `altair new` writes to
	// config/plugin/<name>.yml. Plugins typically supply this via //go:embed.
	SampleConfig() []byte
}

// PluginContext is the single argument handed to a Plugin's Load and
// LoadCommand methods. It hides raw bearers so the framework can enrich calls
// (logger tags, config-subtree decoding) without breaking plugin implementations.
type PluginContext struct {
	// Version is the value of the `version:` field from the plugin's YAML
	// file. Plugins switch on it internally to support multiple schemas.
	Version string

	// DecodeConfig unmarshals this plugin's YAML config into target. The
	// closure hides the yaml dependency from plugin packages and lets the
	// framework evolve parsing without touching plugin code.
	DecodeConfig func(target interface{}) error

	// Database resolves a fabricated db.DB and its config by instance name
	// (the string referenced in the plugin's `database:` config field).
	Database func(instance string) (db.DB, core.DatabaseConfig, error)

	// Logger is a zerolog logger pre-tagged with ["altair","plugin",<name>].
	Logger zerolog.Logger

	AppBearer core.AppBearer
	AppModule App
	ApiError  ApiError
}

// MigrationSet describes one embedded migration source belonging to a plugin.
// Plugins return []MigrationSet so a plugin may own zero, one, or many
// migration sources (e.g. per-driver or per-database).
type MigrationSet struct {
	// DatabaseInstance names the entry in database.yml this migration applies to.
	DatabaseInstance string

	// Driver is the SQL driver name (today: "mysql").
	Driver string

	// FS is the embed.FS (or any io/fs.FS) containing the migration files.
	FS fs.FS

	// SourcePath is the directory inside FS that holds the numbered migrations.
	SourcePath string

	// VersionTable is the golang-migrate migrations table. Convention:
	// "<plugin>_plugin_db_versions". Override only when migrating across
	// plugin renames.
	VersionTable string
}
