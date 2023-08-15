package dao

import (
	"gorm.io/gorm"
)

type BaseDao struct {
	*gorm.DB
}

func NewBaseDao() *BaseDao {
	return &BaseDao{db}
}

func NewReadOnlyBaseDao() *BaseDao {
	return &BaseDao{roDb}
}

func (dao *BaseDao) GetOne(model interface{}, filter interface{}) interface{} {
	result := dao.Where(filter).Find(model)
	if result.RowsAffected == 0 {
		return nil
	}
	return result
}

func (dao *BaseDao) SingleInsert(model interface{}) error {
	result := dao.Create(model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (dao *BaseDao) SingleUpdate(model interface{}) error {
	result := dao.Save(model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (dao *BaseDao) BatchInsert(data map[string]interface{}) error {
	tx := dao.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, value := range data {
		result := tx.Create(value)
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	return tx.Commit().Error
}
