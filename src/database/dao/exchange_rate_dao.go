package dao

import (
	"WebHome/src/database/model"
	"gorm.io/gorm"
)

type ExchangeRateDao struct {
	BaseDao
	Schema *gorm.DB
}

func NewExchangeRateDao() *ExchangeRateDao {
	schema := db.Table(model.NewExchangeRate().TableName())
	return &ExchangeRateDao{*baseDao, schema}
}

func NewReadOnlyExchangeRateDao() *ExchangeRateDao {
	schema := db.Table(model.NewExchangeRate().TableName())
	return &ExchangeRateDao{*readOnlyBaseDao, schema}
}

func (dao *ExchangeRateDao) GetLastTradingDay() (LastTradingDay string) {
	err := dao.Schema.Select("max(trade_date)").Row().Scan(&LastTradingDay)
	if err != nil {
		return ""
	}
	return
}

func (dao *ExchangeRateDao) GetTheLatestUpdateTime() (updatedAt int64) {
	err := dao.Schema.Select("max(updated_at)").Row().Scan(&updatedAt)
	if err != nil {
		return
	}
	return
}

func (dao *ExchangeRateDao) FlushUpdateTime(updatedAt int64) {
	_ = dao.Schema.UpdateColumn("updated_at", updatedAt)
}
