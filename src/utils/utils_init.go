package utils

import (
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"runtime"
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
}
