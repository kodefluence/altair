package downstream

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/entity"
	"github.com/gin-gonic/gin"
)

var InvalidBearerFormatErr = errors.New("Invalid bearer token format")

type oauth struct {
	oauthAccessTokenModel core.OauthAccessTokenModel
}

func Oauth(oauthAccessTokenModel core.OauthAccessTokenModel) core.DownStreamPlugin {
	return &oauth{oauthAccessTokenModel: oauthAccessTokenModel}
}

func (o *oauth) Name() string {
	return "oauth-plugin"
}

func (o *oauth) Intervene(c *gin.Context, proxyReq *http.Request, r entity.RouterPath) error {
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
		return errors.New(fmt.Sprintf("Error connecting to model: %v", err))
	}

	proxyReq.Header.Add("Resource-Owner-ID", strconv.Itoa(token.ResourceOwnerID))
	proxyReq.Header.Add("Oauth-Application-ID", strconv.Itoa(token.OauthApplicationID))
	return nil
}

func (o *oauth) parseToken(c *gin.Context) (string, error) {
	authorizationHeader := c.Request.Header.Get("Authorization")
	splittedToken := strings.Split(authorizationHeader, " ")

	if len(splittedToken) < 2 {
		return "", InvalidBearerFormatErr
	}

	if splittedToken[0] != "Bearer" {
		return "", InvalidBearerFormatErr
	}

	return splittedToken[1], nil
}
