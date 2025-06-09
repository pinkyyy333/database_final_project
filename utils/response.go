package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RespondError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"error":   true,
		"message": message,
		"code":    code,
	})
}

func RespondSuccess(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": message,
	})
}
