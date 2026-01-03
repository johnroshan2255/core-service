package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey string

func SetJWTKey(key string) {
	jwtKey = key
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if jwtKey == "" {
			log.Printf("JWT middleware: JWT key not set")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT authentication not configured"})
			c.Abort()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(jwtKey), nil
		})

		if err != nil {
			log.Printf("JWT middleware: Token validation error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		userUUID, ok := claims["user_uuid"].(string)
		if !ok {
			userUUID, ok = claims["uuid"].(string)
		}
		if ok && userUUID != "" {
			c.Set("user_uuid", userUUID)
		}

		email, ok := claims["email"].(string)
		if ok && email != "" {
			c.Set("email", email)
		}

		c.Next()
	}
}

