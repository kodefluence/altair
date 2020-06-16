package provider

import (
	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/provider/metric"
)

func Metric(appBearer core.AppBearer) {
	metric.Provide(appBearer)
}
