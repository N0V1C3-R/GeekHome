package utils

import (
	"log"
	"net/http"
)

func HttpGet(url string) *http.Response {
	response, err := http.Get(url)
	if err != nil {
		return nil
	}
	if response.StatusCode != http.StatusOK {
		log.Fatalf("Status code error: %d %s", response.StatusCode, response.Status)
		return nil
	}
	return response
}

func GetCurrentCurrencyData() map[string]interface{} {
	resp := HttpGet("https://www.chinamoney.com.cn/chinese/bkccpr/")
	if resp == nil {
		return nil
	}
	return nil
}
