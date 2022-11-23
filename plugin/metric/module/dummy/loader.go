package dummy

import (
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/plugin/metric/module/dummy/controller/metric"
)

func Load(appModule module.App) {
	appModule.Controller().InjectMetric(metric.NewDummy())
}
