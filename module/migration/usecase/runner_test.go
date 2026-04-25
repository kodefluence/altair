package usecase

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/golang-migrate/migrate/v4"
	"github.com/stretchr/testify/assert"

	"github.com/kodefluence/monorepo/db"

	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/module"
)

// Stubs ------------------------------------------------------------------

// stubPluginBearer is a hand-rolled core.PluginBearer that lets each test
// drive ConfigExists/DecodeConfig deterministically without gomock.
type stubPluginBearer struct {
	configExists  map[string]bool
	pluginVersion map[string]string
	decodeErr     error
	decoded       interface{}
}

func newStubPluginBearer() *stubPluginBearer {
	return &stubPluginBearer{
		configExists:  map[string]bool{},
		pluginVersion: map[string]string{},
	}
}
func (s *stubPluginBearer) ConfigExists(name string) bool { return s.configExists[name] }
func (s *stubPluginBearer) PluginVersion(name string) (string, error) {
	if v, ok := s.pluginVersion[name]; ok {
		return v, nil
	}
	return "", errors.New("not found")
}
func (s *stubPluginBearer) DecodeConfig(name string, target interface{}) error {
	s.decoded = target
	return s.decodeErr
}
func (s *stubPluginBearer) ForEach(func(string) error) {}
func (s *stubPluginBearer) Length() int                { return len(s.configExists) }

// stubDatabaseBearer maps instance name -> (db.DB, core.DatabaseConfig, error).
type stubDatabaseBearer struct {
	results map[string]stubDatabaseEntry
}

type stubDatabaseEntry struct {
	sqldb  db.DB
	config core.DatabaseConfig
	err    error
}

func (s *stubDatabaseBearer) Database(name string) (db.DB, core.DatabaseConfig, error) {
	r, ok := s.results[name]
	if !ok {
		return nil, nil, fmt.Errorf("stub: no entry for %q", name)
	}
	return r.sqldb, r.config, r.err
}

// stubMigrator records what was called and returns canned errors.
type stubMigrator struct {
	upErr        error
	stepsErr     error
	forceErr     error
	versionErr   error
	versionV     uint
	versionDirty bool

	upCalled       int
	stepsCalledN   []int
	forceCalledN   []int
	versionCalled  int
	cleanupCalled  bool
	closeForbidden bool // if set, calling cleanup is fine; calling Close on the migrator would be a fail
}

func (m *stubMigrator) Up() error { m.upCalled++; return m.upErr }
func (m *stubMigrator) Steps(n int) error {
	m.stepsCalledN = append(m.stepsCalledN, n)
	return m.stepsErr
}
func (m *stubMigrator) Force(v int) error {
	m.forceCalledN = append(m.forceCalledN, v)
	return m.forceErr
}
func (m *stubMigrator) Version() (uint, bool, error) {
	m.versionCalled++
	return m.versionV, m.versionDirty, m.versionErr
}

// stubFactory returns the supplied migrator and a cleanup that flips a flag.
func stubFactory(m *stubMigrator, factoryErr error) migratorFactory {
	return func(set module.MigrationSet, sqldb db.DB, dbConfig core.DatabaseConfig, versionTable string) (migrator, func(), error) {
		if factoryErr != nil {
			return nil, nil, factoryErr
		}
		return m, func() { m.cleanupCalled = true }, nil
	}
}

// stubPlugin is a minimal module.Plugin for table tests.
type stubPlugin struct {
	name       string
	deps       []string
	migrations []module.MigrationSet
	loadErr    error
	loadCmdErr error
}

func (p *stubPlugin) Name() string                                            { return p.name }
func (p *stubPlugin) DependsOn() []string                                     { return p.deps }
func (p *stubPlugin) Migrations(_ module.PluginContext) []module.MigrationSet { return p.migrations }
func (p *stubPlugin) Load(_ module.PluginContext) error                       { return p.loadErr }
func (p *stubPlugin) LoadCommand(_ module.PluginContext) error                { return p.loadCmdErr }
func (p *stubPlugin) SampleConfig() []byte                                    { return nil }

func makeOauthSet() module.MigrationSet {
	return module.MigrationSet{
		DatabaseInstance: "main",
		Driver:           "mysql",
		FS:               fstest.MapFS{"mysql/1_init.up.sql": {Data: []byte("")}, "mysql/1_init.down.sql": {Data: []byte("")}},
		SourcePath:       "mysql",
		VersionTable:     "oauth_plugin_db_versions",
	}
}

