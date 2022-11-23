package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// OneController control flow of showing oauth applications detail
type OneController struct {
	applicationManager ApplicationManager
	apiError           module.ApiError
}

// NewOne return struct of OneController
func NewOne(applicationManager ApplicationManager, apiError module.ApiError) *OneController {
	return &OneController{
		applicationManager: applicationManager,
		apiError:           apiError,
	}
}

// Method GET
func (o *OneController) Method() string {
	return "GET"
}

// Path /oauth/applications/:id
func (o *OneController) Path() string {
	return "/oauth/applications/:id"
}

// Control find oauth application
func (o *OneController) Control(ktx kontext.Context, c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", c.Value("request_id")).
			Array("tags", zerolog.Arr().Str("controller").Str("application").Str("one").Str("strconv")).
			Msg("Cannot convert ascii to integer")

		c.JSON(http.StatusBadRequest, jsonapi.BuildResponse(o.apiError.BadRequestError("url parameters: id is not integer")))
		return
	}

	oauthApplicationJSON, jsonAPIErr := o.applicationManager.One(ktx, id)
	if jsonAPIErr != nil {
		c.JSON(jsonAPIErr.HTTPStatus(), jsonapi.BuildResponse(jsonapi.WithErrors(jsonAPIErr)))
		return
	}

	c.JSON(http.StatusOK, jsonapi.BuildResponse(jsonapi.WithData(oauthApplicationJSON)))
}
