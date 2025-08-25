package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"ojire/model"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := model.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		// store user ID for later handlers
		c.Set("userID", userID)
		c.Next()
	}
}
