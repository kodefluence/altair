package downstream

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	coreEntity "github.com/kodefluence/altair/entity"
	"github.com/kodefluence/altair/provider/plugin/oauth/entity"
	"github.com/kodefluence/altair/provider/plugin/oauth/interfaces"
	"github.com/kodefluence/monorepo/db"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/kontext"
)

// ApplicationValidation implement downstream plugin interface
type ApplicationValidation struct {
	oauthApplicationModel interfaces.OauthApplicationModel
	sqldb                 db.DB
}

// NewApplicationValidation create new downstream plugin to check the validity of application uid and application secret given by the client
func NewApplicationValidation(oauthApplicationModel interfaces.OauthApplicationModel, sqldb db.DB) *ApplicationValidation {
	return &ApplicationValidation{oauthApplicationModel: oauthApplicationModel, sqldb: sqldb}
}

// Name of downstream plugin
func (o *ApplicationValidation) Name() string {
	return "application-validation-plugin"
}

// Intervene current request to check application_uid and application_secret
func (o *ApplicationValidation) Intervene(c *gin.Context, proxyReq *http.Request, r coreEntity.RouterPath) error {
	if r.Auth != "oauth_application" {
		return nil
	}

	// TODO check header first for CLIENT_UID and CLIENT_SECRET then check the body. If the header exists then continue with header, if not continue with body.
	// For now we just make it like this.
	if proxyReq.Body == nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return exception.Throw(errors.New("invalid request"), exception.WithTitle("bad request"), exception.WithDetail("request body is can't be null"), exception.WithType(exception.BadInput))
	}

	body, err := proxyReq.GetBody()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return exception.Throw(errors.New("internal server error"))
	}

	applicationJSON := entity.OauthApplicationJSON{}
	err = json.NewDecoder(body).Decode(&applicationJSON)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return exception.Throw(errors.New("invalid request"), exception.WithTitle("bad request"), exception.WithDetail("invalid json body"), exception.WithType(exception.BadInput))
	}

	if applicationJSON.ClientUID == nil || applicationJSON.ClientSecret == nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return exception.Throw(errors.New("invalid request"), exception.WithTitle("bad request"), exception.WithDetail("`client_uid` and `client_secret` can't be null"), exception.WithType(exception.BadInput))
	}

	_, exc := o.oauthApplicationModel.OneByUIDandSecret(kontext.Fabricate(kontext.WithDefaultContext(c)), *applicationJSON.ClientUID, *applicationJSON.ClientSecret, o.sqldb)
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
