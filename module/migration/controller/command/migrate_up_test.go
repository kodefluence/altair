package command_test

import (
	"bytes"
	"errors"
	"testing"
	"testing/fstest"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"

	"github.com/kodefluence/monorepo/db"

	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/module/migration/controller/command"
	"github.com/kodefluence/altair/module/migration/usecase"
)

// stubBearers + stubMigrator from runner_test mirror the structure here. We
// avoid importing them by re-declaring the minimum needed: a Runner backed
// by real stub plugin/database bearers + a factory that returns a fake migrator.

type stubBearer struct {
	configExists  map[string]bool
	pluginVersion map[string]string
}

func (s *stubBearer) ConfigExists(name string) bool { return s.configExists[name] }
func (s *stubBearer) PluginVersion(name string) (string, error) {
	if v, ok := s.pluginVersion[name]; ok {
		return v, nil
	}
	return "", errors.New("not found")
}
func (s *stubBearer) DecodeConfig(string, interface{}) error { return nil }
func (s *stubBearer) ForEach(func(string) error)             {}
func (s *stubBearer) Length() int                            { return len(s.configExists) }

type stubDBBearer struct{}

func (s *stubDBBearer) Database(string) (db.DB, core.DatabaseConfig, error) { return nil, nil, nil }

type stubPlugin struct {
	name string
	sets []module.MigrationSet
}

func (p *stubPlugin) Name() string                                            { return p.name }
func (p *stubPlugin) DependsOn() []string                                     { return nil }
func (p *stubPlugin) Migrations(_ module.PluginContext) []module.MigrationSet { return p.sets }
func (p *stubPlugin) Load(_ module.PluginContext) error                       { return nil }
func (p *stubPlugin) LoadCommand(_ module.PluginContext) error                { return nil }
func (p *stubPlugin) SampleConfig() []byte                                    { return nil }

// makeRunner produces a Runner whose factory always returns the supplied
// fakeMigrator and a no-op cleanup. The plugin bearer reports the named
// plugin as configured.
func makeRunner(t *testing.T, pluginName string, fake *fakeMigrator, factoryErr error) *usecase.Runner {
	t.Helper()
	set := module.MigrationSet{
		DatabaseInstance: "main",
		Driver:           "mysql",
		FS:               fstest.MapFS{"mysql/1_x.up.sql": {Data: []byte("")}, "mysql/1_x.down.sql": {Data: []byte("")}},
		SourcePath:       "mysql",
		VersionTable:     pluginName + "_plugin_db_versions",
	}
	plugins := []module.Plugin{&stubPlugin{name: pluginName, sets: []module.MigrationSet{set}}}
	pb := &stubBearer{configExists: map[string]bool{pluginName: true}}
	dbb := &stubDBBearer{}
	r := usecase.NewRunner(plugins, pb, dbb)
	usecase.SetMigratorFactoryForTest(r, fake.factory(factoryErr))
	return r
}

// fakeMigrator records calls; satisfies the unexported migrator interface
// via the SetMigratorFactoryForTest test seam (declared in runner.go).
type fakeMigrator struct {
	upErr      error
	stepsErr   error
	forceErr   error
	versionErr error
	versionV   uint
	dirty      bool

	upCalls    int
	stepsArgs  []int
	forceArgs  []int
	versionGot int
}

func (f *fakeMigrator) Up() error         { f.upCalls++; return f.upErr }
func (f *fakeMigrator) Steps(n int) error { f.stepsArgs = append(f.stepsArgs, n); return f.stepsErr }
func (f *fakeMigrator) Force(v int) error { f.forceArgs = append(f.forceArgs, v); return f.forceErr }
func (f *fakeMigrator) Version() (uint, bool, error) {
	f.versionGot++
	return f.versionV, f.dirty, f.versionErr
}

func (f *fakeMigrator) factory(factoryErr error) usecase.MigratorFactoryForTest {
	return func(_ module.MigrationSet, _ db.DB, _ core.DatabaseConfig, _ string) (usecase.MigratorForTest, func(), error) {
		if factoryErr != nil {
			return nil, nil, factoryErr
		}
		return f, func() {}, nil
	}
}

// MigrateUp -----------------------------------------------------------

func TestMigrateUp_GettersReturnExpectedStrings(t *testing.T) {
	c := command.NewMigrateUp(nil)
	assert.Equal(t, "migrate:up", c.Use())
	assert.NotEmpty(t, c.Short())
	assert.NotEmpty(t, c.Example())
}

func TestMigrateUp_ModifyFlagsRegistersFlags(t *testing.T) {
	c := command.NewMigrateUp(nil)
	cobraCmd := &cobra.Command{}
	c.ModifyFlags(cobraCmd.Flags())
	assert.NotNil(t, cobraCmd.Flags().Lookup("plugin"))
	assert.NotNil(t, cobraCmd.Flags().Lookup("all"))
}

