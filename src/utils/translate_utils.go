package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

const url = "https://api-free.deepl.com/v2/translate"

func TranslateText(text, targetLang string) *Response {
	AuthKey := os.Getenv("DEEPL_API_KEY")
	version := os.Getenv("VERSION")
	payload := bytes.NewBufferString(fmt.Sprintf("text=%s!&target_lang=%s", text, targetLang))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("DeepL-Auth-Key %s", AuthKey))
	req.Header.Set("User-Agent", fmt.Sprintf("WebHome/%s", version))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	response := parseResponse(resp)

	return &response
}

func parseResponse(response *http.Response) (res Response) {
	switch response.Status {
	case strconv.Itoa(http.StatusForbidden):
		res.StatusCode = http.StatusForbidden
		res.ResponseText = "Error: Blocked by CORS policy."
	case strconv.Itoa(http.StatusNotFound):
		res.StatusCode = http.StatusNotFound
		res.ResponseText = "Error: Server not found"
	case strconv.Itoa(http.StatusTooManyRequests):
		res.StatusCode = http.StatusTooManyRequests
		res.ResponseText = "Error: Too many requests, please try again later."
	case strconv.Itoa(456):
		res.StatusCode = http.StatusTooManyRequests
		res.ResponseText = "Error: The free translation quota is exceeded, please provide a new DeepL API key or upgrade DeepL subscription service."
	case strconv.Itoa(http.StatusInternalServerError):
		res.StatusCode = http.StatusInternalServerError
		res.ResponseText = "Error: DeepL server error, please try again later."
	default:
		res.StatusCode = response.StatusCode
	}

	var body map[string][]interface{}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		res.ResponseText = fmt.Sprintf("Error: %s", err)
	}
	res.ResponseText = body["translations"][0].(map[string]interface{})["text"].(string)
	return
}
