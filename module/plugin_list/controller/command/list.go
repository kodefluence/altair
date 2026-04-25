package command

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/module/migration/usecase"
)

// List prints every plugin compiled into this binary with its current
// activation status, dependency declarations, and (if it owns migrations)
// the embedded target version plus the DB's current version.
type List struct {
	registry     []module.Plugin
	appBearer    core.AppBearer
	pluginBearer core.PluginBearer
	runner       *usecase.Runner
}

func NewList(registry []module.Plugin, appBearer core.AppBearer, pluginBearer core.PluginBearer, runner *usecase.Runner) *List {
	return &List{registry: registry, appBearer: appBearer, pluginBearer: pluginBearer, runner: runner}
}

func (*List) Use() string { return "list" }
func (*List) Short() string {
	return "List every plugin compiled into the binary and its activation state"
}
func (*List) Example() string { return "altair plugin list" }

func (*List) ModifyFlags(flags *pflag.FlagSet) {}

func (l *List) Run(cmd *cobra.Command, args []string) {
	statusByPlugin := map[string]usecase.PluginMigrationStatus{}
	if l.runner != nil {
		if statuses, err := l.runner.Status(); err == nil {
			for _, s := range statuses {
				statusByPlugin[s.Plugin] = s
			}
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PLUGIN\tACTIVE\tDEPENDS_ON\tMIGRATIONS\tDB_INSTANCE\tCURRENT\tTARGET")
	for _, p := range l.registry {
		active := l.active(p)
		deps := strings.Join(p.DependsOn(), ",")
		if deps == "" {
			deps = "-"
		}

		migrations := "no"
		dbInstance := "-"
		current := "-"
		target := "-"
		if s, ok := statusByPlugin[p.Name()]; ok && s.HasMigrations {
			migrations = "yes"
			dbInstance = s.DatabaseInstance
			current = fmt.Sprintf("%d", s.CurrentVersion)
			target = fmt.Sprintf("%d", s.TargetVersion)
			if s.CurrentDirty {
				current += " (dirty)"
			}
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			p.Name(), activeLabel(active), deps, migrations, dbInstance, current, target)
	}
	_ = w.Flush()
}

func (l *List) active(p module.Plugin) bool {
	if l.appBearer == nil || l.appBearer.Config() == nil {
		return false
	}
	if !l.appBearer.Config().PluginExists(p.Name()) {
		return false
	}
	if l.pluginBearer == nil {
		return false
	}
	return l.pluginBearer.ConfigExists(p.Name())
}

func activeLabel(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}
