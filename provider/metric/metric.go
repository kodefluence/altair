package metric

import "github.com/kodefluence/altair/core"

func Provide(appBearer core.AppBearer) {
	appBearer.SetMetricProvider(NewPrometheusMetric())
	appBearer.InjectController(NewPrometheusController())
}
