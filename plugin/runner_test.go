package plugin

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kodefluence/monorepo/db"

	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/module"
)

type stubPlugin struct {
	name         string
	deps         []string
	loadErr      error
	loadCmdErr   error
	loadCalls    int
	loadCmdCalls int
}

func (s *stubPlugin) Name() string                                          { return s.name }
func (s *stubPlugin) DependsOn() []string                                   { return s.deps }
func (s *stubPlugin) Migrations(module.PluginContext) []module.MigrationSet { return nil }
func (s *stubPlugin) Load(module.PluginContext) error {
	s.loadCalls++
	return s.loadErr
}
func (s *stubPlugin) LoadCommand(module.PluginContext) error {
	s.loadCmdCalls++
	return s.loadCmdErr
}
func (s *stubPlugin) SampleConfig() []byte { return nil }

// stub bearers --------------------------------------------------------

type stubAppConfig struct {
	enabled map[string]bool
}

func (s *stubAppConfig) Port() int                     { return 0 }
func (s *stubAppConfig) BasicAuthUsername() string     { return "" }
func (s *stubAppConfig) BasicAuthPassword() string     { return "" }
func (s *stubAppConfig) ProxyHost() string             { return "" }
func (s *stubAppConfig) PluginExists(name string) bool { return s.enabled[name] }
func (s *stubAppConfig) Plugins() []string             { return nil }
func (s *stubAppConfig) AutoMigrate() bool             { return false }
func (s *stubAppConfig) Dump() string                  { return "" }

type stubAppBearer struct{ cfg core.AppConfig }

func (s *stubAppBearer) Config() core.AppConfig                         { return s.cfg }
func (s *stubAppBearer) DownStreamPlugins() []core.DownStreamPlugin     { return nil }
func (s *stubAppBearer) InjectDownStreamPlugin(_ core.DownStreamPlugin) {}
func (s *stubAppBearer) SetMetricProvider(_ core.Metric)                {}
func (s *stubAppBearer) MetricProvider() (core.Metric, error)           { return nil, nil }

type stubPluginBearer struct {
	configs  map[string]bool
	versions map[string]string
}

func (s *stubPluginBearer) ConfigExists(name string) bool { return s.configs[name] }
func (s *stubPluginBearer) PluginVersion(name string) (string, error) {
	if v, ok := s.versions[name]; ok {
		return v, nil
	}
	return "", errors.New("not found")
}
func (s *stubPluginBearer) DecodeConfig(string, interface{}) error { return nil }
func (s *stubPluginBearer) ForEach(func(string) error)             {}
func (s *stubPluginBearer) Length() int                            { return len(s.configs) }

type stubDBBearer struct{}

func (s *stubDBBearer) Database(string) (db.DB, core.DatabaseConfig, error) { return nil, nil, nil }

func runWith(_ []module.Plugin, enabled, configured map[string]bool, versions map[string]string) (core.AppBearer, core.PluginBearer, core.DatabaseBearer) {
	return &stubAppBearer{cfg: &stubAppConfig{enabled: enabled}},
		&stubPluginBearer{configs: configured, versions: versions},
		&stubDBBearer{}
}

func names(plugins []module.Plugin) []string {
	out := make([]string, len(plugins))
	for i, p := range plugins {
		out[i] = p.Name()
	}
	return out
}

func TestTopoSort_AlphabeticalTieBreak(t *testing.T) {
	// No dependencies; alphabetical ordering keeps `make test` diffs stable.
	plugins := []module.Plugin{
		&stubPlugin{name: "zeta"},
		&stubPlugin{name: "alpha"},
		&stubPlugin{name: "mu"},
	}

	ordered, err := topoSort(plugins)
	assert.Nil(t, err)
	assert.Equal(t, []string{"alpha", "mu", "zeta"}, names(ordered))
}

func TestTopoSort_RespectsDependsOn(t *testing.T) {
	plugins := []module.Plugin{
		&stubPlugin{name: "oauth", deps: []string{"metric"}},
		&stubPlugin{name: "metric"},
	}

	ordered, err := topoSort(plugins)
	assert.Nil(t, err)
	assert.Equal(t, []string{"metric", "oauth"}, names(ordered))
}

