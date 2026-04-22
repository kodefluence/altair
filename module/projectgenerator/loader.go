package projectgenerator

import (
	"embed"

	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/module/projectgenerator/controller/command"
)

//go:embed template
var fs embed.FS

// Load wires the `altair new` command with the plugin registry so the
// generated config/plugin/*.yml files always match the plugins compiled into
// the binary.
func Load(app module.App, plugins []module.Plugin) {
	app.Controller().InjectCommand(command.NewNew(fs, plugins))
}
