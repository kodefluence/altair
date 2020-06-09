package authorization

import (
	"encoding/json"
	"net/http"

	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/eobject"
	"github.com/codefluence-x/altair/provider/plugin/oauth/interfaces"
	"github.com/codefluence-x/journal"
	"github.com/gin-gonic/gin"
)

type grantController struct {
	authService interfaces.Authorization
}

func Grant(authService interfaces.Authorization) *grantController {
	return &grantController{
		authService: authService,
	}
}

func (o *grantController) Method() string {
	return "POST"
}

func (o *grantController) Path() string {
	return "/oauth/authorizations"
}

func (o *grantController) Control(c *gin.Context) {
	var req entity.AuthorizationRequestJSON

	rawData, err := c.GetRawData()
	if err != nil {
		journal.Error("Cannot get raw data", err).
			SetTags("controller", "authorization", "grant", "get_raw_data").
			Log()

		c.JSON(http.StatusBadRequest, gin.H{
			"errors": eobject.Wrap(eobject.BadRequestError("request body")),
		})
		return
	}

	err = json.Unmarshal(rawData, &req)
	if err != nil {
		journal.Error("Cannot unmarshal json", err).
			SetTags("controller", "authorization", "grant", "unmarshal").
			Log()

		c.JSON(http.StatusBadRequest, gin.H{
			"errors": eobject.Wrap(eobject.BadRequestError("request body")),
		})
		return
	}

	data, entityErr := o.authService.Grantor(c, req)
	if entityErr != nil {
		c.JSON(entityErr.HttpStatus, gin.H{
			"errors": entityErr.Errors,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": data,
	})
}
