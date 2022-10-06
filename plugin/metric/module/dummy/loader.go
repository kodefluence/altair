package dummy

import (
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/plugin/metric/module/dummy/controller/metric"
)

func Load(appBearer core.AppBearer) {
	appBearer.SetMetricProvider(metric.NewDummy())
}
