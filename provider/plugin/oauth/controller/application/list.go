package application

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kodefluence/altair/provider/plugin/oauth/eobject"
	"github.com/kodefluence/altair/provider/plugin/oauth/interfaces"
)

// ListController show list of oauth applications
type ListController struct {
	applicationManager interfaces.ApplicationManager
}

// NewList return struct of ListController
func NewList(applicationManager interfaces.ApplicationManager) *ListController {
	return &ListController{
		applicationManager: applicationManager,
	}
}

// Method GET
func (l *ListController) Method() string {
	return "GET"
}

// Path /oauth/applications
func (l *ListController) Path() string {
	return "/oauth/applications"
}

// Control list of oauth applications
func (l *ListController) Control(c *gin.Context) {
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
