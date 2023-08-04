package middleware

import (
	"WebHome/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"path"
)

func LogMiddleware() gin.HandlerFunc {
	logger, changeSrc := getLogger()
	return func(c *gin.Context) {
		changeSrc()

		startTime := utils.GetCurrentTime()
		c.Next()
		endTime := utils.GetCurrentTime()

		latencyTime := endTime.Sub(startTime)
		reqUri := c.Request.RequestURI
		reqMethod := c.Request.Method
		statusCode := c.Writer.Status()
		reqBody := c.Request.Body

		logger.Infof("| %13v | %s | %s | %3d | %s |",
			latencyTime, reqUri, reqMethod, statusCode, reqBody)
	}
}

func getLogger() (*logrus.Logger, func()) {
	logFilePath := ""
	if dir, err := os.Getwd(); err != nil {
		logFilePath = dir + "/logs/"
	}

	utils.CreateFolder(logFilePath)
	filenameFormat := utils.FormatCurrentTime(utils.GetCurrentTime())
	src := createLogFile(logFilePath, filenameFormat)

	logger := logrus.New()
	logger.Out = src
	logger.SetLevel(logrus.InfoLevel)

	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return logger, func() {
		newFilenameFormat := utils.FormatCurrentTime(utils.GetCurrentTime())
		if newFilenameFormat != filenameFormat {
			filenameFormat = newFilenameFormat
			defer func(src *os.File) {
				err := src.Close()
				if err != nil {

				}
			}(src)
			newSrc := createLogFile(logFilePath, filenameFormat)
			logger.Out = newSrc
		}
	}
}

func createLogFile(logFilePath, filenameFormat string) *os.File {
	logFilename := filenameFormat + ".log"
	logFilename = path.Join(logFilePath, logFilename)

	checkFile(logFilename)

	src, err := os.OpenFile(logFilename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	return src
}

func checkFile(filename string) {
	if _, err := os.Stat(filename); err != nil {
		if _, err := os.Create(filename); err != nil {
			panic(err)
		}
	}
}
