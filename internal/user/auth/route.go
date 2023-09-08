package auth

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthRoute struct {
	Auth Authenticator
}

type RequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (route *AuthRoute) LoginRoute(c *gin.Context) {
	session := sessions.Default(c)
	username, password, ok := c.Request.BasicAuth()

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication error"})
		return
	}

	if err := route.Auth.Login(username, password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	session.Set(username, password)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
}

func (route *AuthRoute) RegisterRoute(c *gin.Context) {
	var body RequestBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication error"})
		return
	}

	id, err := route.Auth.Register(body.Username, body.Password)
	if err != nil {

		target := &UsernameAlreadyTaken{}
		if errors.As(err, &target) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}
