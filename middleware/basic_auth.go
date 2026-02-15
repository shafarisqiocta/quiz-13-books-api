package middleware

import (
	"encoding/base64"
	"net/http"
	"strings"

	"quiz-13-books-api/config"

	"github.com/gin-gonic/gin"
)

func BasicAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		splitToken := strings.Split(authHeader, " ")
		if len(splitToken) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid auth format"})
			c.Abort()
			return
		}

		decoded, _ := base64.StdEncoding.DecodeString(splitToken[1])
		credentials := strings.Split(string(decoded), ":")

		if len(credentials) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials format"})
			c.Abort()
			return
		}

		username := credentials[0]
		password := credentials[1]

		var dbPassword string
		query := `SELECT password FROM users WHERE username=$1`

		err := config.DB.QueryRow(query, username).Scan(&dbPassword)
		if err != nil || dbPassword != password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			c.Abort()
			return
		}

		c.Next()
	}
}
