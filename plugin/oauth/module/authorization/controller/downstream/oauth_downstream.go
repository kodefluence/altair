package downstream

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	coreEntity "github.com/kodefluence/altair/entity"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/monorepo/db"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/kontext"

	"github.com/gin-gonic/gin"
)

// Oauth implement downstream plugin interface
type Oauth struct {
	oauthAccessTokenRepo OauthAccessTokenRepository
	sqldb                db.DB
}

// NewOauth create new downstream plugin to check the validity of access token given by the users
func NewOauth(oauthAccessTokenRepo OauthAccessTokenRepository, sqldb db.DB) *Oauth {
	return &Oauth{oauthAccessTokenRepo: oauthAccessTokenRepo, sqldb: sqldb}
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

	token, exc := o.oauthAccessTokenRepo.OneByToken(kontext.Fabricate(kontext.WithDefaultContext(c)), accessToken, o.sqldb)
	if exc != nil {
		if exc.Type() == exception.NotFound {
			c.AbortWithStatus(http.StatusUnauthorized)
			return exc
		}

		c.AbortWithStatus(http.StatusServiceUnavailable)
		return fmt.Errorf("Error connecting to model: %v", err)
	}

	if time.Now().After(token.ExpiresIn) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return fmt.Errorf("Token already expired: %s", token.ExpiresIn.String())
	}

	if o.validTokenScope(token, r) == false {
		c.AbortWithStatus(http.StatusForbidden)
		return fmt.Errorf("Invalid token scope: %v", token.Scopes.String)
	}

	proxyReq.Header.Add("Resource-Owner-ID", strconv.Itoa(token.ResourceOwnerID))
	proxyReq.Header.Add("Oauth-Application-ID", strconv.Itoa(token.OauthApplicationID))

	return nil
}

func (o *Oauth) validTokenScope(token entity.OauthAccessToken, r coreEntity.RouterPath) bool {
	if r.Scope == "" {
		return true
	}

	if token.Scopes.Valid {
		tokenScopes := strings.Split(token.Scopes.String, " ")
		routeScopes := strings.Split(r.Scope, " ")

		for _, routeScope := range routeScopes {
			for _, tokenScope := range tokenScopes {
				if routeScope == tokenScope {
					return true
				}
			}
		}
	}

	return false
}

func (o *Oauth) parseToken(c *gin.Context) (string, exception.Exception) {
	authorizationHeader := c.Request.Header.Get("Authorization")
	splittedToken := strings.Split(authorizationHeader, " ")

	if len(splittedToken) < 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return "", exception.Throw(errors.New("invalid request"), exception.WithTitle("bad request"), exception.WithDetail("invalid bearer format"), exception.WithType(exception.BadInput))
	}

	if splittedToken[0] != "Bearer" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return "", exception.Throw(errors.New("invalid request"), exception.WithTitle("bad request"), exception.WithDetail("invalid bearer format"), exception.WithType(exception.BadInput))
	}

	return splittedToken[1], nil
}
