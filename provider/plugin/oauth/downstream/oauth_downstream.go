package downstream

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	coreEntity "github.com/codefluence-x/altair/entity"

	"github.com/codefluence-x/altair/provider/plugin/oauth/interfaces"
	"github.com/gin-gonic/gin"
)

// ErrInvalidBearerFormat returned when the header of bearer token is invalid
var ErrInvalidBearerFormat = errors.New("Invalid bearer token format")

// Oauth implement downstream plugin interface
type Oauth struct {
	oauthAccessTokenModel interfaces.OauthAccessTokenModel
}

// NewOauth create new downstream plugin to check the validity of access token given by the users
func NewOauth(oauthAccessTokenModel interfaces.OauthAccessTokenModel) *Oauth {
	return &Oauth{oauthAccessTokenModel: oauthAccessTokenModel}
}

// Name get the name of downstream plugin
func (o *Oauth) Name() string {
	return "oauth-plugin"
}

// Intervene current request to check the bearer token validity
func (o *Oauth) Intervene(c *gin.Context, proxyReq *http.Request, r coreEntity.RouterPath) error {
	if r.Auth != "oauth" {
		return nil
	}

	accessToken, err := o.parseToken(c)
	if err != nil {
		return err
	}

	token, err := o.oauthAccessTokenModel.OneByToken(c, accessToken)
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatus(http.StatusUnauthorized)
			return err
		}

		c.AbortWithStatus(http.StatusServiceUnavailable)
		return fmt.Errorf("Error connecting to model: %v", err)
	}

	proxyReq.Header.Add("Resource-Owner-ID", strconv.Itoa(token.ResourceOwnerID))
	proxyReq.Header.Add("Oauth-Application-ID", strconv.Itoa(token.OauthApplicationID))
	return nil
}

func (o *Oauth) parseToken(c *gin.Context) (string, error) {
	authorizationHeader := c.Request.Header.Get("Authorization")
	splittedToken := strings.Split(authorizationHeader, " ")

	if len(splittedToken) < 2 {
		return "", ErrInvalidBearerFormat
	}

	if splittedToken[0] != "Bearer" {
		return "", ErrInvalidBearerFormat
	}

	return splittedToken[1], nil
}
