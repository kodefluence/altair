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

// CreateController control flow of oauth application creation
type CreateController struct {
	applicationManager ApplicationManager
	apiError           module.ApiError
}

// NewCreate return struct of CreateController
func NewCreate(applicationManager ApplicationManager, apiError module.ApiError) *CreateController {
	return &CreateController{
		applicationManager: applicationManager,
		apiError:           apiError,
	}
}

// Method POST
func (cr *CreateController) Method() string {
	return "POST"
}

// Path /oauth/applications
func (cr *CreateController) Path() string {
	return "/oauth/applications"
}

// Control creation of oauth application
func (cr *CreateController) Control(ktx kontext.Context, c *gin.Context) {
	var oauthApplicationJSON entity.OauthApplicationJSON

	rawData, err := c.GetRawData()
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", c.Value("request_id")).
			Array("tags", zerolog.Arr().Str("controller").Str("application").Str("create").Str("get_raw_data")).
			Msg("Cannot get raw data")
		c.JSON(http.StatusBadRequest, jsonapi.BuildResponse(cr.apiError.BadRequestError("request body")))
		return
	}

	err = json.Unmarshal(rawData, &oauthApplicationJSON)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", c.Value("request_id")).
			Array("tags", zerolog.Arr().Str("controller").Str("application").Str("update").Str("unmarshal")).
			Msg("Cannot unmarshal json")
		c.JSON(http.StatusBadRequest, jsonapi.BuildResponse(cr.apiError.BadRequestError("invalid json format")))
		return
	}

	result, jsonapiErr := cr.applicationManager.Create(ktx, oauthApplicationJSON)
	if jsonapiErr != nil {
		c.JSON(jsonapiErr.HTTPStatus(), jsonapi.BuildResponse(jsonapi.WithErrors(jsonapiErr)))
		return
	}

	c.JSON(http.StatusCreated, jsonapi.BuildResponse(jsonapi.WithData(result)))
}
