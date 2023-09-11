package main

import (
	"WebHome/src/server"
	"WebHome/src/server/command_server"
	"WebHome/src/server/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"html/template"
	"os"
	"path/filepath"
	"runtime"
	"time"
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

func main() {
	r := gin.Default()
	r.Use(middleware.LogMiddleware())
	r.Use(middleware.UpdateAuthMiddleware())
	r.MaxMultipartMemory = 2 << 20
	templateFuncRegister(r)
	routingRegistration(r)
	_ = r.Run(":5466")
}

func routingRegistration(router *gin.Engine) {
	_, file, _, _ := runtime.Caller(0)
	_ = os.Chdir(filepath.Dir(file))
	fmt.Println(file)
	router.LoadHTMLGlob("../templates/html/*")
	router.StaticFile("/favicon.ico", "/favicon.ico")
	router.Static("/templates", "../templates")

	router.NoRoute(server.PageNotFound)

	registerGroupRegistration(router)
	blogGroupRegistration(router)

	router.GET("/", server.HomeHandle)
	router.GET("/hello", server.WelcomeHandle)
	router.Any("/login", server.LoginHandle)
	router.POST("/logout", server.LogoutHandle)
	router.POST("/verify", server.VerifyCode)
	router.GET("/codev", server.CodeDevGetHandle)
	router.POST("/codev", server.CodeDevPostHandle)
	router.POST("/command", command_server.CommandHandle)
}

func registerGroupRegistration(router *gin.Engine) {
	registerGroup := router.Group("/register")
	{
		registerGroup.Any("", server.RegisterHandle)
		registerGroup.POST("/verify_code", server.EmailVerify)
		registerGroup.POST("/create_user", server.RegisterUser)
	}
}

func blogGroupRegistration(router *gin.Engine) {
	blogGroup := router.Group("/blogs")
	{
		blogGroup.GET("", server.BlogListHandler)
		blogGroup.GET("/read/*title", server.ReadHandle)
		blogGroup.GET("/img/:guid", server.LoadLocalImageHandle)
		blogGroup.GET("../img/:guid", server.LoadLocalImageHandle)
		blogGroup.GET("/new", server.EditArticle)
		blogGroup.GET("/edit/*title", server.EditArticle)
		blogGroup.POST("/save/*title", server.SaveArticle)
		blogGroup.POST("/upload_image", server.UploadImageHandle)
	}
}

func templateFuncRegister(router *gin.Engine) {
	router.SetFuncMap(template.FuncMap{
		"formatTimestamp": func(ts int64) string {
			location, _ := time.LoadLocation("Asia/Shanghai")
			t := time.UnixMilli(ts).In(location)
			return t.Format("2006-01-02 15:04:05")
		},
		"getUsernameById": func(userMap map[int64]string, id int64) string {
			return userMap[id]
		},
	})
}
