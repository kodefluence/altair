package prometheus

import (
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/plugin/metric/module/prometheus/controller/http"
	"github.com/kodefluence/altair/plugin/metric/module/prometheus/controller/metric"
)

func Load(appModule module.App) {
	appModule.Controller().InjectMetric(metric.NewPrometheus())
	appModule.Controller().InjectHTTP(http.NewPrometheusController())
}
