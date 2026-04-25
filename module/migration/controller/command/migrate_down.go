package command

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/kodefluence/altair/module/migration/usecase"
)

type MigrateDown struct {
	runner *usecase.Runner
	plugin string
	steps  int
}

func NewMigrateDown(runner *usecase.Runner) *MigrateDown {
	return &MigrateDown{runner: runner, steps: 1}
}

func (*MigrateDown) Use() string     { return "migrate:down" }
func (*MigrateDown) Short() string   { return "Roll back the last N migrations for a plugin" }
func (*MigrateDown) Example() string { return "altair plugin migrate:down --plugin oauth --steps 1" }

func (m *MigrateDown) ModifyFlags(flags *pflag.FlagSet) {
	flags.StringVar(&m.plugin, "plugin", "", "Plugin name whose migrations to roll back")
	flags.IntVar(&m.steps, "steps", 1, "Number of migrations to roll back")
}

func (m *MigrateDown) Run(cmd *cobra.Command, args []string) {
	if m.plugin == "" {
		fmt.Fprintln(cmd.ErrOrStderr(), "pass --plugin <name>")
		return
	}
	if err := m.runner.Down(m.plugin, m.steps); err != nil {
		fmt.Fprintln(cmd.ErrOrStderr(), "error:", err)
	}
}
