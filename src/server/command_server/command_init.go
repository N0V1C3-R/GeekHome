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
	Options map[string]string
}

var (
	deepLAPIURL string
	commands    map[string]reflect.Type
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	_ = os.Chdir(filepath.Dir(file))
	var configPath string
	if os.Getenv("ENVIRONMENT") == "local" {
		configPath = filepath.Join("..", "..", "config", ".env_local")
	} else {
		configPath = filepath.Join("..", "..", "config", ".env")
	}
	_ = godotenv.Load(configPath)

	deepLAPIURL = os.Getenv("DEEPL_URL")

	commands = map[string]reflect.Type{
		// Help
		"help": reflect.TypeOf(&HelpServer{}),

		// Man
		"man": reflect.TypeOf(&ManServer{}),

		// Translation
		"translate": reflect.TypeOf(&TranslateServer{}),
		"trans":     reflect.TypeOf(&TranslateServer{}),
		"ts":        reflect.TypeOf(&TranslateServer{}),

		// Exchange rate conversion
		"fx": reflect.TypeOf(&CurrencyConvertServer{}),
		"ex": reflect.TypeOf(&CurrencyConvertServer{}),

		// Morse code encoding
		"morse": reflect.TypeOf(&MorseServer{}),
		// Morse code decoding
		"esrom": reflect.TypeOf(&EsromServer{}),

		// Base64 conversion
		"base64": reflect.TypeOf(&Base64Server{}),
		"b64":    reflect.TypeOf(&Base64Server{}),

		// Base Conversion
		"bc": reflect.TypeOf(&BaseConversionServer{}),

		// Time to Timestamp
		"tts": reflect.TypeOf(&TimeConvertServer{}),
		// Timestamp to time
		"tst": reflect.TypeOf(&TimestampConvertServer{}),

		// Generate a password
		"genpwd": reflect.TypeOf(&GenpwdServer{}),

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

		// Access to the online IDE interface
		"codev": reflect.TypeOf(&CodeVServer{}),

		// Visit the blog
		"blogs": reflect.TypeOf(&BlogServer{}),

		// Open a link
		"open": reflect.TypeOf(&OpenServer{}),

		// Google search
		"google": reflect.TypeOf(&GoogleServer{}),
		"go":     reflect.TypeOf(&GoogleServer{}),

		// Bing search
		"bing": reflect.TypeOf(&BingServer{}),

		// GitHub search
		"github": reflect.TypeOf(&GitHubServer{}),
	}
}
