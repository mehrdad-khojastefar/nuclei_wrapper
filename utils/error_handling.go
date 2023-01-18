package utils

import (
	"github.com/gin-gonic/gin"
)

func ThrowJsonError(code int, message string, c *gin.Context) {
	c.JSON(code, gin.H{
		"error": message,
	})
	c.Abort()
}
