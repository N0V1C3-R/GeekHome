package command_server

import (
	"WebHome/src/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
)

type CodeVServer struct {
	BaseCommand
}

type BlogServer struct {
	BaseCommand
}

type OpenServer struct {
	BaseCommand
}

type GoogleServer struct {
	BaseCommand
}

type BingServer struct {
	BaseCommand
}

type GitHubServer struct {
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
	cvs.Options = make(map[string]string)
	rawParts := strings.Split(stdin, " ")
	parts := utils.RemoveElements(rawParts, "").([]string)
	for i := 0; i < len(parts); i++ {
		arg := parts[i]
		if cvs.Options["languageID"] == "" {
			languageCode := languageMap[strings.ToUpper(arg)]
			if languageCode != "" {
				cvs.Options["languageID"] = languageCode
				continue
			}
		}
		if cvs.Options["indentType"] == "" {
			indentType := indentTypeMap[strings.ToUpper(arg)]
			if indentType != "" {
				cvs.Options["indentType"] = indentType
				continue
			}
		}
		if cvs.Options["indentUnit"] == "" {
			indentUnit := indentUnitMap[arg]
			if indentUnit != "" {
				cvs.Options["indentUnit"] = indentUnit
				continue
			}
		}
	}
	if cvs.Options["languageID"] == "" {
		cvs.Options["languageID"] = "71"
	}
	if cvs.Options["indentType"] == "" {
		cvs.Options["indentType"] = "space"
	}
	if cvs.Options["indentUnit"] == "" {
		cvs.Options["indentUnit"] = "4"
	}
}

func (cvs *CodeVServer) ExecuteCommand(c *gin.Context) {
	languageID := cvs.Options["languageID"]
	indentType := cvs.Options["indentType"]
	indentUnit := cvs.Options["indentUnit"]
	urlStr := fmt.Sprintf("/codev?languageID=%s&indentType=%s&indentUnit=%s", languageID, indentType, indentUnit)
	c.Redirect(http.StatusFound, urlStr)
}

func (bs *BlogServer) ParseCommand(stdin string) {
	bs.Options = make(map[string]string)
	rawParts := strings.Split(stdin, " ")
	parts := utils.RemoveElements(rawParts, "").([]string)
	for i := 0; i < len(parts); i++ {
		arg := parts[i]
		if bs.Options["author"] == "" && strings.ToUpper(arg) == "-A" {
			if i+1 < len(parts) {
				i++
				bs.Options["author"] = parts[i]
				continue
			}
		}
		if bs.Options["title"] == "" && strings.ToUpper(arg) == "-T" {
			if i+1 < len(parts) {
				i++
				bs.Options["title"] = parts[i]
				continue
			}
		}
	}
}

func (bs *BlogServer) ExecuteCommand(c *gin.Context) {
	urlStr := "/blogs"
	if bs.Options["author"] != "" {
		urlStr += "&authorName=" + bs.Options["author"]
	}
	if bs.Options["title"] != "" {
		if urlStr == "/blogs" {
			urlStr += "?title=" + bs.Options["title"]
		} else {
			urlStr += "&title=" + bs.Options["title"]
		}
	}
	c.Redirect(http.StatusFound, urlStr)
}

func (os *OpenServer) ParseCommand(stdin string) {
	os.Options = make(map[string]string)
	rawParts := strings.Split(stdin, " ")
	parts := utils.RemoveElements(rawParts, "").([]string)
	for i := 0; i < len(parts); i++ {
		arg := parts[i]
		if strings.ToUpper(arg) == "-N" {
			os.Options["mode"] = "newTab"
			continue
		}
		os.Options["url"] += arg
	}
}

func (os *OpenServer) ExecuteCommand(c *gin.Context) {
	urlStr := os.Options["url"]
	if urlStr == "" {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Please enter the URL or use man open to see how the 'open' command is used."})
		return
	}
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Unresolvable URL."})
		return
	}
	if parsedURL.Scheme != "" {
		urlStr = parsedURL.String()
	} else {
		urlStr = "http://" + parsedURL.String()
	}
	if os.Options["mode"] != "" {
		if parsedURL.RawQuery != "" {
			urlStr += "&mode=newTab"
		} else {
			urlStr += "?mode=newTab"
		}
	}
	c.JSON(http.StatusFound, gin.H{"response": urlStr})
}

func (gs *GoogleServer) ParseCommand(stdin string) {
	gs.Options = make(map[string]string)
	if stdin == "" {
		return
	} else {
		gs.Options["search"] = url.QueryEscape(stdin)
	}
}

func (gs *GoogleServer) ExecuteCommand(c *gin.Context) {
	urlStr := "https://www.google.com?mode=newTab"
	if gs.Options["search"] != "" {
		urlStr = "https://www.google.com/search?mode=newTab&q=" + gs.Options["search"]
	}
	c.JSON(http.StatusFound, gin.H{"response": urlStr})
}

func (bs *BingServer) ParseCommand(stdin string) {
	bs.Options = make(map[string]string)
	if stdin == "" {
		return
	} else {
		bs.Options["search"] = url.QueryEscape(stdin)
	}
}

func (bs *BingServer) ExecuteCommand(c *gin.Context) {
	urlStr := "https://www.bing.com?mode=newTab"
	if bs.Options["search"] != "" {
		urlStr = "https://www.bing.com/search?mode=newTab&q=" + bs.Options["search"]
	}
	c.JSON(http.StatusFound, gin.H{"response": urlStr})
}

func (gs *GitHubServer) ParseCommand(stdin string) {
	gs.Options = make(map[string]string)
	if stdin == "" {
		return
	} else {
		gs.Options["search"] = url.QueryEscape(stdin)
	}
}

func (gs *GitHubServer) ExecuteCommand(c *gin.Context) {
	urlStr := "https://github.com/?mode=newTab"
	if gs.Options["search"] != "" {
		urlStr = "https://github.com/search?mode=newTab&q=" + gs.Options["search"]
	}
	c.JSON(http.StatusFound, gin.H{"response": urlStr})
}
