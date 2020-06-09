package application

import (
	"net/http"
	"strconv"

	"github.com/codefluence-x/altair/provider/plugin/oauth/eobject"
	"github.com/codefluence-x/altair/provider/plugin/oauth/interfaces"
	"github.com/codefluence-x/journal"
	"github.com/gin-gonic/gin"
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
		journal.Error("Cannot convert ascii to integer", err).
			SetTags("controller", "application", "one", "strconv").
			Log()

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
