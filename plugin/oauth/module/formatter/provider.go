package formatter

import "github.com/kodefluence/altair/plugin/oauth/module/formatter/usecase"

func Provide() *usecase.Formatter {
	return usecase.NewFormatter()
}
