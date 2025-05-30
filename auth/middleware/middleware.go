package middleware

import (
	"net/http"
	"strings"

	"github.com/XRS0/ToTalkB/auth"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func UserIdentity(c *gin.Context) {
	authorization := c.GetHeader("Authorization")
	parts := strings.Split(authorization, " ")
	if len(parts) < 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "authorization header format must be Bearer {token}"})
		return
	}

	userId, err := auth.ParseAccessToken(parts[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
		return
	}

	c.Set("userId", userId)
	c.Next()
}
