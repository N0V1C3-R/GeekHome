package command_server

import (
	"WebHome/src/database/dao"
	"WebHome/src/server/middleware"
	"WebHome/src/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type TranslateCommand struct {
	BaseCommand
}

var supportLanguages = map[string]string{
	"BG":    "Bulgarian",
	"CS":    "Czech",
	"DA":    "Danish",
	"DE":    "German",
	"EL":    "Greek",
	"EN":    "English(unspecified variant for backward compatibility; please select EN-GB or EN-US instead)",
	"EN-GB": "English(British)",
	"EN-US": "English(American)",
	"ES":    "Spanish",
	"ET":    "Estonian",
	"FI":    "Finish",
	"FR":    "French",
	"HU":    "Hungarian",
	"ID":    "Indonesian",
	"IT":    "Italian",
	"JA":    "Japanese",
	"KO":    "Korean",
	"LT":    "Lithuanian",
	"LV":    "Latvian",
	"NB":    "Norwegian(Bokmål)",
	"NL":    "Dutch",
	"PL":    "Polish",
	"PT":    "Portuguese",
	"PT-BR": "Portuguese(Brazilian)",
	"PT-PT": "Portuguese(all Portuguese varieties excluding Brazilian Portuguese)",
	"RO":    "Romanian",
	"RU":    "Russian",
	"SK":    "Slovak",
	"SL":    "Slovenian",
	"SV":    "Swedish",
	"TR":    "Turkish",
	"UK":    "Ukrainian",
	"ZH":    "Chinese(simplified)",
}

func (tc *TranslateCommand) ParseCommand(stdin string) {
	rawParts := strings.Split(stdin, " ")
	parts := utils.RemoveElements(rawParts, "").([]string)
	tc.Required = make(map[string]string)
	for i := 0; i < len(parts); i++ {
		arg := parts[i]
		if tc.Required["text"] == "" {
			switch strings.ToUpper(arg) {
			case "-TARGET", "-T":
				if i+1 < len(parts) {
					i++
					ok := validationSupportLanguage(strings.ToUpper(parts[i]))
					if !ok {
						return
					}
					tc.Required["targetLang"] = strings.ToUpper(parts[i])
				}
				continue
			}
		}
		tc.Required["text"] += arg + " "
	}
	tc.Required["text"] = strings.TrimRight(tc.Required["text"], " ")
	if tc.Required["targetLang"] == "" {
		tc.Required["targetLang"] = "EN-US"
	}
}

func (tc *TranslateCommand) ExecuteCommand(c *gin.Context) {
	ok, returnValue := getDeepLAPIKey(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"response": returnValue})
		return
	}

	targetLang := tc.Required["targetLang"]
	if targetLang == "" {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Unsupported target languages."})
		return
	}

	text := tc.Required["text"]
	if text == "" {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Please enter the text to be translated."})
		return
	}

	version := os.Getenv("VERSION")
	payload := bytes.NewBufferString(fmt.Sprintf("text=%s&target_lang=%s", text, targetLang))

	req, _ := http.NewRequest("POST", deepLAPIURL, payload)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("DeepL-Auth-Key %s", returnValue))
	req.Header.Set("User-Agent", fmt.Sprintf("WebHome/%s", version))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(resp.StatusCode, gin.H{"response": err.Error()})
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	response := parseTranslateResponse(resp)
	if response != nil {
		c.JSON(response.StatusCode, gin.H{"response": response.ResponseText})
	}
}

func getDeepLAPIKey(c *gin.Context) (bool, string) {
	key := os.Getenv("DEEPL_API_KEY")
	if key == "" {
		userAuth := middleware.GetUserAuth(c)
		if userAuth.UserId == 0 {
			return false, "ERROR: Failed to get the API Key."
		}
		key = getUserDeepLAPIKey(userAuth)
		if key == "" {
			return false, "ERROR: Please set a valid API Key for the user."
		}
	}
	return true, key
}

func getUserDeepLAPIKey(userAuth middleware.UserAuth) string {
	userApiKeyDao := dao.NewUserAPIKeyDao()
	entity := userApiKeyDao.GetAPIKey(userAuth.UserId, "DeepL", true)
	if entity.IsEnabled == false {
		return ""
	}
	key := utils.DecryptCipherText(entity.APIKey, userAuth.Username)
	return key
}

func validationSupportLanguage(targetLang string) bool {
	_, ok := supportLanguages[targetLang]
	return ok
}

func parseTranslateResponse(response *http.Response) *utils.Response {
	resp := utils.Response{
		StatusCode:   200,
		ResponseText: "",
	}
	switch response.Status {
	case strconv.Itoa(http.StatusForbidden):
		resp.StatusCode = http.StatusForbidden
		resp.ResponseText = "ERROR: Blocked by CORS policy."
	case strconv.Itoa(http.StatusNotFound):
		resp.StatusCode = http.StatusNotFound
		resp.ResponseText = "ERROR: Server not found"
	case strconv.Itoa(http.StatusTooManyRequests):
		resp.StatusCode = http.StatusTooManyRequests
		resp.ResponseText = "ERROR: Too many requests, please try again later."
	case strconv.Itoa(456):
		resp.StatusCode = http.StatusTooManyRequests
		resp.ResponseText = "ERROR: The free translation quota is exceeded, please provide a new DeepL API key or upgrade DeepL subscription service."
	case strconv.Itoa(http.StatusInternalServerError):
		resp.StatusCode = http.StatusInternalServerError
		resp.ResponseText = "ERROR: DeepL server error, please try again later."
	default:
		resp.StatusCode = response.StatusCode
	}

	var body map[string][]interface{}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		resp.ResponseText = fmt.Sprintf("ERROR: %s", err)
	}
	resp.ResponseText = body["translations"][0].(map[string]interface{})["text"].(string)
	return &resp
}
