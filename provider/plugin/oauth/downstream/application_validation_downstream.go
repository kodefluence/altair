package downstream

import (
	"database/sql"
	"encoding/json"
	"net/http"

	coreEntity "github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/interfaces"
	"github.com/gin-gonic/gin"
)

// ApplicationValidation implement downstream plugin interface
type ApplicationValidation struct {
	oauthApplicationModel interfaces.OauthApplicationModel
}

// NewApplicationValidation create new downstream plugin to check the validity of application uid and application secret given by the client
func NewApplicationValidation(oauthApplicationModel interfaces.OauthApplicationModel) *ApplicationValidation {
	return &ApplicationValidation{oauthApplicationModel: oauthApplicationModel}
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
		return ErrInvalidRequest
	}

	body, err := proxyReq.GetBody()
	if err != nil {
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return ErrUnavailable
	}

	applicationJSON := entity.OauthApplicationJSON{}
	err = json.NewDecoder(body).Decode(&applicationJSON)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return ErrBadRequest
	}

	if applicationJSON.ClientUID == nil || applicationJSON.ClientSecret == nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return ErrInvalidRequest
	}

	_, err = o.oauthApplicationModel.OneByUIDandSecret(c, *applicationJSON.ClientUID, *applicationJSON.ClientSecret)
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatus(http.StatusUnauthorized)
			return ErrApplicationNotExists
		}

		c.AbortWithStatus(http.StatusServiceUnavailable)
		return ErrUnavailable
	}

	return nil
}
