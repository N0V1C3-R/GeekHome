package model

import (
	"WebHome/src/utils"
)

type BaseModel struct {
	Id        int64 `gorm:"primaryKey;column:id"`
	CreatedAt int64 `gorm:"column:created_at;not null"`
	UpdatedAt int64 `gorm:"autoUpdateTime:milli;column:updated_at;not null;default:0"`
	DeletedAt int64 `gorm:"column:deleted_at;not null;default:0"`
}

func NewBaseModel() *BaseModel {
	return &BaseModel{
		Id:        utils.GenerateSnowflake(),
		CreatedAt: utils.ConvertToMilliTime(utils.GetCurrentTime()),
		UpdatedAt: 0,
		DeletedAt: 0,
	}
}
