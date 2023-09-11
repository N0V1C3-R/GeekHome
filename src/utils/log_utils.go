package utils

import (
	"github.com/sirupsen/logrus"
	"os"
	"path"
)

func GetLogger() (*logrus.Logger, func()) {
	logFilePath := os.Getenv("LOG_PATH")
	CreateFolder(logFilePath)

	logFileName := os.Getenv("PROJECT_NAME")
	logFileName += "_" + GetCurrentTime().Format("20060102")
	src := createLogFile(logFilePath, logFileName)

	logger := logrus.New()
	logger.Out = src
	logger.SetLevel(logrus.InfoLevel)

	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return logger, func() {
		newFilenameFormat := os.Getenv("PROJECT_NAME") + "_" + GetCurrentTime().Format("20060102")
		if newFilenameFormat != logFileName {
			logFileName = newFilenameFormat
			defer func(src *os.File) {
				err := src.Close()
				if err != nil {
				}
			}(src)
			newSrc := createLogFile(logFilePath, logFileName)
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
