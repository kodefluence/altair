package authorization

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kodefluence/altair/provider/plugin/oauth/entity"
	"github.com/kodefluence/altair/provider/plugin/oauth/eobject"
	"github.com/kodefluence/altair/provider/plugin/oauth/interfaces"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// TokenController control flow of creating access token
type TokenController struct {
	authService interfaces.Authorization
}

// NewToken create new token controller
func NewToken(authService interfaces.Authorization) *TokenController {
	return &TokenController{
		authService: authService,
	}
}

// Method POST
func (o *TokenController) Method() string {
	return "POST"
}

// Path /oauth/authorizations/token
func (o *TokenController) Path() string {
	return "/oauth/authorizations/token"
}

// Control creating access token based on access token request
func (o *TokenController) Control(c *gin.Context) {
	var req entity.AccessTokenRequestJSON

	rawData, err := c.GetRawData()
	if err != nil {
		log.Error().
			Err(err).
			Stack().
			Interface("request_id", c.Value("request_id")).
			Array("tags", zerolog.Arr().Str("controller").Str("authorization").Str("token").Str("get_raw_data")).
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
			Array("tags", zerolog.Arr().Str("controller").Str("authorization").Str("token").Str("unmarshal")).
			Msg("Cannot unmarshal json")

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
