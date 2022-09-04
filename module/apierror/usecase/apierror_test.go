package usecase_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/kodefluence/altair/module/apierror/usecase"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/stretchr/testify/assert"
)

func TestApiError(t *testing.T) {

	t.Run("Internal server error", func(t *testing.T) {
		ktx := kontext.Fabricate()
		uuid := uuid.New()
		ktx.Set("request_id", uuid)

		response := jsonapi.BuildResponse(usecase.NewApiError().InternalServerError(ktx))

		assert.Equal(t, http.StatusInternalServerError, response.HTTPStatus())
		assert.Equal(
			t,
			fmt.Sprintf("JSONAPI Error:\n[Internal server error] Detail: Something is not right, help us fix this problem. Contribute to https://github.com/kodefluence/altair. Tracing code: '%v', Code: ERR0500\n", uuid),
			response.Errors.Error(),
		)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		response := jsonapi.BuildResponse(usecase.NewApiError().UnauthorizedError())

		assert.Equal(t, http.StatusUnauthorized, response.HTTPStatus())
		assert.Equal(
			t,
			"JSONAPI Error:\n[Unauthorized error] Detail: You are unauthorized, Code: ERR0401\n",
			response.Errors.Error(),
		)
	})

	t.Run("Bad request error", func(t *testing.T) {
		response := jsonapi.BuildResponse(usecase.NewApiError().BadRequestError("json"))

		assert.Equal(t, http.StatusBadRequest, response.HTTPStatus())
		assert.Equal(
			t,
			"JSONAPI Error:\n[Bad request error] Detail: You've send malformed request in your `json`, Code: ERR0400\n",
			response.Errors.Error(),
		)
	})

	t.Run("Not found error", func(t *testing.T) {
		ktx := kontext.Fabricate()
		uuid := uuid.New()
		ktx.Set("request_id", uuid)
		entityType := "oauth_applications"

		response := jsonapi.BuildResponse(usecase.NewApiError().NotFoundError(ktx, entityType))

		assert.Equal(t, http.StatusNotFound, response.HTTPStatus())
		assert.Equal(
			t,
			fmt.Sprintf("JSONAPI Error:\n[Not found error] Detail: Resource of `oauth_applications` is not found. Tracing code: `%v`, Code: ERR0404\n", ktx.GetWithoutCheck("request_id")),
			response.Errors.Error(),
		)
	})

	t.Run("Forbidden error", func(t *testing.T) {
		ktx := kontext.Fabricate()
		uuid := uuid.New()
		ktx.Set("request_id", uuid)
		entityType := "oauth_applications"
		reason := "not have access"

		response := jsonapi.BuildResponse(usecase.NewApiError().ForbiddenError(ktx, entityType, reason))

		assert.Equal(t, http.StatusForbidden, response.HTTPStatus())
		assert.Equal(
			t,
			fmt.Sprintf("JSONAPI Error:\n[Forbidden error] Detail: Resource of `oauth_applications` is forbidden to be accessed, because of: not have access. Tracing code: `%v`, Code: ERR0403\n", ktx.GetWithoutCheck("request_id")),
			response.Errors.Error(),
		)
	})

	t.Run("Validation error", func(t *testing.T) {
		response := jsonapi.BuildResponse(usecase.NewApiError().ValidationError("some validation messages goes here"))

		assert.Equal(t, http.StatusUnprocessableEntity, response.HTTPStatus())
		assert.Equal(
			t,
			"JSONAPI Error:\n[Validation error] Detail: Validation error because of: some validation messages goes here, Code: ERR1442\n",
			response.Errors.Error(),
		)
	})
}
