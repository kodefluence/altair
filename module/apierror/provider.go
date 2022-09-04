package apierror

import "github.com/kodefluence/altair/module/apierror/usecase"

func Provide() *usecase.ApiError {
	return usecase.NewApiError()
}
