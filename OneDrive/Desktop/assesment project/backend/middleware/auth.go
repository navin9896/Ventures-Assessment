package middleware

import (
	"net/http"
	"strings"

	"shopping-cart/database"
	"shopping-cart/models"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" or just use the header value
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			// If no Bearer prefix, try using the header value directly
			token = authHeader
		}
		token = strings.TrimSpace(token)

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
			c.Abort()
			return
		}

		// Find user by token
		var user models.User
		if err := database.DB.Where("token = ?", token).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Inject user into context
		c.Set("user", &user)
		c.Set("user_id", user.ID)
		c.Next()
	}
}

