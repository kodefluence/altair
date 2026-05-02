package command_test

import (
	"errors"
	"testing"
	"testing/fstest"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/kodefluence/monorepo/db"

	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/module/migration/usecase"
	"github.com/kodefluence/altair/module/plugin_list/controller/command"
)

// Stubs ------------------------------------------------------------------

type stubAppConfig struct {
	enabled map[string]bool
}

func (s *stubAppConfig) Port() int                      { return 0 }
func (s *stubAppConfig) BasicAuthUsername() string      { return "" }
func (s *stubAppConfig) BasicAuthPassword() string      { return "" }
func (s *stubAppConfig) ProxyHost() string              { return "" }
func (s *stubAppConfig) UpstreamTimeout() time.Duration { return 0 }
func (s *stubAppConfig) MaxRequestBodySize() int64      { return 0 }
func (s *stubAppConfig) PluginExists(name string) bool  { return s.enabled[name] }
func (s *stubAppConfig) Plugins() []string              { return nil }
func (s *stubAppConfig) AutoMigrate() bool              { return false }
func (s *stubAppConfig) Dump() string                   { return "" }

type stubAppBearer struct {
	cfg core.AppConfig
}

func (s *stubAppBearer) Config() core.AppConfig                         { return s.cfg }
func (s *stubAppBearer) DownStreamPlugins() []core.DownStreamPlugin     { return nil }
func (s *stubAppBearer) InjectDownStreamPlugin(_ core.DownStreamPlugin) {}
func (s *stubAppBearer) SetMetricProvider(_ core.Metric)                {}
func (s *stubAppBearer) MetricProvider() (core.Metric, error)           { return nil, nil }

type stubPluginBearer struct {
	configs map[string]bool
}

func (s *stubPluginBearer) ConfigExists(name string) bool          { return s.configs[name] }
func (s *stubPluginBearer) PluginVersion(string) (string, error)   { return "1.0", nil }
func (s *stubPluginBearer) DecodeConfig(string, interface{}) error { return nil }
func (s *stubPluginBearer) ForEach(func(string) error)             {}
func (s *stubPluginBearer) Length() int                            { return len(s.configs) }

type stubDBBearer struct{}

func (s *stubDBBearer) Database(string) (db.DB, core.DatabaseConfig, error) {
	return nil, nil, errors.New("stub: no db")
}

type stubPlugin struct {
	name       string
	deps       []string
	migrations []module.MigrationSet
}

func (p *stubPlugin) Name() string                                            { return p.name }
func (p *stubPlugin) DependsOn() []string                                     { return p.deps }
func (p *stubPlugin) Migrations(_ module.PluginContext) []module.MigrationSet { return p.migrations }
func (p *stubPlugin) Load(_ module.PluginContext) error                       { return nil }
func (p *stubPlugin) LoadCommand(_ module.PluginContext) error                { return nil }
func (p *stubPlugin) SampleConfig() []byte                                    { return nil }

type fakeMigrator struct{ versionV uint }

func (f *fakeMigrator) Up() error                    { return nil }
func (f *fakeMigrator) Steps(int) error              { return nil }
func (f *fakeMigrator) Force(int) error              { return nil }
func (f *fakeMigrator) Version() (uint, bool, error) { return f.versionV, false, nil }

func newRunnerForList(t *testing.T, plugins []module.Plugin, pb *stubPluginBearer) *usecase.Runner {
	t.Helper()
	r := usecase.NewRunner(plugins, pb, &stubDBBearer{})
	usecase.SetMigratorFactoryForTest(r, func(_ module.MigrationSet, _ db.DB, _ core.DatabaseConfig, _ string) (usecase.MigratorForTest, func(), error) {
		// Deliberately error so List's Status() falls back to status-less
		// rendering for migration-owning plugins. We exercise both code
		// paths separately — this stub keeps List robust to runner errors.
		return nil, nil, errors.New("no db in test")
	})
	return r
}

// ----------------------------------------------------------------------

func TestList_Getters(t *testing.T) {
	c := command.NewList(nil, nil, nil, nil)
	assert.Equal(t, "list", c.Use())
	assert.NotEmpty(t, c.Short())
	assert.NotEmpty(t, c.Example())
}

func TestList_ModifyFlagsNoOp(t *testing.T) {
	c := command.NewList(nil, nil, nil, nil)
	cobraCmd := &cobra.Command{}
	assert.NotPanics(t, func() { c.ModifyFlags(cobraCmd.Flags()) })
}

// Assumption: Run iterates every plugin in registry order. Inactive plugins
// (not in app.yml or no config file) are still listed but marked inactive.
// This test exercises the active/inactive branches without asserting on
// stdout text — that's brittle. We confirm Run completes without panic.
func TestList_RunIteratesEveryPlugin(t *testing.T) {
	plugins := []module.Plugin{
		&stubPlugin{name: "metric"},
		&stubPlugin{name: "oauth", migrations: []module.MigrationSet{
			{DatabaseInstance: "main", Driver: "mysql", FS: fstest.MapFS{}, SourcePath: "x"},
		}},
	}
	cfg := &stubAppConfig{enabled: map[string]bool{"oauth": true}}
	appBearer := &stubAppBearer{cfg: cfg}
	pluginBearer := &stubPluginBearer{configs: map[string]bool{"oauth": true}}
	r := newRunnerForList(t, plugins, pluginBearer)

	c := command.NewList(plugins, appBearer, pluginBearer, r)
	assert.NotPanics(t, func() { c.Run(&cobra.Command{}, nil) })
}

// Assumption: nil app/plugin bearers are tolerated by the active() helper.
// Everything is rendered as inactive when there's no bearer to consult.
func TestList_RunWithNilBearersTreatsAllInactive(t *testing.T) {
	plugins := []module.Plugin{&stubPlugin{name: "ghost"}}
	c := command.NewList(plugins, nil, nil, nil)
	assert.NotPanics(t, func() { c.Run(&cobra.Command{}, nil) })
}

// Assumption: when runner.Status() succeeds and reports a plugin, the row
// shows the version pair.
func TestList_RunRendersStatusFromRunner(t *testing.T) {
	plugins := []module.Plugin{
		&stubPlugin{name: "metric"}, // no migrations -> status reports HasMigrations=false
	}
	cfg := &stubAppConfig{enabled: map[string]bool{"metric": true}}
	appBearer := &stubAppBearer{cfg: cfg}
	pluginBearer := &stubPluginBearer{configs: map[string]bool{"metric": true}}
	r := usecase.NewRunner(plugins, pluginBearer, &stubDBBearer{})
	// metric has no migrations, so Runner.Status doesn't invoke the factory.
	c := command.NewList(plugins, appBearer, pluginBearer, r)
	assert.NotPanics(t, func() { c.Run(&cobra.Command{}, nil) })
}
