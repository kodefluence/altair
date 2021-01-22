package application

import (
	"net/http"
	"strconv"

	"github.com/codefluence-x/altair/provider/plugin/oauth/eobject"
	"github.com/codefluence-x/altair/provider/plugin/oauth/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type oneController struct {
	applicationManager interfaces.ApplicationManager
}

func One(applicationManager interfaces.ApplicationManager) *oneController {
	return &oneController{
		applicationManager: applicationManager,
	}
}

func (o *oneController) Method() string {
	return "GET"
}

func (o *oneController) Path() string {
	return "/oauth/applications/:id"
}

func (o *oneController) Control(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", c.Value("request_id")).
			Array("tags", zerolog.Arr().Str("controller").Str("application").Str("one").Str("strconv")).
			Msg("Cannot convert ascii to integer")
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": eobject.Wrap(eobject.BadRequestError("url parameters: id is not integer")),
		})
		return
	}

	oauthApplicationJSON, entityError := o.applicationManager.One(c, id)
	if entityError != nil {
		c.JSON(entityError.HttpStatus, gin.H{
			"errors": entityError.Errors,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": oauthApplicationJSON,
	})
}