func TestTopoSort_SkipsSoftMissingDependency(t *testing.T) {
	// oauth depends on metric, but metric is not in the active set;
	// soft dep semantics say we load oauth anyway without error.
	plugins := []module.Plugin{
		&stubPlugin{name: "oauth", deps: []string{"metric"}},
	}

	ordered, err := topoSort(plugins)
	assert.Nil(t, err)
	assert.Equal(t, []string{"oauth"}, names(ordered))
}

func TestTopoSort_RejectsCycles(t *testing.T) {
	plugins := []module.Plugin{
		&stubPlugin{name: "a", deps: []string{"b"}},
		&stubPlugin{name: "b", deps: []string{"a"}},
	}

	ordered, err := topoSort(plugins)
	assert.Nil(t, ordered)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "cycle"), "expected cycle error, got: %v", err)
	assert.True(t, strings.Contains(err.Error(), "a"), "expected residual 'a' in error")
	assert.True(t, strings.Contains(err.Error(), "b"), "expected residual 'b' in error")
}

func TestTopoSort_RejectsSelfDependency(t *testing.T) {
	plugins := []module.Plugin{
		&stubPlugin{name: "loner", deps: []string{"loner"}},
	}

	ordered, err := topoSort(plugins)
	assert.Nil(t, ordered)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "loner"), "expected plugin name in self-dep error")
}

// activePlugins -------------------------------------------------------

// Assumption: a plugin missing from app.yml `plugins:` is silently skipped.
func TestActivePlugins_SkipsPluginsAbsentFromAppYaml(t *testing.T) {
	plugins := []module.Plugin{&stubPlugin{name: "oauth"}}
	appBearer, pluginBearer, _ := runWith(plugins, map[string]bool{}, map[string]bool{"oauth": true}, nil)
	got, err := activePlugins(plugins, appBearer, pluginBearer)
	assert.Nil(t, err)
	assert.Empty(t, got)
}

// Assumption: a plugin in app.yml but missing from config/plugin/ is a hard
// error with a message naming the plugin.
func TestActivePlugins_ErrorsWhenAppEnabledButConfigMissing(t *testing.T) {
	plugins := []module.Plugin{&stubPlugin{name: "oauth"}}
	appBearer, pluginBearer, _ := runWith(plugins, map[string]bool{"oauth": true}, map[string]bool{}, nil)
	got, err := activePlugins(plugins, appBearer, pluginBearer)
	assert.Nil(t, got)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "oauth")
}

// Assumption: only AND-gated plugins are returned, in registry order.
func TestActivePlugins_ReturnsAndGatedSubsetInOrder(t *testing.T) {
	a := &stubPlugin{name: "a"}
	b := &stubPlugin{name: "b"}
	c := &stubPlugin{name: "c"}
	plugins := []module.Plugin{a, b, c}
	enabled := map[string]bool{"a": true, "c": true}
	configured := map[string]bool{"a": true, "c": true}
	appBearer, pluginBearer, _ := runWith(plugins, enabled, configured, nil)
	got, err := activePlugins(plugins, appBearer, pluginBearer)
	assert.Nil(t, err)
	assert.Equal(t, []string{"a", "c"}, names(got))
}

// run() — the internal orchestrator -----------------------------------
//
// Note: tests target run() directly rather than the public Load/LoadCommand
// because those use the package-level Registry() and would invoke real
// plugins. run() takes the registry as a parameter, which is the testable
// seam.

// Assumption: run invokes the action callback on every active plugin in
// topo-sorted order.
func TestRun_InvokesEveryActivePlugin(t *testing.T) {
	a := &stubPlugin{name: "a"}
	b := &stubPlugin{name: "b"}
	plugins := []module.Plugin{a, b}
	enabled := map[string]bool{"a": true, "b": true}
	configured := map[string]bool{"a": true, "b": true}
	versions := map[string]string{"a": "1.0", "b": "1.0"}
	appBearer, pluginBearer, dbBearer := runWith(plugins, enabled, configured, versions)
	err := run(plugins, appBearer, pluginBearer, dbBearer, nil, nil, func(p module.Plugin, ctx module.PluginContext) error {
		return p.Load(ctx)
	})
	assert.Nil(t, err)
	assert.Equal(t, 1, a.loadCalls)
	assert.Equal(t, 1, b.loadCalls)
}

