package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func VerifyCode(c *gin.Context) {

}

func EmailVerify(c *gin.Context) {
	requestID := c.Request.Header.Get("X-Request-ID")
	var emailVerifyForm emailVerifyForm
	if err := c.ShouldBind(&emailVerifyForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"response": err.Error()})
		return
	}
	value, err := rdb.Get(ctx, requestID).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"response": "Invalid request."})
		return
	} else {
		var r registerRedis
		err = json.Unmarshal([]byte(value), &r)
		if err != nil {
			c.JSON(419, gin.H{"response": "Request Expired."})
			return
		}
		if r.VerifyCode == emailVerifyForm.Code {
			c.JSON(http.StatusOK, gin.H{"response": "Verification success!"})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"response": "Error in inputting verification code."})
		}
	}
}
