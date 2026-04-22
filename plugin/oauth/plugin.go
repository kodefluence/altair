package oauth

import (
	"embed"
	_ "embed"
	"fmt"

	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/plugin/oauth/entity"
)

//go:embed config.sample.yml
var sampleConfig []byte

//go:embed migrations/mysql/*.sql
var migrationsFS embed.FS

// Plugin implements module.Plugin for the oauth plugin. Wiring lives in
// version_1_0.go and is selected by Load/LoadCommand based on PluginContext.Version.
type Plugin struct{}

// Name implements module.Plugin.
func (*Plugin) Name() string { return "oauth" }

// DependsOn implements module.Plugin. The oauth plugin does not hard-depend on
// any other plugin; metric collection degrades to a no-op fallback when metric
// is inactive.
func (*Plugin) DependsOn() []string { return nil }

// Migrations implements module.Plugin. The target database instance is read
// from the plugin's `config.database` field; the version-tracking table is
// set explicitly to "oauth_plugin_db_versions" so the value is stable even
// if the plugin is ever renamed.
func (*Plugin) Migrations(ctx module.PluginContext) []module.MigrationSet {
	if ctx.DecodeConfig == nil {
		return nil
	}
	var cfg entity.OauthPlugin
	if err := ctx.DecodeConfig(&cfg); err != nil {
		return nil
	}
	return []module.MigrationSet{{
		DatabaseInstance: cfg.DatabaseInstance(),
		Driver:           "mysql",
		FS:               migrationsFS,
		SourcePath:       "migrations/mysql",
		VersionTable:     "oauth_plugin_db_versions",
	}}
}

// SampleConfig implements module.Plugin.
func (*Plugin) SampleConfig() []byte { return sampleConfig }

// Load implements module.Plugin and dispatches on PluginContext.Version.
func (*Plugin) Load(ctx module.PluginContext) error {
	switch ctx.Version {
	case "1.0":
		return loadV1_0(ctx)
	default:
		return fmt.Errorf("undefined template version: %s for oauth plugin", ctx.Version)
	}
}

// LoadCommand implements module.Plugin and dispatches on PluginContext.Version.
func (*Plugin) LoadCommand(ctx module.PluginContext) error {
	switch ctx.Version {
	case "1.0":
		return loadCommandV1_0(ctx)
	default:
		return fmt.Errorf("undefined template version: %s for oauth plugin", ctx.Version)
	}
}
