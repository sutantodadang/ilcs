package middlewares

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		headerSplit := strings.Split(authHeader, " ")
		if len(headerSplit) != 2 {
			c.JSON(401, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		if headerSplit[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		_, err := jwt.Parse(headerSplit[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				c.JSON(401, gin.H{"error": "Invalid token"})
				c.Abort()
				return nil, nil
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
