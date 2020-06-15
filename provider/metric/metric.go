package metric

import "github.com/codefluence-x/altair/core"

func Provide(appBearer core.AppBearer) {
	appBearer.SetMetricProvider(NewPrometheusMetric())
	appBearer.InjectController(NewPrometheusController())
}
