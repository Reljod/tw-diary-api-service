package middleware

import (
	b64 "encoding/base64"
	"net/http"

	"github.com/Reljod/tw-diary-api-service/internal/user/auth"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	SessionHandler auth.SessionHandler
}

func (middleware *AuthMiddleware) Authenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("SESSION")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Cannot access request"})
			return
		}

		var cookieBytes []byte = make([]byte, len(cookie))

		n, err := b64.StdEncoding.Decode(cookieBytes, []byte(cookie))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}

		session, err := middleware.SessionHandler.Decode(string(cookieBytes[:n]))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}

		isValid, err := middleware.SessionHandler.IsValid(session.Id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}

		if !isValid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Session already invalid"})
			return
		}

		c.Set("session", session)
		c.Next()
	}
}
