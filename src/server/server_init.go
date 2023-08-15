package server

import (
	"WebHome/src/database/model"
	redisC "WebHome/src/redis"
	"context"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type registerForm struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type emailVerifyForm struct {
	Code string `json:"code"`
}

type userCookie struct {
	UserId   int64  `json:"uid"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type loginInfo struct {
	Username string `json:"username"`
	IP       string `json:"IP"`
}

type codeMsg struct {
	InputValue string `json:"inputValue"`
	LanguageID int    `json:"languageID"`
	SourceCode string `json:"sourceCode"`
}

var (
	once                sync.Once
	rdb                 *redis.Client
	ctx                 context.Context
	blogImagePath       string
	judge0API           string
	judge0TokenUrl      string
	blogClassifications []model.BlogClassification
	secretKey           = "123"
	username            = "Visitor"
)

const (
	cookieName   = "token"
	cookieMaxAge = 3600
)

func init() {
	once.Do(
		func() {
			_, file, _, _ := runtime.Caller(0)
			_ = os.Chdir(filepath.Dir(file))
			var configPath string
			if os.Getenv("ENVIRONMENT") == "local" {
				configPath = filepath.Join("..", "..", "config", ".env_local")
			} else {
				configPath = filepath.Join("..", "..", "config", ".env")
			}
			_ = godotenv.Load(configPath)
			rdb = redisC.ConnectionRedis()
			ctx = context.Background()
			blogImagePath = os.Getenv("BLOG_IMAGE_PATH")
			judge0API = os.Getenv("JUDGE0_API_URL")
			judge0TokenUrl = os.Getenv("JUDGE0_TOKEN_URL")
			blogClassifications = []model.BlogClassification{
				"Programming Languages",
				"Operating System",
				"Database",
				"News",
				"Network Security",
				"NMiscellaneous Discussions",
			}
		})
}
