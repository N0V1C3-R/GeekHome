package command_server

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
)

type CommandImpl interface {
	ParseCommand(parts string)
	ExecuteCommand(c *gin.Context)
}

type CommandRequest struct {
	Stdin string `json:"stdin"`
}

type BaseCommand struct {
	Options  map[string]string
	Required map[string]string
}

var (
	deepLAPIURL string
	commands    map[string]reflect.Type
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	_ = os.Chdir(filepath.Dir(file))
	configPath := filepath.Join("..", "..", "config", ".env_local")
	_ = godotenv.Load(configPath)

	deepLAPIURL = os.Getenv("DEEPL_URL")

	commands = map[string]reflect.Type{
		// Translation
		"translate": reflect.TypeOf(&TranslateCommand{}),
		"trans":     reflect.TypeOf(&TranslateCommand{}),
		"ts":        reflect.TypeOf(&TranslateCommand{}),

		// Exchange rate conversion
		"fx": reflect.TypeOf(&CurrencyConvertCommand{}),
		"ex": reflect.TypeOf(&CurrencyConvertCommand{}),

		// Morse code encoding
		"morse": reflect.TypeOf(&MorseServer{}),
		// Morse code decoding
		"esrom": reflect.TypeOf(&EsromServer{}),

		// Base64 conversion
		"base64": reflect.TypeOf(&Base64Server{}),
		"b64":    reflect.TypeOf(&Base64Server{}),

		// Time to Timestamp
		"tts": reflect.TypeOf(&TimeConvertServer{}),
		// Timestamp to time
		"tst": reflect.TypeOf(&TimestampConvertServer{}),

		// Access to the online IDE interface
		"codev": reflect.TypeOf(&CodeVServer{}),

		// Modify username
		"rename": reflect.TypeOf(&RenameServer{}),

		// Add 3rd-party service API
		"adk": reflect.TypeOf(&AddAPIKeyServer{}),
		// Find 3rd-party service API
		"fdk": reflect.TypeOf(&FindAPIKeyServer{}),
		// Update 3rd-party service API
		"upk": reflect.TypeOf(&UpdateAPIKeyServer{}),
		// Disable 3rd-party service API
		"ban": reflect.TypeOf(&BanAPIKeyServer{}),
	}
}
