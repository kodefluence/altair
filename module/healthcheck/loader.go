package healthcheck

import (
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/module/healthcheck/controller/http"
)

func Load(app module.App) {
	app.Controller().InjectHTTP(http.NewHealthController())
}
