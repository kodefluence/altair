package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// UpdateController control flow of update oauth application
type UpdateController struct {
	applicationManager ApplicationManager
	apiError           module.ApiError
}

// NewUpdate create struct of UpdateController
func NewUpdate(applicationManager ApplicationManager, apiError module.ApiError) *UpdateController {
	return &UpdateController{
		applicationManager: applicationManager,
		apiError:           apiError,
	}
}

// Method PUT
func (uc *UpdateController) Method() string {
	return "PUT"
}

// Path /oauth/applications/:id
func (uc *UpdateController) Path() string {
	return "/oauth/applications/:id"
}

// Control update oauth applications
func (uc *UpdateController) Control(c *gin.Context) {
	ktx := kontext.Fabricate(kontext.WithDefaultContext(c))
	ktx.Set("request_id", c.GetString("request_id"))

	var oauthApplicationUpdateJSON entity.OauthApplicationUpdateJSON

	rawData, err := c.GetRawData()
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", c.Value("request_id")).
			Array("tags", zerolog.Arr().Str("controller").Str("application").Str("update").Str("get_raw_data")).
			Msg("Cannot get raw data")

		c.JSON(http.StatusBadRequest, jsonapi.BuildResponse(uc.apiError.BadRequestError("request body")))
		return
	}

	err = json.Unmarshal(rawData, &oauthApplicationUpdateJSON)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", c.Value("request_id")).
			Array("tags", zerolog.Arr().Str("controller").Str("application").Str("update").Str("unmarshal")).
			Msg("Cannot unmarshal json")
		c.JSON(http.StatusBadRequest, jsonapi.BuildResponse(uc.apiError.BadRequestError("invalid json format")))
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", c.Value("request_id")).
			Array("tags", zerolog.Arr().Str("controller").Str("application").Str("update").Str("strconv")).
			Msg("Cannot convert ascii to integer")

		c.JSON(http.StatusBadRequest, jsonapi.BuildResponse(uc.apiError.BadRequestError("url parameters: id is not integer")))
		return
	}

	result, jsonapiErr := uc.applicationManager.Update(ktx, id, oauthApplicationUpdateJSON)
	if jsonapiErr != nil {
		c.JSON(jsonapiErr.HTTPStatus(), jsonapi.BuildResponse(jsonapi.WithErrors(jsonapiErr)))
		return
	}

	c.JSON(http.StatusOK, jsonapi.BuildResponse(jsonapi.WithData(result)))
}
