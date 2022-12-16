package downstream

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kodefluence/altair/module"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/util"
	"github.com/kodefluence/monorepo/db"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/kontext"
)

type ApplicationValidation struct {
	oauthApplicationRepo OauthApplicationRepository
	sqldb                db.DB
}

// NewApplicationValidation create new downstream plugin to check the validity of application uid and application secret given by the client
func NewApplicationValidation(oauthApplicationRepo OauthApplicationRepository, sqldb db.DB) *ApplicationValidation {
	return &ApplicationValidation{oauthApplicationRepo: oauthApplicationRepo, sqldb: sqldb}
}

// Name of downstream plugin
func (o *ApplicationValidation) Name() string {
	return "application-validation-plugin"
}

// Intervene current request to check application_uid and application_secret
func (o *ApplicationValidation) Intervene(c *gin.Context, proxyReq *http.Request, r module.RouterPath) error {
	if r.GetAuth() != "oauth_application" {
		return nil
	}

	applicationJSON := entity.OauthApplicationJSON{}

	if proxyReq.Body == nil {
		if clientUID := c.GetHeader("CLIENT_UID"); clientUID != "" {
			applicationJSON.ClientUID = util.ValueToPointer(clientUID)
		}

		if clientSecret := c.GetHeader("CLIENT_SECRET"); clientSecret != "" {
			applicationJSON.ClientSecret = util.ValueToPointer(clientSecret)
		}

	} else {
		body, err := proxyReq.GetBody()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return exception.Throw(errors.New("internal server error"))
		}

		err = json.NewDecoder(body).Decode(&applicationJSON)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return exception.Throw(errors.New("invalid request"), exception.WithTitle("bad request"), exception.WithDetail("invalid json body"), exception.WithType(exception.BadInput))
		}
	}

	if applicationJSON.ClientUID == nil || applicationJSON.ClientSecret == nil {

		c.AbortWithStatus(http.StatusUnprocessableEntity)

		return exception.Throw(errors.New("invalid request"), exception.WithTitle("bad request"), exception.WithDetail("`client_uid` and `client_secret` can't be null"), exception.WithType(exception.BadInput))
	}

	_, exc := o.oauthApplicationRepo.OneByUIDandSecret(kontext.Fabricate(kontext.WithDefaultContext(c)), *applicationJSON.ClientUID, *applicationJSON.ClientSecret, o.sqldb)
	if exc != nil {
		if exc.Type() == exception.NotFound {
			c.AbortWithStatus(http.StatusUnauthorized)
			return exception.Throw(exc, exception.WithType(exception.Unauthorized))
		}

		c.AbortWithStatus(http.StatusServiceUnavailable)
		return exc
	}

	return nil
}
