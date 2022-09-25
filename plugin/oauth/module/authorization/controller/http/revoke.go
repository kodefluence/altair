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

// RevokeController control flow of revoking access token
type RevokeController struct {
	authorizationUsecase Authorization
	apiError             module.ApiError
}

// NewRevoke create new revoke controller
func NewRevoke(authorizationUsecase Authorization, apiError module.ApiError) *RevokeController {
	return &RevokeController{
		authorizationUsecase: authorizationUsecase,
		apiError:             apiError,
	}
}

// Method POST
func (o *RevokeController) Method() string {
	return "POST"
}

// Path /oauth/authorizations/revoke
func (o *RevokeController) Path() string {
	return "/oauth/authorizations/revoke"
}

// Control revoking access token
func (o *RevokeController) Control(c *gin.Context) {
	ktx := kontext.Fabricate(kontext.WithDefaultContext(c))
	ktx.Set("request_id", c.GetString("request_id"))

	var req entity.RevokeAccessTokenRequestJSON

	rawData, err := c.GetRawData()
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", c.Value("request_id")).
			Array("tags", zerolog.Arr().Str("controller").Str("authorization").Str("revoke").Str("get_raw_data")).
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
			Array("tags", zerolog.Arr().Str("controller").Str("authorization").Str("revoke").Str("unmarshal")).
			Msg("Cannot unmarshal json")

		c.JSON(http.StatusBadRequest, jsonapi.BuildResponse(o.apiError.BadRequestError("request body")))
		return
	}

	jsonapierr := o.authorizationUsecase.RevokeToken(ktx, req)
	if jsonapierr != nil {
		c.JSON(jsonapierr.HTTPStatus(), jsonapi.BuildResponse(jsonapi.WithErrors(jsonapierr)))
		return
	}

	c.JSON(http.StatusOK, jsonapi.BuildResponse())
}
