package redis

import (
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
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

var (
	addr     string
	password string
)

func ConnectionRedis() *redis.Client {
	getConnectionInfo()
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	return client
}

func getConnectionInfo() {
	addr = os.Getenv("REDIS_POINT")
	password = os.Getenv("REDIS_PWD")
}