func TestMigrateUp_RunWithoutFlagsErrorsToStderr(t *testing.T) {
	c := command.NewMigrateUp(nil)
	stderr := &bytes.Buffer{}
	cobraCmd := &cobra.Command{}
	cobraCmd.SetErr(stderr)
	c.ModifyFlags(cobraCmd.Flags())
	c.Run(cobraCmd, nil)
	assert.Contains(t, stderr.String(), "--plugin")
}

func TestMigrateUp_RunWithPluginFlagInvokesUp(t *testing.T) {
	fake := &fakeMigrator{}
	r := makeRunner(t, "oauth", fake, nil)
	c := command.NewMigrateUp(r)
	cobraCmd := &cobra.Command{}
	c.ModifyFlags(cobraCmd.Flags())
	_ = cobraCmd.Flags().Set("plugin", "oauth")
	c.Run(cobraCmd, nil)
	assert.Equal(t, 1, fake.upCalls)
}

func TestMigrateUp_RunWithAllFlagInvokesUpAll(t *testing.T) {
	fake := &fakeMigrator{}
	r := makeRunner(t, "oauth", fake, nil)
	c := command.NewMigrateUp(r)
	cobraCmd := &cobra.Command{}
	c.ModifyFlags(cobraCmd.Flags())
	_ = cobraCmd.Flags().Set("all", "true")
	c.Run(cobraCmd, nil)
	assert.Equal(t, 1, fake.upCalls)
}

func TestMigrateUp_RunReportsRunnerError(t *testing.T) {
	fake := &fakeMigrator{upErr: errors.New("boom")}
	r := makeRunner(t, "oauth", fake, nil)
	c := command.NewMigrateUp(r)
	stderr := &bytes.Buffer{}
	cobraCmd := &cobra.Command{}
	cobraCmd.SetErr(stderr)
	c.ModifyFlags(cobraCmd.Flags())
	_ = cobraCmd.Flags().Set("plugin", "oauth")
	c.Run(cobraCmd, nil)
	assert.Contains(t, stderr.String(), "boom")
}

// MigrateDown ---------------------------------------------------------

func TestMigrateDown_Getters(t *testing.T) {
	c := command.NewMigrateDown(nil)
	assert.Equal(t, "migrate:down", c.Use())
	assert.NotEmpty(t, c.Short())
	assert.NotEmpty(t, c.Example())
}

func TestMigrateDown_ModifyFlagsRegistersFlags(t *testing.T) {
	c := command.NewMigrateDown(nil)
	cobraCmd := &cobra.Command{}
	c.ModifyFlags(cobraCmd.Flags())
	assert.NotNil(t, cobraCmd.Flags().Lookup("plugin"))
	assert.NotNil(t, cobraCmd.Flags().Lookup("steps"))
}

func TestMigrateDown_RunWithoutPluginErrors(t *testing.T) {
	c := command.NewMigrateDown(nil)
	stderr := &bytes.Buffer{}
	cobraCmd := &cobra.Command{}
	cobraCmd.SetErr(stderr)
	c.ModifyFlags(cobraCmd.Flags())
	c.Run(cobraCmd, nil)
	assert.Contains(t, stderr.String(), "--plugin")
}

func TestMigrateDown_RunPassesStepsThrough(t *testing.T) {
	fake := &fakeMigrator{}
	r := makeRunner(t, "oauth", fake, nil)
	c := command.NewMigrateDown(r)
	cobraCmd := &cobra.Command{}
	c.ModifyFlags(cobraCmd.Flags())
	_ = cobraCmd.Flags().Set("plugin", "oauth")
	_ = cobraCmd.Flags().Set("steps", "3")
	c.Run(cobraCmd, nil)
	assert.Equal(t, []int{-3}, fake.stepsArgs)
}

func TestMigrateDown_RunReportsRunnerError(t *testing.T) {
	fake := &fakeMigrator{stepsErr: errors.New("rollback failed")}
	r := makeRunner(t, "oauth", fake, nil)
	c := command.NewMigrateDown(r)
	stderr := &bytes.Buffer{}
	cobraCmd := &cobra.Command{}
	cobraCmd.SetErr(stderr)
	c.ModifyFlags(cobraCmd.Flags())
	_ = cobraCmd.Flags().Set("plugin", "oauth")
	c.Run(cobraCmd, nil)
	assert.Contains(t, stderr.String(), "rollback failed")
}

// MigrateStatus -------------------------------------------------------

