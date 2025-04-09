package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// Check for token in Cookie
		cookieToken, err := c.Cookie("token")
		if err == nil {
			tokenString = cookieToken
		} else {
			// Check for token in Authorization header
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenString = parts[1]
				}
			}
		}

		// If no token found
		if tokenString == "" {
			// Return a specific status code and clear message
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "unauthorized",
				"message": "Authentication required",
			})
			c.Abort()
			return
		}

		// Verify and decode the token
		claims, err := ValidateJWT(tokenString)
		if err != nil {
			// For expired tokens, give a specific message
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "token_expired",
				"message": "Your session has expired, please login again",
			})
			c.Abort()
			return
		}

		encodedKey := claims["key"].(string)
		binaryKey, err := DecodeFromBase64(encodedKey)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unknown error occured with the binary key"})
		}
		// Set user ID in context for use in handlers
		c.Set("userID", claims["sub"])
		c.Set("encryptionKey", binaryKey)
		c.Next()
	}
}
