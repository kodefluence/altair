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

type revokeController struct {
	authService interfaces.Authorization
}

func Revoke(authService interfaces.Authorization) *revokeController {
	return &revokeController{
		authService: authService,
	}
}

func (o *revokeController) Method() string {
	return "POST"
}

func (o *revokeController) Path() string {
	return "/oauth/authorizations/revoke"
}

func (o *revokeController) Control(c *gin.Context) {
	var req entity.RevokeAccessTokenRequestJSON

	rawData, err := c.GetRawData()
	if err != nil {
		journal.Error("Cannot get raw data", err).
			SetTags("controller", "authorization", "revoke", "get_raw_data").
			Log()

		c.JSON(http.StatusBadRequest, gin.H{
			"errors": eobject.Wrap(eobject.BadRequestError("request body")),
		})
		return
	}

	err = json.Unmarshal(rawData, &req)
	if err != nil {
		journal.Error("Cannot unmarshal json", err).
			SetTags("controller", "authorization", "revoke", "unmarshal").
			Log()

		c.JSON(http.StatusBadRequest, gin.H{
			"errors": eobject.Wrap(eobject.BadRequestError("request body")),
		})
		return
	}

	entityErr := o.authService.RevokeToken(c, req)
	if entityErr != nil {
		c.JSON(entityErr.HttpStatus, gin.H{
			"errors": entityErr.Errors,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Access token has been successfully revoked.",
	})
}
