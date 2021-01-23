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

// GrantController control flow of grant access token / authorization code
type GrantController struct {
	authService interfaces.Authorization
}

// NewGrant return struct ob GrantController
func NewGrant(authService interfaces.Authorization) *GrantController {
	return &GrantController{
		authService: authService,
	}
}

// Method Post
func (o *GrantController) Method() string {
	return "POST"
}

// Path /oauth/authorizations
func (o *GrantController) Path() string {
	return "/oauth/authorizations"
}

// Control granting access token / authorization code
func (o *GrantController) Control(c *gin.Context) {
	var req entity.AuthorizationRequestJSON

	rawData, err := c.GetRawData()
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", c.Value("request_id")).
			Array("tags", zerolog.Arr().Str("controller").Str("authorization").Str("grant").Str("get_raw_data")).
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
			Array("tags", zerolog.Arr().Str("controller").Str("authorization").Str("grant").Str("unmarshal")).
			Msg("Cannot unmarshal json")
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
