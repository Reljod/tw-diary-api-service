package main

import (
	"context"
	"log"

	"github.com/Reljod/tw-diary-api-service/internal/database"
	"github.com/Reljod/tw-diary-api-service/internal/user/auth"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var secret = []byte("secret")

func main() {
	r := engine()
	r.Use(gin.Logger())
	if err := engine().Run(":8080"); err != nil {
		log.Fatal("Unable to start:", err)
	}

	defer database.Conn.Close(context.Background())
}

func engine() *gin.Engine {
	r := gin.New()

	var bcryptPwManager auth.PasswordManager = &auth.BCryptPasswordManager{}
	var authService auth.Authenticator = &auth.SimpleSessionBasedAuth{Db: database.Conn, PasswordManager: bcryptPwManager}
	authRoutes := auth.AuthRoute{Auth: authService}

	r.Use(sessions.Sessions("mysession", cookie.NewStore(secret)))
	r.POST("/login", authRoutes.LoginRoute)
	r.POST("/register", authRoutes.RegisterRoute)

	return r
}
