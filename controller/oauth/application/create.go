package application

import (
	"encoding/json"
	"net/http"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/eobject"
	"github.com/codefluence-x/journal"
	"github.com/gin-gonic/gin"
)

type createController struct {
	applicationManager core.ApplicationManager
}

func Create(applicationManager core.ApplicationManager) core.Controller {
	return &createController{
		applicationManager: applicationManager,
	}
}

func (cr *createController) Method() string {
	return "POST"
}

func (cr *createController) Path() string {
	return "/oauth/applications"
}

func (cr *createController) Control(c *gin.Context) {
	var oauthApplicationJSON entity.OauthApplicationJSON

	rawData, err := c.GetRawData()
	if err != nil {
		journal.Error("Cannot get raw data", err).
			SetTags("controller", "application", "create", "get_raw_data").
			Log()

		c.JSON(http.StatusBadRequest, gin.H{
			"errors": eobject.Wrap(eobject.BadRequestError("request body")),
		})
		return
	}

	err = json.Unmarshal(rawData, &oauthApplicationJSON)
	if err != nil {
		journal.Error("Cannot unmarshal json", err).
			SetTags("controller", "application", "create", "unmarshal").
			Log()

		c.JSON(http.StatusBadRequest, gin.H{
			"errors": eobject.Wrap(eobject.BadRequestError("request body")),
		})
		return
	}

	result, entityError := cr.applicationManager.Create(c, oauthApplicationJSON)
	if entityError != nil {
		c.JSON(entityError.HttpStatus, gin.H{
			"errors": entityError.Errors,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": result,
	})
}
