package server

import (
	"WebHome/src/server/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomeHandle(c *gin.Context) {
	var ip string
	userAuth := middleware.GetUserAuth(c)
	if userAuth.UserId != 0 {
		username = userAuth.Username
		ip = c.ClientIP()
	} else {
		username = "Visitor"
		ip = "127.0.0.1"
	}
	c.HTML(http.StatusOK, "home.html", gin.H{"title": "Terminal", "user": username, "ip": ip})
}
