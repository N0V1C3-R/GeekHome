package model

type UserAPIKey struct {
	BaseModel
	UserID     int64                 `gorm:"column:user_id;not null"`
	ServerName ThirdPartyServiceName `gorm:"column:server_name;not null"`
	APIKey     string                `gorm:"column:api_key;not null"`
	IsEnabled  bool                  `gorm:"column:is_enabled;not null"`
}

func (*UserAPIKey) TableName() string {
	return "user_api_key"
}

func NewUserApiKey() *UserAPIKey {
	return &UserAPIKey{
		BaseModel: *NewBaseModel(),
	}
}
