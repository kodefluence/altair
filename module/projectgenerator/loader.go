package projectgenerator

import (
	"embed"

	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/module/projectgenerator/controller/command"
)

//go:embed template
var fs embed.FS

func Load(app module.App) {
	app.Controller().InjectCommand(command.NewNew(fs))
}
