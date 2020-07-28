package authorization

import (
	"encoding/json"
	"net/http"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/eobject"
	"github.com/codefluence-x/altair/provider/plugin/oauth/interfaces"
	"github.com/codefluence-x/journal"
	"github.com/gin-gonic/gin"
)

type tokenController struct {
	authService interfaces.Authorization
}

// Token create new token controller
func Token(authService interfaces.Authorization) core.Controller {
	return &tokenController{
		authService: authService,
	}
}

func (o *tokenController) Method() string {
	return "POST"
}

func (o *tokenController) Path() string {
	return "/oauth/authorizations/token"
}

func (o *tokenController) Control(c *gin.Context) {
	var req entity.AccessTokenRequestJSON

	rawData, err := c.GetRawData()
	if err != nil {
		journal.Error("Cannot get raw data", err).
			SetTags("controller", "authorization", "token", "get_raw_data").
			SetTrackId(c.Value("track_id")).
			Log()

		c.JSON(http.StatusBadRequest, gin.H{
			"errors": eobject.Wrap(eobject.BadRequestError("request body")),
		})
		return
	}

	err = json.Unmarshal(rawData, &req)
	if err != nil {
		journal.Error("Cannot unmarshal json", err).
			SetTags("controller", "authorization", "token", "unmarshal").
			SetTrackId(c.Value("track_id")).
			Log()

		c.JSON(http.StatusBadRequest, gin.H{
			"errors": eobject.Wrap(eobject.BadRequestError("request body")),
		})
		return
	}

	data, entityErr := o.authService.Token(c, req)
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
