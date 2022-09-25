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

// GrantController control flow of grant access token / authorization code
type GrantController struct {
	authorizationUsecase Authorization
	apiError             module.ApiError
}

// NewGrant return struct ob GrantController
func NewGrant(authorizationUsecase Authorization, apiError module.ApiError) *GrantController {
	return &GrantController{
		authorizationUsecase: authorizationUsecase,
		apiError:             apiError,
	}
}

// Method Post
func (o *GrantController) Method() string {
	return "POST"
}

// Path /oauth/authorizations
func (o *GrantController) Path() string {
	return "/oauth/authorizations"
}

// Control granting access token / authorization code
func (o *GrantController) Control(c *gin.Context) {
	ktx := kontext.Fabricate(kontext.WithDefaultContext(c))
	ktx.Set("request_id", c.GetString("request_id"))

	var req entity.AuthorizationRequestJSON

	rawData, err := c.GetRawData()
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", c.Value("request_id")).
			Array("tags", zerolog.Arr().Str("controller").Str("authorization").Str("grant").Str("get_raw_data")).
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
			Array("tags", zerolog.Arr().Str("controller").Str("authorization").Str("grant").Str("unmarshal")).
			Msg("Cannot unmarshal json")
		c.JSON(http.StatusBadRequest, jsonapi.BuildResponse(o.apiError.BadRequestError("request body")))
		return
	}

	data, jsonapierr := o.authorizationUsecase.GrantAuthorizationCode(ktx, req)
	if jsonapierr != nil {
		c.JSON(jsonapierr.HTTPStatus(), jsonapi.BuildResponse(jsonapi.WithErrors(jsonapierr)))
		return
	}

	c.JSON(http.StatusOK, jsonapi.BuildResponse(
		jsonapi.WithData(data),
	))
}
