package prometheus

import (
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/plugin/metric/module/prometheus/controller/http"
	"github.com/kodefluence/altair/plugin/metric/module/prometheus/controller/metric"
)

func Load(appBearer core.AppBearer) {
	appBearer.SetMetricProvider(metric.NewPrometheus())
	appBearer.InjectController(http.NewPrometheusController())
}
