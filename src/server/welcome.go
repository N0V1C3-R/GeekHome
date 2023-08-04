package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func WelcomeHandle(c *gin.Context) {
	c.HTML(http.StatusOK, "welcome.html", gin.H{"title": "Hello World"})
}
