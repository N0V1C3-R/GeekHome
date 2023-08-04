package command_server

import (
	"WebHome/src/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type CodeVServer struct {
	BaseCommand
}

var languageMap = map[string]string{
	"C":           "1001",
	"C#":          "1022",
	"C++":         "1002",
	"GO":          "60",
	"JAVA":        "1004",
	"JS":          "63",
	"JAVASCRIPT":  "63",
	"OBJECTIVE-C": "79",
	"PHP":         "68",
	"PY2":         "70",
	"PY3":         "71",
	"PYTHON2":     "70",
	"PYTHON3":     "71",
	"R":           "80",
	"RUBY":        "72",
	"RUST":        "73",
	"SWIFT":       "83",
	"TYPESCRIPT":  "74",
}

var indentTypeMap = map[string]string{
	"-T": "Tab",
	"-S": "Space",
}

var indentUnitMap = map[string]string{
	"2": "2",
	"4": "4",
}

func (cvs *CodeVServer) ParseCommand(stdin string) {
	cvs.Required = make(map[string]string)
	rawParts := strings.Split(stdin, " ")
	parts := utils.RemoveElements(rawParts, "").([]string)
	for i := 0; i < len(parts); i++ {
		arg := parts[i]
		if cvs.Required["languageID"] == "" {
			languageCode := languageMap[strings.ToUpper(arg)]
			if languageCode != "" {
				cvs.Required["languageID"] = languageCode
				continue
			}
		}
		if cvs.Required["indentType"] == "" {
			indentType := indentTypeMap[strings.ToUpper(arg)]
			if indentType != "" {
				cvs.Required["indentType"] = indentType
				continue
			}
		}
		if cvs.Required["indentUnit"] == "" {
			indentUnit := indentUnitMap[arg]
			if indentUnit != "" {
				cvs.Required["indentUnit"] = indentUnit
				continue
			}
		}
	}
	if cvs.Required["languageID"] == "" {
		cvs.Required["languageID"] = "71"
	}
	if cvs.Required["indentType"] == "" {
		cvs.Required["indentType"] = "space"
	}
	if cvs.Required["indentUnit"] == "" {
		cvs.Required["indentUnit"] = "4"
	}
}

func (cvs *CodeVServer) ExecuteCommand(c *gin.Context) {
	languageID := cvs.Required["languageID"]
	indentType := cvs.Required["indentType"]
	indentUnit := cvs.Required["indentUnit"]
	url := fmt.Sprintf("/codev?languageID=%s&indentType=%s&indentUnit=%s", languageID, indentType, indentUnit)
	c.Redirect(http.StatusFound, url)
}
