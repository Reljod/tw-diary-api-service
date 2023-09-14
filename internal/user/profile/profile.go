package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileRoute struct{}

func (route *ProfileRoute) GetProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"userId": 1})
}
