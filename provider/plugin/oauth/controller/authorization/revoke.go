package authorization

import (
	"encoding/json"
	"net/http"

	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/eobject"
	"github.com/codefluence-x/altair/provider/plugin/oauth/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// RevokeController control flow of revoking access token
type RevokeController struct {
	authService interfaces.Authorization
}

// NewRevoke create new revoke controller
func NewRevoke(authService interfaces.Authorization) *RevokeController {
	return &RevokeController{
		authService: authService,
	}
}

// Method POST
func (o *RevokeController) Method() string {
	return "POST"
}

// Path /oauth/authorizations/revoke
func (o *RevokeController) Path() string {
	return "/oauth/authorizations/revoke"
}

// Control revoking access token
func (o *RevokeController) Control(c *gin.Context) {
	var req entity.RevokeAccessTokenRequestJSON

	rawData, err := c.GetRawData()
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", c.Value("request_id")).
			Array("tags", zerolog.Arr().Str("controller").Str("authorization").Str("revoke").Str("get_raw_data")).
			Msg("Cannot get raw data")

		c.JSON(http.StatusBadRequest, gin.H{
			"errors": eobject.Wrap(eobject.BadRequestError("request body")),
		})
		return
	}

	err = json.Unmarshal(rawData, &req)
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", c.Value("request_id")).
			Array("tags", zerolog.Arr().Str("controller").Str("authorization").Str("revoke").Str("unmarshal")).
			Msg("Cannot unmarshal json")

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
