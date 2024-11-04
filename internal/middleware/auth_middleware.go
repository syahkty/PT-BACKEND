package middleware

import (
	"fmt"
	"library_api/pkg/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is used to validate JWT and set the role of the user
func AuthMiddleware(jwtAuth *auth.JWTAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}
		
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwtAuth.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set the username and role from token claims
		c.Set("id",claims.ID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		fmt.Printf("Authenticated user. ID: %d, Username: %s, Role: %s\n",claims.ID, claims.Username, claims.Role)
		c.Next()
	}
}

// AdminOnlyMiddleware checks if the user is an admin
func AdminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			fmt.Printf("Unauthorized access attempt. Role: %v\n", role)
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this resource"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// UserOnlyMiddleware checks if the user is a regular user
func UserOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "user" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this resource"})
			c.Abort()
			return
		}
		c.Next()
	}
}