// Helpers ---------------------------------------------------------------

func runnerWith(t *testing.T, plugins []module.Plugin, mig *stubMigrator, factoryErr error) (*Runner, *stubPluginBearer, *stubDatabaseBearer) {
	t.Helper()
	pb := newStubPluginBearer()
	for _, p := range plugins {
		pb.configExists[p.Name()] = true
	}
	dbb := &stubDatabaseBearer{results: map[string]stubDatabaseEntry{
		"main": {sqldb: nil, config: nil, err: nil},
	}}
	r := NewRunner(plugins, pb, dbb)
	r.migratorFactory = stubFactory(mig, factoryErr)
	return r, pb, dbb
}

// findPlugin -----------------------------------------------------------

func TestFindPlugin_FoundReturnsPlugin(t *testing.T) {
	oauth := &stubPlugin{name: "oauth"}
	r, _, _ := runnerWith(t, []module.Plugin{oauth}, &stubMigrator{}, nil)
	got, ok := r.findPlugin("oauth")
	assert.True(t, ok)
	assert.Equal(t, oauth, got)
}

func TestFindPlugin_MissingReturnsFalse(t *testing.T) {
	r, _, _ := runnerWith(t, nil, &stubMigrator{}, nil)
	got, ok := r.findPlugin("ghost")
	assert.Nil(t, got)
	assert.False(t, ok)
}

// resolveVersionTable --------------------------------------------------

func TestResolveVersionTable_UsesExplicitWhenSet(t *testing.T) {
	r, _, _ := runnerWith(t, nil, &stubMigrator{}, nil)
	set := module.MigrationSet{VersionTable: "custom_versions"}
	assert.Equal(t, "custom_versions", r.resolveVersionTable("oauth", set))
}

func TestResolveVersionTable_DefaultsToPluginPrefixed(t *testing.T) {
	r, _, _ := runnerWith(t, nil, &stubMigrator{}, nil)
	set := module.MigrationSet{VersionTable: ""}
	assert.Equal(t, "oauth_plugin_db_versions", r.resolveVersionTable("oauth", set))
}

// configured -----------------------------------------------------------

func TestConfigured_ReportsBearerView(t *testing.T) {
	r, pb, _ := runnerWith(t, nil, &stubMigrator{}, nil)
	pb.configExists["oauth"] = true
	pb.configExists["metric"] = false
	assert.True(t, r.configured(&stubPlugin{name: "oauth"}))
	assert.False(t, r.configured(&stubPlugin{name: "metric"}))
	assert.False(t, r.configured(&stubPlugin{name: "ghost"}))
}

// pluginMigrationSets --------------------------------------------------

func TestPluginMigrationSets_ReturnsPluginsSets(t *testing.T) {
	set := makeOauthSet()
	oauth := &stubPlugin{name: "oauth", migrations: []module.MigrationSet{set}}
	r, _, _ := runnerWith(t, []module.Plugin{oauth}, &stubMigrator{}, nil)
	sets, err := r.pluginMigrationSets("oauth")
	assert.Nil(t, err)
	assert.Equal(t, []module.MigrationSet{set}, sets)
}

func TestPluginMigrationSets_UnknownPluginReturnsErr(t *testing.T) {
	r, _, _ := runnerWith(t, nil, &stubMigrator{}, nil)
	sets, err := r.pluginMigrationSets("ghost")
	assert.Nil(t, sets)
	assert.True(t, errors.Is(err, ErrPluginNotInRegistry), "got %v", err)
}

// runOne ---------------------------------------------------------------

func TestRunOne_EmptyDatabaseInstanceErrors(t *testing.T) {
	r, _, _ := runnerWith(t, nil, &stubMigrator{}, nil)
	err := r.runOne("oauth", module.MigrationSet{}, func(migrator) error { return nil })
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "DatabaseInstance")
}

func TestRunOne_DBBearerErrorPropagates(t *testing.T) {
	r, _, dbb := runnerWith(t, nil, &stubMigrator{}, nil)
	dbb.results["main"] = stubDatabaseEntry{err: errors.New("db down")}
	err := r.runOne("oauth", makeOauthSet(), func(migrator) error { return nil })
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "db down")
}

func TestRunOne_FactoryErrorPropagates(t *testing.T) {
	r, _, _ := runnerWith(t, nil, nil, errors.New("driver missing"))
	err := r.runOne("oauth", makeOauthSet(), func(migrator) error { return nil })
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "driver missing")
}

