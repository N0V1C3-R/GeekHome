package model

import (
	"WebHome/src/utils"
)

type UserEntity struct {
	BaseModel
	Username      string   `gorm:"column:username;not null"`
	Email         string   `gorm:"column:email;not null"`
	Password      string   `gorm:"column:password;not null"`
	Role          UserRole `gorm:"column:role;not null;default:1"`
	Active        bool     `gorm:"column:active;not null"`
	LastLoginTime int64    `gorm:"column:last_login_at;not null;default:0"`
}

func (*UserEntity) TableName() string {
	return "user_entity"
}

func NewUserEntity() *UserEntity {
	return &UserEntity{
		BaseModel: *NewBaseModel(),
	}
}

func (user *UserEntity) CreateUser(username, email, password string) UserEntity {
	user.Username = username
	user.Password = utils.EncryptString(password)
	user.Email = email
	user.Active = true
	user.LastLoginTime = 0
	return *user
}

func (user *UserEntity) InitPassword() string {
	password := utils.GeneratePassword(15, 12)
	user.Password = utils.EncryptString(password)
	user.UpdatedAt = utils.ConvertToMilliTime(utils.GetCurrentTime())
	return password
}

func (user *UserEntity) DeactivateUser() {
	user.Active = false
	user.UpdatedAt = utils.ConvertToMilliTime(utils.GetCurrentTime())
}
