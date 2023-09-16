package auth

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

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
	var cookieSession *Session = nil
	cookie, cErr := c.Cookie("SESSION")

	if cookie != "" && cErr == nil {
		var cookieBytes []byte = make([]byte, len(cookie))

		n, err := b64.StdEncoding.Decode(cookieBytes, []byte(cookie))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication error"})
			return
		}

		err = json.Unmarshal(cookieBytes[:n], &cookieSession)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unmarshal? %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication error"})
			return
		}
	}

	if cookieSession != nil {
		isValid, err := route.Auth.IsSessionValid(cookieSession.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication error"})
			return
		}

		if isValid {
			c.JSON(http.StatusOK, gin.H{"message": "User already authenticated"})
			return
		}
	}

	username, password, ok := c.Request.BasicAuth()

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication error"})
		return
	}

	if err := route.Auth.Login(username, password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newSession, err := route.Auth.CreateSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newSessionJson, err := json.Marshal(newSession)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	sessionEnc := b64.StdEncoding.EncodeToString(newSessionJson)

	c.SetCookie("SESSION", sessionEnc, int(newSession.Expiry), "/", "localhost", false, false)
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

func (route *AuthRoute) LogoutRoute(c *gin.Context) {
	sessionId, isExists := c.Get("sessionId")
	if sessionId == nil || !isExists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Cannot access request"})
		return
	}

	err := route.Auth.Logout(fmt.Sprintf("%v", sessionId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication error"})
		return
	}

	c.Set("sessionId", nil)
	c.JSON(http.StatusOK, gin.H{"message": "Logout success"})
}
