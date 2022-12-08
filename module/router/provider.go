package router

import (
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/module/router/usecase"
)

func Provide(downStreamPlugin []module.DownstreamController, metric []module.MetricController) (*usecase.Compiler, *usecase.Generator) {
	return usecase.NewCompiler(), usecase.NewGenerator(downStreamPlugin, metric)
}
