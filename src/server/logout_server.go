package server

import (
	"WebHome/src/server/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LogoutHandle(c *gin.Context) {
	userAuth := middleware.GetUserAuth(c)
	if userAuth.UserId != 0 {
		c.SetCookie("__userInfo", "", -1, "/", "", true, true)
		c.SetCookie("userAuthorization", "", -1, "/", "", true, true)
		c.JSON(http.StatusOK, gin.H{"response": true})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"response": false})
	}
}