func TestRunOne_ActionInvokedAndCleanupAlwaysCalled(t *testing.T) {
	mig := &stubMigrator{}
	r, _, _ := runnerWith(t, nil, mig, nil)
	called := false
	err := r.runOne("oauth", makeOauthSet(), func(m migrator) error {
		called = true
		assert.Equal(t, mig, m)
		return errors.New("action error")
	})
	assert.NotNil(t, err)
	assert.True(t, called, "action must be invoked")
	assert.True(t, mig.cleanupCalled, "cleanup must run even if action errors")
}

// Up -------------------------------------------------------------------

func TestUp_HappyPathCallsMigratorUp(t *testing.T) {
	mig := &stubMigrator{}
	oauth := &stubPlugin{name: "oauth", migrations: []module.MigrationSet{makeOauthSet()}}
	r, _, _ := runnerWith(t, []module.Plugin{oauth}, mig, nil)
	assert.Nil(t, r.Up("oauth"))
	assert.Equal(t, 1, mig.upCalled)
}

func TestUp_ErrNoChangeIsTreatedAsSuccess(t *testing.T) {
	mig := &stubMigrator{upErr: migrate.ErrNoChange}
	oauth := &stubPlugin{name: "oauth", migrations: []module.MigrationSet{makeOauthSet()}}
	r, _, _ := runnerWith(t, []module.Plugin{oauth}, mig, nil)
	assert.Nil(t, r.Up("oauth"))
}

func TestUp_OtherErrorsArePropagatedWithPluginName(t *testing.T) {
	mig := &stubMigrator{upErr: errors.New("boom")}
	oauth := &stubPlugin{name: "oauth", migrations: []module.MigrationSet{makeOauthSet()}}
	r, _, _ := runnerWith(t, []module.Plugin{oauth}, mig, nil)
	err := r.Up("oauth")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "oauth")
	assert.Contains(t, err.Error(), "boom")
}

func TestUp_UnknownPluginErrors(t *testing.T) {
	r, _, _ := runnerWith(t, nil, &stubMigrator{}, nil)
	assert.NotNil(t, r.Up("ghost"))
}

func TestUp_PluginWithNoMigrationSetsIsNoop(t *testing.T) {
	mig := &stubMigrator{}
	oauth := &stubPlugin{name: "oauth", migrations: nil}
	r, _, _ := runnerWith(t, []module.Plugin{oauth}, mig, nil)
	assert.Nil(t, r.Up("oauth"))
	assert.Equal(t, 0, mig.upCalled)
}

// Down -----------------------------------------------------------------

func TestDown_StepsZeroOrLessDefaultsToOne(t *testing.T) {
	for _, n := range []int{0, -3} {
		mig := &stubMigrator{}
		oauth := &stubPlugin{name: "oauth", migrations: []module.MigrationSet{makeOauthSet()}}
		r, _, _ := runnerWith(t, []module.Plugin{oauth}, mig, nil)
		assert.Nil(t, r.Down("oauth", n))
		assert.Equal(t, []int{-1}, mig.stepsCalledN)
	}
}

func TestDown_PassesNegativeStepsToMigrator(t *testing.T) {
	mig := &stubMigrator{}
	oauth := &stubPlugin{name: "oauth", migrations: []module.MigrationSet{makeOauthSet()}}
	r, _, _ := runnerWith(t, []module.Plugin{oauth}, mig, nil)
	assert.Nil(t, r.Down("oauth", 3))
	assert.Equal(t, []int{-3}, mig.stepsCalledN)
}

func TestDown_ErrNoChangeIsSuccess(t *testing.T) {
	mig := &stubMigrator{stepsErr: migrate.ErrNoChange}
	oauth := &stubPlugin{name: "oauth", migrations: []module.MigrationSet{makeOauthSet()}}
	r, _, _ := runnerWith(t, []module.Plugin{oauth}, mig, nil)
	assert.Nil(t, r.Down("oauth", 1))
}

func TestDown_OtherErrorPropagated(t *testing.T) {
	mig := &stubMigrator{stepsErr: errors.New("rollback failed")}
	oauth := &stubPlugin{name: "oauth", migrations: []module.MigrationSet{makeOauthSet()}}
	r, _, _ := runnerWith(t, []module.Plugin{oauth}, mig, nil)
	err := r.Down("oauth", 1)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "rollback failed")
}

