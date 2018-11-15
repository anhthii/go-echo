package middlewares

import (
	"os"

	"github.com/anhthii/go-echo/db/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// TokenAuthMiddleware authenticate routes with token
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header.Get("Authorization") // Grab the token from the header

		if tokenHeader == "" {
			c.JSON(403, "No access token provided")
			c.Abort()
			return
		}

		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenHeader, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_SECRET")), nil
		})

		if err != nil { // malformed token, returns with http code 403
			c.JSON(403, "Malformed token")
			c.Abort()
			return
		}

		if !token.Valid { // token is invalid, but not signed by this server
			c.JSON(403, "Invalid token")
			c.Abort()
			return
		}

		c.Set("user", tk.Username)
		c.Next()
	}
}

// IsValidUser check if a user really owns this token
func IsValidUser(c *gin.Context) {
	userNameFromToken := c.MustGet("user")
	username := c.Param("username")
	if userNameFromToken != username {
		c.JSON(401, "You are not allowed to access this route")
		c.Abort()
		return
	}
	c.Next()
}
