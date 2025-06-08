package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Unauthorized", "code": 401})
			c.Abort()
			return
		}
		// 實務上要驗 token 有效性並把 patient_id 放入 context
		c.Next()
	}
}
