package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type UserAuth struct {
	UserId      int64  `json:"uid"`
	Username    string `json:"username"`
	Role        string `json:"role"`
	WorkingPath string `json:"workingPath"`
}

func UpdateAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("userAuthorization")
		if err != nil || cookie == "" {
			c.Next()
			return
		}
		c.SetCookie("userAuthorization", cookie, 3600, "/", "", true, true)
		c.Next()
	}
}

func GetUserAuth(c *gin.Context) (userAuth UserAuth) {
	authToken, err := c.Cookie("userAuthorization")
	if err != nil {
		return
	}
	if err := json.Unmarshal([]byte(authToken), &userAuth); err != nil {
		return
	}
	return
}