func TestDown_UnknownPluginErrors(t *testing.T) {
	r, _, _ := runnerWith(t, nil, &stubMigrator{}, nil)
	assert.NotNil(t, r.Down("ghost", 1))
}

// Force ----------------------------------------------------------------

func TestForce_PassesVersionThrough(t *testing.T) {
	mig := &stubMigrator{}
	oauth := &stubPlugin{name: "oauth", migrations: []module.MigrationSet{makeOauthSet()}}
	r, _, _ := runnerWith(t, []module.Plugin{oauth}, mig, nil)
	assert.Nil(t, r.Force("oauth", 7))
	assert.Equal(t, []int{7}, mig.forceCalledN)
}

func TestForce_MigratorErrorPropagated(t *testing.T) {
	mig := &stubMigrator{forceErr: errors.New("nope")}
	oauth := &stubPlugin{name: "oauth", migrations: []module.MigrationSet{makeOauthSet()}}
	r, _, _ := runnerWith(t, []module.Plugin{oauth}, mig, nil)
	assert.NotNil(t, r.Force("oauth", 7))
}

func TestForce_UnknownPluginErrors(t *testing.T) {
	r, _, _ := runnerWith(t, nil, &stubMigrator{}, nil)
	assert.NotNil(t, r.Force("ghost", 1))
}

// Status ---------------------------------------------------------------

func TestStatus_NotConfiguredPluginsAreSkipped(t *testing.T) {
	mig := &stubMigrator{versionV: 4}
	oauth := &stubPlugin{name: "oauth", migrations: []module.MigrationSet{makeOauthSet()}}
	r, pb, _ := runnerWith(t, []module.Plugin{oauth}, mig, nil)
	pb.configExists["oauth"] = false
	statuses, err := r.Status()
	assert.Nil(t, err)
	assert.Empty(t, statuses)
}

func TestStatus_PluginWithNoMigrationsReportsHasMigrationsFalse(t *testing.T) {
	metric := &stubPlugin{name: "metric", migrations: nil}
	r, _, _ := runnerWith(t, []module.Plugin{metric}, &stubMigrator{}, nil)
	statuses, err := r.Status()
	assert.Nil(t, err)
	assert.Len(t, statuses, 1)
	assert.False(t, statuses[0].HasMigrations)
	assert.Equal(t, "metric", statuses[0].Plugin)
}

func TestStatus_HappyPathReportsCurrentAndTarget(t *testing.T) {
	mig := &stubMigrator{versionV: 3}
	oauth := &stubPlugin{name: "oauth", migrations: []module.MigrationSet{makeOauthSet()}}
	r, _, _ := runnerWith(t, []module.Plugin{oauth}, mig, nil)
	statuses, err := r.Status()
	assert.Nil(t, err)
	assert.Len(t, statuses, 1)
	assert.True(t, statuses[0].HasMigrations)
	assert.Equal(t, uint(3), statuses[0].CurrentVersion)
	assert.Equal(t, uint(1), statuses[0].TargetVersion)
	assert.False(t, statuses[0].CurrentDirty)
	assert.Equal(t, "oauth_plugin_db_versions", statuses[0].VersionTable)
	assert.Equal(t, "main", statuses[0].DatabaseInstance)
	assert.Equal(t, "mysql", statuses[0].Driver)
}

func TestStatus_ErrNilVersionLeavesCurrentZero(t *testing.T) {
	mig := &stubMigrator{versionErr: migrate.ErrNilVersion}
	oauth := &stubPlugin{name: "oauth", migrations: []module.MigrationSet{makeOauthSet()}}
	r, _, _ := runnerWith(t, []module.Plugin{oauth}, mig, nil)
	statuses, err := r.Status()
	assert.Nil(t, err)
	assert.Len(t, statuses, 1)
	assert.Equal(t, uint(0), statuses[0].CurrentVersion)
	assert.False(t, statuses[0].CurrentDirty)
}

func TestStatus_VersionErrorPropagated(t *testing.T) {
	mig := &stubMigrator{versionErr: errors.New("schema lookup failed")}
	oauth := &stubPlugin{name: "oauth", migrations: []module.MigrationSet{makeOauthSet()}}
	r, _, _ := runnerWith(t, []module.Plugin{oauth}, mig, nil)
	statuses, err := r.Status()
	assert.Nil(t, statuses)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "schema lookup")
}

// Drift ----------------------------------------------------------------

