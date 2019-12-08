package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"meta":    gin.H{"http_status": http.StatusOK},
	})
}
