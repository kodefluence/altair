package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health provide altair status
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}
