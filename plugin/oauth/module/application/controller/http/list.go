package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
)

// ListController show list of oauth applications
type ListController struct {
	applicationManager ApplicationManager
	apiError           ApiError
}

// NewList return struct of ListController
func NewList(applicationManager ApplicationManager, apiError ApiError) *ListController {
	return &ListController{
		applicationManager: applicationManager,
		apiError:           apiError,
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
	ktx := kontext.Fabricate(kontext.WithDefaultContext(c))
	ktx.Set("request_id", c.GetString("request_id"))

	var offset, limit int
	var err error

	offset, err = strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, jsonapi.BuildResponse(l.apiError.BadRequestError("query parameters: offset")))
		return
	}

	limit, err = strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, jsonapi.BuildResponse(l.apiError.BadRequestError("query parameters: limit")))
		return

	}

	oauthApplicationJSON, total, jsonapiErr := l.applicationManager.List(ktx, offset, limit)
	if jsonapiErr != nil {
		c.JSON(jsonapiErr.HTTPStatus(), jsonapi.BuildResponse(jsonapi.WithErrors(jsonapiErr)))
		return
	}

	c.JSON(http.StatusOK, jsonapi.BuildResponse(
		jsonapi.WithData(oauthApplicationJSON),
		jsonapi.WithMeta("offset", offset),
		jsonapi.WithMeta("limit", limit),
		jsonapi.WithMeta("total", total),
	))
}
