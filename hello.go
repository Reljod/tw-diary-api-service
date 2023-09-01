package main

import (
	"net/http"

	"github.com/Reljod/tw-diary-api-service/greetings"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		message := greetings.Hello(name)
		c.JSON(http.StatusOK, gin.H{
			"message": message,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