func TestDrift_NoDriftWhenCurrentEqualsTarget(t *testing.T) {
	mig := &stubMigrator{versionV: 1} // target is 1 from makeOauthSet
	oauth := &stubPlugin{name: "oauth", migrations: []module.MigrationSet{makeOauthSet()}}
	r, _, _ := runnerWith(t, []module.Plugin{oauth}, mig, nil)
	reports, err := r.Drift()
	assert.Nil(t, err)
	assert.Empty(t, reports)
}

func TestDrift_ReportsWhenCurrentBelowTarget(t *testing.T) {
	mig := &stubMigrator{versionV: 0, versionErr: migrate.ErrNilVersion}
	oauth := &stubPlugin{name: "oauth", migrations: []module.MigrationSet{makeOauthSet()}}
	r, _, _ := runnerWith(t, []module.Plugin{oauth}, mig, nil)
	reports, err := r.Drift()
	assert.Nil(t, err)
	assert.Len(t, reports, 1)
	assert.Equal(t, "oauth", reports[0].Plugin)
	assert.Equal(t, uint(0), reports[0].CurrentVersion)
	assert.Equal(t, uint(1), reports[0].TargetVersion)
	assert.False(t, reports[0].Dirty)
}

func TestDrift_ReportsWhenDirtyEvenIfVersionsMatch(t *testing.T) {
	mig := &stubMigrator{versionV: 1, versionDirty: true}
	oauth := &stubPlugin{name: "oauth", migrations: []module.MigrationSet{makeOauthSet()}}
	r, _, _ := runnerWith(t, []module.Plugin{oauth}, mig, nil)
	reports, err := r.Drift()
	assert.Nil(t, err)
	assert.Len(t, reports, 1)
	assert.True(t, reports[0].Dirty)
}

func TestDrift_PluginsWithoutMigrationsAreIgnored(t *testing.T) {
	metric := &stubPlugin{name: "metric", migrations: nil}
	r, _, _ := runnerWith(t, []module.Plugin{metric}, &stubMigrator{}, nil)
	reports, err := r.Drift()
	assert.Nil(t, err)
	assert.Empty(t, reports)
}

// UpAll ---------------------------------------------------------------

func TestUpAll_RunsConfiguredPluginsInRegistryOrder(t *testing.T) {
	mig := &stubMigrator{}
	a := &stubPlugin{name: "a", migrations: []module.MigrationSet{makeOauthSet()}}
	b := &stubPlugin{name: "b", migrations: []module.MigrationSet{makeOauthSet()}}
	r, _, _ := runnerWith(t, []module.Plugin{a, b}, mig, nil)
	assert.Nil(t, r.UpAll())
	assert.Equal(t, 2, mig.upCalled, "Up called once per plugin")
}

func TestUpAll_SkipsUnconfiguredPlugins(t *testing.T) {
	mig := &stubMigrator{}
	a := &stubPlugin{name: "a", migrations: []module.MigrationSet{makeOauthSet()}}
	b := &stubPlugin{name: "b", migrations: []module.MigrationSet{makeOauthSet()}}
	r, pb, _ := runnerWith(t, []module.Plugin{a, b}, mig, nil)
	pb.configExists["b"] = false // a configured, b not
	assert.Nil(t, r.UpAll())
	assert.Equal(t, 1, mig.upCalled)
}

func TestUpAll_StopsAtFirstFailure(t *testing.T) {
	// Use distinct migrators per plugin would be cleaner; share one and
	// flip its error after the first call by capturing call count.
	mig := &stubMigrator{upErr: errors.New("boom")}
	a := &stubPlugin{name: "a", migrations: []module.MigrationSet{makeOauthSet()}}
	b := &stubPlugin{name: "b", migrations: []module.MigrationSet{makeOauthSet()}}
	r, _, _ := runnerWith(t, []module.Plugin{a, b}, mig, nil)
	err := r.UpAll()
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "a"), "error should name first failing plugin")
	assert.Equal(t, 1, mig.upCalled, "should not advance past first failure")
}

// realMigratorFactory --------------------------------------------------

func TestRealMigratorFactory_RejectsUnsupportedDriver(t *testing.T) {
	set := module.MigrationSet{Driver: "postgres", DatabaseInstance: "main"}
	m, cleanup, err := realMigratorFactory(set, nil, nil, "oauth_plugin_db_versions")
	assert.Nil(t, m)
	assert.Nil(t, cleanup)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "postgres")
	assert.Contains(t, err.Error(), "not supported")
}
