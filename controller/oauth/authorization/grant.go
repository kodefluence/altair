package authorization

import (
	"net/http"

	"github.com/codefluence-x/altair/core"
	"github.com/gin-gonic/gin"
)

type grantController struct {
	authService core.Authorization
}

func Grant(authService core.Authorization) core.Controller {
	return &grantController{
		authService: authService,
	}
}

func (o *grantController) Method() string {
	return "POST"
}

func (o *grantController) Path() string {
	return "/oauth/authorizations"
}

func (o *grantController) Control(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{})
}
