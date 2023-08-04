package dao

import (
	"WebHome/src/database"
	"gorm.io/gorm"
	"sync"
)

var (
	once    sync.Once
	db      *gorm.DB
	baseDao *BaseDao
)

func init() {
	once.Do(
		func() {
			db, _ = database.ConnectDB(database.Mysql)
			baseDao = NewBaseDao()
		})
}
