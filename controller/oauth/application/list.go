package application

import (
	"net/http"
	"strconv"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/eobject"
	"github.com/gin-gonic/gin"
)

type listController struct {
	applicationManager core.ApplicationManager
}

func List(applicationManager core.ApplicationManager) core.Controller {
	return &listController{
		applicationManager: applicationManager,
	}
}

func (l *listController) Method() string {
	return "GET"
}

func (l *listController) Path() string {
	return "/oauth/applications"
}

func (l *listController) Control(c *gin.Context) {
	var offset, limit int
	var err error

	offset, err = strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": eobject.Wrap(eobject.BadRequestError("query parameters: offset")),
		})
		return
	}

	limit, err = strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": eobject.Wrap(eobject.BadRequestError("query parameters: limit")),
		})
		return
	}

	oauthApplicationJSON, total, entityError := l.applicationManager.List(c, offset, limit)
	if entityError != nil {
		c.JSON(entityError.HttpStatus, gin.H{
			"errors": entityError.Errors,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": oauthApplicationJSON,
		"meta": gin.H{
			"offset": offset,
			"limit":  limit,
			"total":  total,
		},
	})
}
