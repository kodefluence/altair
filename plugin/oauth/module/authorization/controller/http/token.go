package http

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// TokenController control flow of creating access token
type TokenController struct {
	authorizationUsecase Authorization
	apiError             module.ApiError
}

// NewToken create new token controller
func NewToken(authorizationUsecase Authorization, apiError module.ApiError) *TokenController {
	return &TokenController{
		authorizationUsecase: authorizationUsecase,
		apiError:             apiError,
	}
}

// Method POST
func (o *TokenController) Method() string {
	return "POST"
}

// Path /oauth/authorizations/token
func (o *TokenController) Path() string {
	return "/oauth/authorizations/token"
}

// Control creating access token based on access token request
func (o *TokenController) Control(ktx kontext.Context, c *gin.Context) {
	var req entity.AccessTokenRequestJSON

	rawData, err := c.GetRawData()
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", c.Value("request_id")).
			Array("tags", zerolog.Arr().Str("controller").Str("authorization").Str("token").Str("get_raw_data")).
			Msg("Cannot get raw data")

		c.JSON(http.StatusBadRequest, jsonapi.BuildResponse(o.apiError.BadRequestError("request body")))
		return
	}

	err = json.Unmarshal(rawData, &req)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", c.Value("request_id")).
			Array("tags", zerolog.Arr().Str("controller").Str("authorization").Str("token").Str("unmarshal")).
			Msg("Cannot unmarshal json")

		c.JSON(http.StatusBadRequest, jsonapi.BuildResponse(o.apiError.BadRequestError("request body")))
		return
	}

	data, jsonapierr := o.authorizationUsecase.GrantToken(ktx, req)
	if jsonapierr != nil {
		c.JSON(jsonapierr.HTTPStatus(), jsonapi.BuildResponse(jsonapi.WithErrors(jsonapierr)))
		return
	}

	c.JSON(http.StatusOK, jsonapi.BuildResponse(
		jsonapi.WithData(data),
	))
}
