package formatter

import (
	"time"

	"github.com/kodefluence/altair/plugin/oauth/module/formatter/usecase"
)

func Provide(tokenExpiresIn time.Duration, codeExpiresIn time.Duration, refreshTokenExpiresIn time.Duration) *usecase.Formatter {
	return usecase.NewFormatter(tokenExpiresIn, codeExpiresIn, refreshTokenExpiresIn)
}
