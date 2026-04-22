package command

import (
	"fmt"

	"github.com/kodefluence/altair/module/migration/usecase"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// MigrateForce is the escape hatch for a dirty migrations table: sets the
// schema version pointer without running SQL. Use after manually reconciling
// a half-applied migration.
type MigrateForce struct {
	runner  *usecase.Runner
	plugin  string
	version int
}

func NewMigrateForce(runner *usecase.Runner) *MigrateForce {
	return &MigrateForce{runner: runner, version: -1}
}

func (*MigrateForce) Use() string     { return "migrate:force" }
func (*MigrateForce) Short() string   { return "Force-set a plugin's migration version (clears dirty state)" }
func (*MigrateForce) Example() string { return "altair plugin migrate:force --plugin oauth --version 3" }

func (m *MigrateForce) ModifyFlags(flags *pflag.FlagSet) {
	flags.StringVar(&m.plugin, "plugin", "", "Plugin whose schema_version to force")
	flags.IntVar(&m.version, "version", -1, "Target version to set")
}

func (m *MigrateForce) Run(cmd *cobra.Command, args []string) {
	if m.plugin == "" {
		fmt.Fprintln(cmd.ErrOrStderr(), "pass --plugin <name>")
		return
	}
	if m.version < 0 {
		fmt.Fprintln(cmd.ErrOrStderr(), "pass --version <n>")
		return
	}
	if err := m.runner.Force(m.plugin, m.version); err != nil {
		fmt.Fprintln(cmd.ErrOrStderr(), "error:", err)
	}
}