// Assumption: run surfaces a plugin's error wrapped with the plugin's name.
func TestRun_PluginErrorsAreWrappedWithPluginName(t *testing.T) {
	a := &stubPlugin{name: "a", loadErr: errors.New("init blew up")}
	plugins := []module.Plugin{a}
	enabled := map[string]bool{"a": true}
	configured := map[string]bool{"a": true}
	versions := map[string]string{"a": "1.0"}
	appBearer, pluginBearer, dbBearer := runWith(plugins, enabled, configured, versions)
	err := run(plugins, appBearer, pluginBearer, dbBearer, nil, nil, func(p module.Plugin, ctx module.PluginContext) error {
		return p.Load(ctx)
	})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "a")
	assert.Contains(t, err.Error(), "init blew up")
}

// Assumption: run is a no-op when nothing is active.
func TestRun_NoActivePluginsIsNoop(t *testing.T) {
	a := &stubPlugin{name: "a"}
	plugins := []module.Plugin{a}
	appBearer, pluginBearer, dbBearer := runWith(plugins, map[string]bool{}, map[string]bool{}, nil)
	err := run(plugins, appBearer, pluginBearer, dbBearer, nil, nil, func(p module.Plugin, ctx module.PluginContext) error {
		return p.Load(ctx)
	})
	assert.Nil(t, err)
	assert.Equal(t, 0, a.loadCalls)
}

// Assumption: missing PluginVersion entry surfaces wrapped with plugin name.
func TestRun_MissingPluginVersionErrors(t *testing.T) {
	a := &stubPlugin{name: "a"}
	plugins := []module.Plugin{a}
	enabled := map[string]bool{"a": true}
	configured := map[string]bool{"a": true}
	appBearer, pluginBearer, dbBearer := runWith(plugins, enabled, configured, nil)
	err := run(plugins, appBearer, pluginBearer, dbBearer, nil, nil, func(p module.Plugin, ctx module.PluginContext) error {
		return p.Load(ctx)
	})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "a")
}

// Assumption: AND-gate violation (plugin in app.yml, no config file)
// errors before any plugin Load runs.
func TestRun_AndGateViolationErrorsBeforeLoad(t *testing.T) {
	a := &stubPlugin{name: "a"}
	plugins := []module.Plugin{a}
	enabled := map[string]bool{"a": true}
	configured := map[string]bool{} // missing
	appBearer, pluginBearer, dbBearer := runWith(plugins, enabled, configured, nil)
	err := run(plugins, appBearer, pluginBearer, dbBearer, nil, nil, func(p module.Plugin, ctx module.PluginContext) error {
		return p.Load(ctx)
	})
	assert.NotNil(t, err)
	assert.Equal(t, 0, a.loadCalls, "AND-gate failure must short-circuit before Load runs")
}

// Assumption: the public Load entry compiles and is wired to Registry +
// run; the smoke test exercises real plugins. This tests just that calling
// Load with empty bearers returns nil (no real plugins active in test).
func TestPublicLoad_SmokeWithEmptyBearersReturnsNil(t *testing.T) {
	appBearer, pluginBearer, dbBearer := runWith(nil, map[string]bool{}, map[string]bool{}, nil)
	assert.Nil(t, Load(appBearer, pluginBearer, dbBearer, nil, nil))
	assert.Nil(t, LoadCommand(appBearer, pluginBearer, dbBearer, nil, nil))
}

// buildContext --------------------------------------------------------

// Assumption: buildContext returns a PluginContext with non-nil DecodeConfig
// and Database closures, the supplied version, and pre-tagged logger.
func TestBuildContext_PopulatesAllClosures(t *testing.T) {
	pb := &stubPluginBearer{configs: map[string]bool{"oauth": true}, versions: map[string]string{"oauth": "1.0"}}
	dbb := &stubDBBearer{}
	ab := &stubAppBearer{cfg: &stubAppConfig{}}
	p := &stubPlugin{name: "oauth"}
	ctx := buildContext(p, "1.0", pb, dbb, ab, nil, nil)
	assert.Equal(t, "1.0", ctx.Version)
	assert.NotNil(t, ctx.DecodeConfig)
	assert.NotNil(t, ctx.Database)
	assert.Equal(t, ab, ctx.AppBearer)
}
