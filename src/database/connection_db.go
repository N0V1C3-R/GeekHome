package database

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	_ = os.Chdir(filepath.Dir(file))
	configPath := filepath.Join("..", "..", "config", ".env_local")
	_ = godotenv.Load(configPath)
}

type DBType string

const (
	Mysql DBType = "mysql"
)

func ConnectDB(dbType DBType) (*gorm.DB, error) {
	var (
		dataSource string
		dialer     gorm.Dialector
	)
	switch dbType {
	case Mysql:
		dataSource = getDataSource(dbType)
		dialer = mysql.Open(dataSource)
	default:
		return nil, errors.New("unsupported db type")
	}
	db, err := gorm.Open(dialer, &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	return db, nil
}

func getDataSource(dbType DBType) string {
	var database = strings.ToUpper(string(dbType))
	var (
		dbPoint    string = os.Getenv(fmt.Sprintf("%s_DB_POINT", database))
		dbUsername string = os.Getenv(fmt.Sprintf("%s_DB_USERNAME", database))
		dbPassword string = os.Getenv(fmt.Sprintf("%s_DB_PASSWORD", database))
		dbSchema   string = os.Getenv(fmt.Sprintf("%s_DB_SCHEMA", database))
	)
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4", dbUsername, dbPassword, dbPoint, dbSchema)
	log.Println("-----datasource-----")
	log.Println(dataSource)
	log.Println("----------")
	return dataSource
}
