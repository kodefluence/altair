package pluginlist

import (
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/module/migration/usecase"
	"github.com/kodefluence/altair/module/plugin_list/controller/command"
)

// Provide wires `altair plugin list` onto appModule. Injected once from the
// `altair plugin` subcommand handler alongside the migration commands.
func Provide(appModule module.App, registry []module.Plugin, appBearer core.AppBearer, pluginBearer core.PluginBearer, runner *usecase.Runner) {
	appModule.Controller().InjectCommand(
		command.NewList(registry, appBearer, pluginBearer, runner),
	)
}
