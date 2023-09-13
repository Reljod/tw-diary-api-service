package main

import (
	"context"
	"log"

	"github.com/Reljod/tw-diary-api-service/config"
	"github.com/Reljod/tw-diary-api-service/internal/cache"
	"github.com/Reljod/tw-diary-api-service/internal/database"
	"github.com/Reljod/tw-diary-api-service/internal/user/auth"

	"github.com/gin-gonic/gin"
)

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

	var redisCache *cache.RedisCache = cache.CreateRedisCache(&config.Config)

	var options cache.SessionCacheOptions = cache.SessionCacheOptions{Expiry: config.Config.Session.Expiry, Prefix: "session:"}
	var sessionCache cache.SessionCache = &cache.SessionRedisCache{Redis: redisCache, Options: &options}
	var sessionHandler auth.SessionHandler = &auth.SimpleSessionHandler{Cache: sessionCache, Config: &config.Config}
	var bcryptPwManager auth.PasswordManager = &auth.BCryptPasswordManager{}
	var authService auth.Authenticator = &auth.SimpleSessionBasedAuth{
		Db: database.Conn, PasswordManager: bcryptPwManager, SessionHandler: sessionHandler}
	authRoutes := auth.AuthRoute{Auth: authService}

	r.Use(sessions.Sessions("mysession", cookie.NewStore(secret)))
	r.POST("/login", authRoutes.LoginRoute)
	r.POST("/register", authRoutes.RegisterRoute)

	return r
}
