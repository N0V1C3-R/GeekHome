package dao

import (
	"WebHome/src/database"
	"gorm.io/gorm"
	"sync"
)

var (
	once            sync.Once
	roDb            *gorm.DB
	db              *gorm.DB
	baseDao         *BaseDao
	readOnlyBaseDao *BaseDao
)

func init() {
	once.Do(
		func() {
			roDb, _ = database.ConnectDB()
			db, _ = database.ConnectDB()
			baseDao = NewBaseDao()
			readOnlyBaseDao = NewReadOnlyBaseDao()
		})
}
