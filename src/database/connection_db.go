package database

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

func ConnectDB() (*gorm.DB, error) {
	var (
		dataSource string
		dialer     gorm.Dialector
	)
	dataSource = getDataSource()
	dialer = mysql.Open(dataSource)
	db, err := gorm.Open(dialer, &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	return db, nil
}

func getDataSource() string {
	var (
		dbPoint    string = os.Getenv("MYSQL_DB_POINT")
		dbUsername string = os.Getenv("MYSQL_GEEK_HOME_USER")
		dbPassword string = os.Getenv("MYSQL_GEEK_HOME_PWD")
		dbSchema   string = os.Getenv("MYSQL_DB_SCHEMA")
	)
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4", dbUsername, dbPassword, dbPoint, dbSchema)
	return dataSource
}

func getReadOnlyDataSource() string {
	var (
		dbPoint    string = os.Getenv("MYSQL_DB_POINT")
		dbUsername string = os.Getenv("MYSQL_GEEK_HOME_READ_ONLY_USER")
		dbPassword string = os.Getenv("MYSQL_GEEK_HOME_READ_ONLY_PWD")
		dbSchema   string = os.Getenv("MYSQL_DB_SCHEMA")
	)
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4", dbUsername, dbPassword, dbPoint, dbSchema)
	return dataSource
}