func TestMigrateStatus_Getters(t *testing.T) {
	c := command.NewMigrateStatus(nil)
	assert.Equal(t, "migrate:status", c.Use())
	assert.NotEmpty(t, c.Short())
	assert.NotEmpty(t, c.Example())
}

func TestMigrateStatus_ModifyFlagsRegistersNoFlags(t *testing.T) {
	c := command.NewMigrateStatus(nil)
	cobraCmd := &cobra.Command{}
	c.ModifyFlags(cobraCmd.Flags())
	count := 0
	cobraCmd.Flags().VisitAll(func(*pflag.Flag) { count++ })
	assert.Equal(t, 0, count, "migrate:status takes no flags")
}

func TestMigrateStatus_RunPrintsTableAndReturns(t *testing.T) {
	fake := &fakeMigrator{versionV: 5}
	r := makeRunner(t, "oauth", fake, nil)
	c := command.NewMigrateStatus(r)
	cobraCmd := &cobra.Command{}
	c.ModifyFlags(cobraCmd.Flags())
	// Status writes to os.Stdout (tabwriter); we just confirm Run completes
	// without error and that Version() got queried via the runner.
	assert.NotPanics(t, func() { c.Run(cobraCmd, nil) })
	assert.GreaterOrEqual(t, fake.versionGot, 1)
}

func TestMigrateStatus_RunPrintsErrorOnFailure(t *testing.T) {
	fake := &fakeMigrator{versionErr: errors.New("boom")}
	r := makeRunner(t, "oauth", fake, nil)
	c := command.NewMigrateStatus(r)
	stderr := &bytes.Buffer{}
	cobraCmd := &cobra.Command{}
	cobraCmd.SetErr(stderr)
	c.ModifyFlags(cobraCmd.Flags())
	c.Run(cobraCmd, nil)
	assert.Contains(t, stderr.String(), "boom")
}

// MigrateForce --------------------------------------------------------

func TestMigrateForce_Getters(t *testing.T) {
	c := command.NewMigrateForce(nil)
	assert.Equal(t, "migrate:force", c.Use())
	assert.NotEmpty(t, c.Short())
	assert.NotEmpty(t, c.Example())
}

func TestMigrateForce_ModifyFlagsRegistersFlags(t *testing.T) {
	c := command.NewMigrateForce(nil)
	cobraCmd := &cobra.Command{}
	c.ModifyFlags(cobraCmd.Flags())
	assert.NotNil(t, cobraCmd.Flags().Lookup("plugin"))
	assert.NotNil(t, cobraCmd.Flags().Lookup("version"))
}

func TestMigrateForce_RunWithoutPluginErrors(t *testing.T) {
	c := command.NewMigrateForce(nil)
	stderr := &bytes.Buffer{}
	cobraCmd := &cobra.Command{}
	cobraCmd.SetErr(stderr)
	c.ModifyFlags(cobraCmd.Flags())
	c.Run(cobraCmd, nil)
	assert.Contains(t, stderr.String(), "--plugin")
}

func TestMigrateForce_RunWithoutVersionErrors(t *testing.T) {
	c := command.NewMigrateForce(nil)
	stderr := &bytes.Buffer{}
	cobraCmd := &cobra.Command{}
	cobraCmd.SetErr(stderr)
	c.ModifyFlags(cobraCmd.Flags())
	_ = cobraCmd.Flags().Set("plugin", "oauth")
	c.Run(cobraCmd, nil)
	assert.Contains(t, stderr.String(), "--version")
}

func TestMigrateForce_RunPassesVersionThrough(t *testing.T) {
	fake := &fakeMigrator{}
	r := makeRunner(t, "oauth", fake, nil)
	c := command.NewMigrateForce(r)
	cobraCmd := &cobra.Command{}
	c.ModifyFlags(cobraCmd.Flags())
	_ = cobraCmd.Flags().Set("plugin", "oauth")
	_ = cobraCmd.Flags().Set("version", "9")
	c.Run(cobraCmd, nil)
	assert.Equal(t, []int{9}, fake.forceArgs)
}

func TestMigrateForce_RunReportsRunnerError(t *testing.T) {
	fake := &fakeMigrator{forceErr: errors.New("nope")}
	r := makeRunner(t, "oauth", fake, nil)
	c := command.NewMigrateForce(r)
	stderr := &bytes.Buffer{}
	cobraCmd := &cobra.Command{}
	cobraCmd.SetErr(stderr)
	c.ModifyFlags(cobraCmd.Flags())
	_ = cobraCmd.Flags().Set("plugin", "oauth")
	_ = cobraCmd.Flags().Set("version", "1")
	c.Run(cobraCmd, nil)
	assert.Contains(t, stderr.String(), "nope")
}
