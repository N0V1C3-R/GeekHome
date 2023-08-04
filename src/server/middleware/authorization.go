package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

type UserAuth struct {
	UserId   int64  `json:"uid"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken, err := c.Cookie("userAuthorization")
		if err != nil {
			c.Redirect(302, "/login")
			return
		}
		fmt.Println(authToken)
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
