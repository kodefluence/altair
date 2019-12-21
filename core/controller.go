package core

import "github.com/gin-gonic/gin"

type Controller interface {
	Control(c *gin.Context)

	// Relative path
	// /oauth/applications
	Path() string

	// GET PUT POST
	Method() string
}

type APIEngine interface {
	Handle(httpMethod, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
}
