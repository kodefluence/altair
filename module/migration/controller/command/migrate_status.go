package command

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/kodefluence/altair/module/migration/usecase"
)

type MigrateStatus struct {
	runner *usecase.Runner
}

func NewMigrateStatus(runner *usecase.Runner) *MigrateStatus {
	return &MigrateStatus{runner: runner}
}

func (*MigrateStatus) Use() string     { return "migrate:status" }
func (*MigrateStatus) Short() string   { return "Show current vs target migration versions per plugin" }
func (*MigrateStatus) Example() string { return "altair plugin migrate:status" }

func (*MigrateStatus) ModifyFlags(flags *pflag.FlagSet) {}

func (m *MigrateStatus) Run(cmd *cobra.Command, args []string) {
	statuses, err := m.runner.Status()
	if err != nil {
		fmt.Fprintln(cmd.ErrOrStderr(), "error:", err)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PLUGIN\tDB\tDRIVER\tVERSION_TABLE\tCURRENT\tTARGET\tDRIFT\tDIRTY")
	for _, s := range statuses {
		if !s.HasMigrations {
			fmt.Fprintf(w, "%s\t-\t-\t-\tn/a\tn/a\tno-migrations\tno\n", s.Plugin)
			continue
		}
		drift := "up-to-date"
		if s.CurrentVersion < s.TargetVersion {
			drift = "behind"
		}
		dirty := "no"
		if s.CurrentDirty {
			dirty = "yes"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%d\t%s\t%s\n",
			s.Plugin, s.DatabaseInstance, s.Driver, s.VersionTable,
			s.CurrentVersion, s.TargetVersion, drift, dirty)
	}
	_ = w.Flush()
}
