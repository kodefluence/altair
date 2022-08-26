package dummy

import (
	"github.com/kodefluence/altair/core"
	"github.com/kodefluence/altair/plugin/metric/module/dummy/usecase"
)

func Provide(appBearer core.AppBearer) {
	appBearer.SetMetricProvider(usecase.NewDummy())
}
