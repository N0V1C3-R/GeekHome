package model

type FavoritesFolder struct {
	BaseModel
	UserId   int64  `gorm:"column:user_id;not null"`
	IsFolder bool   `gorm:"column:is_folder;not null"`
	Nickname string `gorm:"column:nickname;not null"`
	Upper    string `gorm:"column:upper;not null;default:'/'"`
	Path     string `gorm:"column:path;not null;default:''"`
}

func (*FavoritesFolder) TableName() string {
	return "favorites_folder"
}

func NewFavoritesFolder() *FavoritesFolder {
	return &FavoritesFolder{
		BaseModel: *NewBaseModel(),
	}
}
