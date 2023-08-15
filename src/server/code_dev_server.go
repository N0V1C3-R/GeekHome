package server

import (
	"WebHome/src/database/dao"
	"WebHome/src/server/middleware"
	"WebHome/src/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
	"net/http"
	"os"
)

type submissionRequest struct {
	Stdin      string `json:"stdin"`
	LanguageID int    `json:"language_id"`
	SourceCode string `json:"source_code"`
}

type submissionToken struct {
	Token string `json:"token"`
}

type resultRequest struct {
	Stdout        string      `json:"stdout"`
	Time          string      `json:"time"`
	Memory        int         `json:"memory"`
	Stderr        string      `json:"stderr"`
	Token         string      `json:"token"`
	CompileOutput string      `json:"compile_output"`
	Message       string      `json:"message"`
	Status        interface{} `json:"status"`
}

var languageMap = map[string]string{
	"1001": "C (Clang 10.0.1)",
	"50":   "C (GCC 9.2.0)",
	"1022": "C# (Mono 6.10.0.104)",
	"1021": "C# (.NET Core SDK 3.1.302)",
	"1023": "C# Test (.Net Core SDK 3.1.302, NUnit 3.12.0)",
	"1002": "C++ (Clang 10.0.1)",
	"54":   "C++ (GCC 9.2.0)",
	"1015": "C++ Test (Clang 10.0.1, Google Test 1.8.1)",
	"1012": "C++ Test (GCC 8.4.0, Google Test1.8.1)",
	"60":   "Go(1.13.5)",
	"1004": "Java (OpenJDK 14.0.1)",
	"63":   "JavaScript (Node.js 12.14.0)",
	"79":   "Objective-C (Clang 7.0.1)",
	"68":   "PHP (7.4.1)",
	"70":   "Python (2.7.17)",
	"71":   "Python (3.8.1)",
	"80":   "R (4.0.0)",
	"72":   "Ruby (2.7.0)",
	"73":   "Rust (1.40.0)",
	"83":   "Swift (5.2.3)",
	"74":   "TypeScript (3.7.4)",
}

var indentTypeMap = map[string]string{
	"tab":   "Tab",
	"space": "Space",
}

var indentUnitMap = map[string]string{
	"2": "2",
	"4": "4",
}

func CodeDevGetHandle(c *gin.Context) {
	token, err := utils.GetToken(secretKey)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.SetCookie(cookieName, token, cookieMaxAge, "/", "", false, true)
	languageID := c.Query("languageID")
	indentType := c.Query("indentType")
	indentUnit := c.Query("indentUnit")
	languageHTML := generateLanguageSelectHTML(languageID)
	indentTypeHTML := generateIndentTypeSelectHTML(indentType)
	indentUnitHTML := generateIndentUnitSelectHTML(indentUnit)
	c.HTML(http.StatusOK, "codedev.html",
		gin.H{
			"title":          "CodeV",
			"languageHTML":   template.HTML(languageHTML),
			"indentTypeHTML": template.HTML(indentTypeHTML),
			"indentUnitHTML": template.HTML(indentUnitHTML),
		})
}

func CodeDevPostHandle(c *gin.Context) {
	postJudge0API(c)
}

func generateLanguageSelectHTML(languageID string) string {
	var buffer bytes.Buffer
	if languageID == "" {
		languageID = "71"
	}

	for key, val := range languageMap {
		buffer.WriteString(fmt.Sprintf(`<option value="%s"`, key))
		if key == languageID {
			buffer.WriteString(` selected`)
		}
		buffer.WriteString(fmt.Sprintf(`>%s</option>`, val))

	}

	return buffer.String()
}

func generateIndentTypeSelectHTML(indentType string) string {
	var buffer bytes.Buffer
	if indentType == "" {
		indentType = "space"
	}

	for key, val := range indentTypeMap {
		buffer.WriteString(fmt.Sprintf(`<option value="%s"`, key))
		if key == indentType {
			buffer.WriteString(` selected`)
		}
		buffer.WriteString(fmt.Sprintf(`>%s</option>`, val))

	}

	return buffer.String()
}

func generateIndentUnitSelectHTML(indentUnit string) string {
	var buffer bytes.Buffer
	if indentUnit == "" {
		indentUnit = "4"
	}

	for key, val := range indentUnitMap {
		buffer.WriteString(fmt.Sprintf(`<option value="%s"`, key))
		if key == indentUnit {
			buffer.WriteString(` selected`)
		}
		buffer.WriteString(fmt.Sprintf(`>%s</option>`, val))

	}

	return buffer.String()
}

func postJudge0API(c *gin.Context) {
	reqBody, _ := io.ReadAll(c.Request.Body)
	reqMap := &codeMsg{}
	err := json.Unmarshal(reqBody, reqMap)
	if err != nil {
		panic(err)
	}
	inputValue := reqMap.InputValue
	languageId := reqMap.LanguageID
	sourceCode := reqMap.SourceCode

	authKey := os.Getenv("JUDGE0_API_KEY")
	if authKey == "" {
		userAuth := middleware.GetUserAuth(c)
		if userAuth.UserId == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"response": "ERROR: Failed to get user information."})
			return
		}
		authKey = parseJudge0APIKey(authKey, userAuth, c)
		if authKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"response": "ERROR: Failed to get the API Key."})
			return
		}
	}

	if sourceCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"response": "ERROR: Source code can't be empty!"})
		return
	}

	if (inputValue != "" && utils.IsBase64String(inputValue)) || utils.IsBase64String(sourceCode) {
		c.JSON(http.StatusBadRequest, gin.H{"response": "ERROR: Please enter the correct parameters."})
		return
	}

	token := requestToken(inputValue, sourceCode, authKey, languageId)
	requestResult(token, authKey, c)
}

func parseJudge0APIKey(authKey string, userAuth middleware.UserAuth, c *gin.Context) string {
	userAPIKeyDao := dao.NewUserAPIKeyDao()
	entity := userAPIKeyDao.GetAPIKey(userAuth.UserId, "Judge0", true)
	if entity.IsEnabled == false {
		return ""
	}
	authKey = utils.DecryptCipherText(entity.APIKey, userAuth.Username)
	if authKey == "" {
		return ""
	}
	return authKey
}

func requestToken(inputValue, sourceCode, authKey string, languageId int) string {
	submissionRequest := submissionRequest{
		Stdin:      utils.Base64EncodeString([]byte(inputValue)),
		LanguageID: languageId,
		SourceCode: utils.Base64EncodeString([]byte(sourceCode)),
	}
	requestBody, err := json.Marshal(submissionRequest)

	payload := bytes.NewBuffer(requestBody)

	req, _ := http.NewRequest("POST", judge0TokenUrl, payload)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-RapidAPI-Key", authKey)
	req.Header.Add("X-RapidAPI-Host", "judge0-ce.p.rapidapi.com")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	var submissionToken submissionToken
	err = json.NewDecoder(resp.Body).Decode(&submissionToken)

	return submissionToken.Token
}

func requestResult(token, authKey string, c *gin.Context) {
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"response": "ERROR: Illegal token."})
	}
	req, _ := http.NewRequest("GET", judge0API+token, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-RapidAPI-Key", authKey)
	req.Header.Add("X-RapidAPI-Host", "judge0-ce.p.rapidapi.com")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"response": "ERROR: Bad requests."})
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	var resultRequest resultRequest
	err = json.NewDecoder(resp.Body).Decode(&resultRequest)
	if err != nil {

	}

	c.JSON(http.StatusOK, gin.H{"response": resultRequest})
}
