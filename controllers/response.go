package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RespondOK sends a successful 200 OK response with standardized format.
func RespondOK(c *gin.Context, data gin.H) {
	data["success"] = true
	data["error"] = false
	c.JSON(http.StatusOK, data)
}

// RespondCreated sends a successful 201 Created response with standardized format.
func RespondCreated(c *gin.Context, data gin.H) {
	data["success"] = true
	data["error"] = false
	c.JSON(http.StatusCreated, data)
}

// RespondError sends an error response with standardized format.
// It will include success=false, error=true, message, and code.
func RespondError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"success": false,
		"error":   true,
		"message": message,
		"code":    code,
	})
}
