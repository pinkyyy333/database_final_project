package middleware

import (
	"net/http"
	"os"
	"strings"

	"clinic-backend/controllers"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// AuthMiddleware 驗證 JWT 並將 patient_id 放進 context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": true, "message": "未提供授權標頭", "code": 401})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// 解析 token
		token, err := jwt.ParseWithClaims(tokenString, &controllers.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": true, "message": "驗證失敗", "code": 401})
			return
		}
		// 取出 claims，並寫入 context
		claims, ok := token.Claims.(*controllers.JWTClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": true, "message": "解析 Claims 失敗", "code": 500})
			return
		}
		c.Set("patient_id", claims.PatientID)
		c.Next()
	}
}
