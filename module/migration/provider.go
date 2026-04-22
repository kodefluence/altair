package migration

import (
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/module/migration/controller/command"
	"github.com/kodefluence/altair/module/migration/usecase"
)

// Provide wires the migration runner and its four CLI commands onto appModule.
// Call once from the `altair plugin` subcommand handler after plugin.LoadCommand.
func Provide(appModule module.App, registry []module.Plugin, pluginBearer core.PluginBearer, dbBearer core.DatabaseBearer) *usecase.Runner {
	runner := usecase.NewRunner(registry, pluginBearer, dbBearer)
	appModule.Controller().InjectCommand(
		command.NewMigrateUp(runner),
		command.NewMigrateDown(runner),
		command.NewMigrateStatus(runner),
		command.NewMigrateForce(runner),
	)
	return runner
}

// Runner is the exported runner type so callers (e.g. altair.go) can do drift
// detection and auto-migrate without re-constructing.
func Runner(registry []module.Plugin, pluginBearer core.PluginBearer, dbBearer core.DatabaseBearer) *usecase.Runner {
	return usecase.NewRunner(registry, pluginBearer, dbBearer)
}
