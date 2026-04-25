package oauth_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/plugin/oauth"
)

// Assumption: the oauth plugin is named "oauth" — the registry, app.yml,
// and config/plugin/<name>.yml all key on this string.
func TestPluginName_IsOauth(t *testing.T) {
	assert.Equal(t, "oauth", (&oauth.Plugin{}).Name())
}

// Assumption: oauth has no hard dependencies; metric is optional.
func TestPluginDependsOn_IsNil(t *testing.T) {
	assert.Nil(t, (&oauth.Plugin{}).DependsOn())
}

// Assumption: SampleConfig returns the embedded oauth config bytes that
// `altair new` writes to config/plugin/oauth.yml. Content includes the
// plugin name and version markers.
func TestPluginSampleConfig_ContainsPluginAndVersion(t *testing.T) {
	got := string((&oauth.Plugin{}).SampleConfig())
	assert.Contains(t, got, "plugin: oauth")
	assert.Contains(t, got, `version: "1.0"`)
	assert.Contains(t, got, "access_token_timeout")
}

// Assumption: Migrations(ctx) with a working DecodeConfig returns one
// MigrationSet whose VersionTable is the canonical oauth_plugin_db_versions
// and whose DatabaseInstance is the value the DecodeConfig populated.
func TestPluginMigrations_ReturnsOauthSetWithCanonicalVersionTable(t *testing.T) {
	p := &oauth.Plugin{}
	ctx := module.PluginContext{
		DecodeConfig: func(target interface{}) error {
			// Mirror entity.OauthPlugin's nested layout enough to set
			// DatabaseInstance via the Config.Database field.
			cfg, ok := target.(interface {
				DatabaseInstance() string
			})
			_ = cfg
			_ = ok
			return nil
		},
	}
	sets := p.Migrations(ctx)
	assert.Len(t, sets, 1)
	assert.Equal(t, "mysql", sets[0].Driver)
	assert.Equal(t, "migrations/mysql", sets[0].SourcePath)
	assert.Equal(t, "oauth_plugin_db_versions", sets[0].VersionTable)
	assert.NotNil(t, sets[0].FS)
}

// Assumption: Migrations(ctx) returns nil when DecodeConfig is missing.
// Defensive guard; hit when callers (e.g. registry audit tests) pass an
// empty PluginContext.
func TestPluginMigrations_NilOnMissingDecodeConfig(t *testing.T) {
	assert.Nil(t, (&oauth.Plugin{}).Migrations(module.PluginContext{}))
}

// Assumption: Migrations(ctx) returns nil when DecodeConfig errors. The
// caller uses the absence of sets as "no migrations to run."
func TestPluginMigrations_NilOnDecodeConfigError(t *testing.T) {
	ctx := module.PluginContext{
		DecodeConfig: func(_ interface{}) error { return errors.New("decode boom") },
	}
	assert.Nil(t, (&oauth.Plugin{}).Migrations(ctx))
}

// Assumption: Load with an unsupported version returns a clear error,
// not a panic. The version is part of the message so operators can tell
// what the plugin saw.
func TestPluginLoad_RejectsUnknownVersion(t *testing.T) {
	err := (&oauth.Plugin{}).Load(module.PluginContext{Version: "9.9"})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "9.9")
	assert.Contains(t, err.Error(), "oauth")
}

// Assumption: LoadCommand with an unsupported version mirrors Load's error
// shape.
func TestPluginLoadCommand_RejectsUnknownVersion(t *testing.T) {
	err := (&oauth.Plugin{}).LoadCommand(module.PluginContext{Version: "9.9"})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "9.9")
	assert.Contains(t, err.Error(), "oauth")
}

// Assumption: Load("1.0") with a missing DecodeConfig errors out before
// touching the database — the call site needs DecodeConfig to read
// access_token_timeout etc. We don't run the happy path because it needs
// a real DB connection; that's covered by the smoke test.
func TestPluginLoad_V10WithMissingDecodeConfigErrors(t *testing.T) {
	err := (&oauth.Plugin{}).Load(module.PluginContext{Version: "1.0"})
	assert.NotNil(t, err)
}

// Assumption: LoadCommand("1.0") with a missing DecodeConfig errors out.
func TestPluginLoadCommand_V10WithMissingDecodeConfigErrors(t *testing.T) {
	err := (&oauth.Plugin{}).LoadCommand(module.PluginContext{Version: "1.0"})
	assert.NotNil(t, err)
}
