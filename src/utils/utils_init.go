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
	configPath := filepath.Join("..", "..", "config", ".env_local")
	_ = godotenv.Load(configPath)
}
