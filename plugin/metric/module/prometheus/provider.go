package prometheus

import (
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/plugin/metric/module/prometheus/controller/http"
	"github.com/kodefluence/altair/plugin/metric/module/prometheus/usecase"
)

func Provide(appBearer core.AppBearer) {
	appBearer.SetMetricProvider(usecase.NewPrometheus())
	appBearer.InjectController(http.NewPrometheusController())
}
