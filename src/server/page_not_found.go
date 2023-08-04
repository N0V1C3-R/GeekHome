package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func PageNotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", gin.H{"title": "Page Not Found", "text": "404 - Page Not Found"})
}
