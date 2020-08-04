package application

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/eobject"
	"github.com/codefluence-x/altair/provider/plugin/oauth/interfaces"
	"github.com/codefluence-x/journal"
	"github.com/gin-gonic/gin"
)

type updateController struct {
	applicationManager interfaces.ApplicationManager
}

func Update(applicationManager interfaces.ApplicationManager) core.Controller {
	return &updateController{
		applicationManager: applicationManager,
	}
}

func (cr *updateController) Method() string {
	return "PUT"
}

func (cr *updateController) Path() string {
	return "/oauth/applications/:id"
}

func (cr *updateController) Control(c *gin.Context) {
	var oauthApplicationUpdateJSON entity.OauthApplicationUpdateJSON

	rawData, err := c.GetRawData()
	if err != nil {
		journal.Error("Cannot get raw data", err).
			SetTags("controller", "application", "update", "get_raw_data").
			SetTrackId(c.Value("track_id")).
			Log()

		c.JSON(http.StatusBadRequest, gin.H{
			"errors": eobject.Wrap(eobject.BadRequestError("request body")),
		})
		return
	}

	err = json.Unmarshal(rawData, &oauthApplicationUpdateJSON)
	if err != nil {
		journal.Error("Cannot unmarshal json", err).
			SetTags("controller", "application", "update", "unmarshal").
			SetTrackId(c.Value("track_id")).
			Log()

		c.JSON(http.StatusBadRequest, gin.H{
			"errors": eobject.Wrap(eobject.BadRequestError("request body")),
		})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		journal.Error("Cannot convert ascii to integer", err).
			SetTags("controller", "application", "update", "strconv").
			SetTrackId(c.Value("track_id")).
			Log()

		c.JSON(http.StatusBadRequest, gin.H{
			"errors": eobject.Wrap(eobject.BadRequestError("url parameters: id is not integer")),
		})
		return
	}

	result, entityError := cr.applicationManager.Update(c, id, oauthApplicationUpdateJSON)
	if entityError != nil {
		c.JSON(entityError.HttpStatus, gin.H{
			"errors": entityError.Errors,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}
