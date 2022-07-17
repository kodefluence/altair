package provider

import (
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/provider/metric"
)

func Metric(appBearer core.AppBearer) {
	metric.Provide(appBearer)
}
