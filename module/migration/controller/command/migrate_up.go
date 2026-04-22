package command

import (
	"fmt"

	"github.com/kodefluence/altair/module/migration/usecase"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type MigrateUp struct {
	runner *usecase.Runner
	plugin string
	all    bool
}

func NewMigrateUp(runner *usecase.Runner) *MigrateUp {
	return &MigrateUp{runner: runner}
}

func (*MigrateUp) Use() string     { return "migrate:up" }
func (*MigrateUp) Short() string   { return "Apply all pending migrations for one or all plugins" }
func (*MigrateUp) Example() string { return "altair plugin migrate:up --plugin oauth | --all" }

func (m *MigrateUp) ModifyFlags(flags *pflag.FlagSet) {
	flags.StringVar(&m.plugin, "plugin", "", "Plugin name whose migrations to apply")
	flags.BoolVar(&m.all, "all", false, "Apply migrations for every plugin with migrations")
}

func (m *MigrateUp) Run(cmd *cobra.Command, args []string) {
	if m.all {
		if err := m.runner.UpAll(); err != nil {
			fmt.Fprintln(cmd.ErrOrStderr(), "error:", err)
		}
		return
	}
	if m.plugin == "" {
		fmt.Fprintln(cmd.ErrOrStderr(), "pass --plugin <name> or --all")
		return
	}
	if err := m.runner.Up(m.plugin); err != nil {
		fmt.Fprintln(cmd.ErrOrStderr(), "error:", err)
	}
}
