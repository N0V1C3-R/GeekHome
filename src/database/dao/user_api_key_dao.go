package dao

import (
	"WebHome/src/database/model"
	"WebHome/src/utils"
	"gorm.io/gorm"
)

type UserAPIKeyDao struct {
	BaseDao
	Schema *gorm.DB
}

func NewUserAPIKeyDao() *UserAPIKeyDao {
	schema := db.Table(model.NewUserApiKey().TableName())
	return &UserAPIKeyDao{*baseDao, schema}
}

func (dao *UserAPIKeyDao) GetAPIKey(userId int64, serviceName model.ThirdPartyServiceName, isUsable bool) model.UserAPIKey {
	var conditions string
	if isUsable {
		conditions = "user_id = ? AND server_name = ? AND is_enabled = 1 AND deleted_at = 0"
	} else {
		conditions = "user_id = ? AND server_name = ? AND deleted_at = 0"
	}
	values := []interface{}{userId, serviceName}
	var UserApiKeyModel model.UserAPIKey
	err := dao.Schema.Where(conditions, values...).First(&UserApiKeyModel)
	if err.Error != nil {
		return UserApiKeyModel
	}
	return UserApiKeyModel
}

func (dao *UserAPIKeyDao) CreateUserAPIKey(userId int64, serverName model.ThirdPartyServiceName, apiKey, username string) bool {
	userApiKeyModel := model.NewUserApiKey()
	userApiKeyModel.UserID = userId
	userApiKeyModel.ServerName = serverName
	userApiKeyModel.APIKey = utils.EncryptPlainText([]byte(apiKey), username)
	userApiKeyModel.IsEnabled = true
	err := dao.SingleInsert(&userApiKeyModel)
	if err != nil {
		return false
	}
	return true
}

func (dao *UserAPIKeyDao) GetUserAllAPIKeys(userId int64) []model.UserAPIKey {
	conditions := "user_id = ? AND deleted_at = 0"
	values := []interface{}{userId}
	var userAPIKeyModels []model.UserAPIKey
	err := dao.Schema.Where(conditions, values...).Find(&userAPIKeyModels)
	if err.Error != nil {
		return []model.UserAPIKey{}
	}
	return userAPIKeyModels
}
