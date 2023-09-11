package middleware

import (
	"WebHome/src/utils"
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
)

func LogMiddleware() gin.HandlerFunc {
	logger, changeSrc := utils.GetLogger()
	return func(c *gin.Context) {
		changeSrc()

		body, _ := c.GetRawData()
		c.Request.Body = io.NopCloser(bytes.NewReader(body))

		startTime := utils.GetCurrentTime()
		c.Next()
		endTime := utils.GetCurrentTime()

		copyContext := c.Copy()
		latencyTime := endTime.Sub(startTime)
		reqUri := c.FullPath()
		reqMethod := copyContext.Request.Method
		statusCode := copyContext.Writer.Status()

		logger.Infof("| %13v | %s | %s | %3d | %s |",
			latencyTime, reqUri, reqMethod, statusCode, body)
	}
}
