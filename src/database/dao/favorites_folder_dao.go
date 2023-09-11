package dao

import (
	"WebHome/src/database/model"
	"gorm.io/gorm"
)

type FavoritesFolderDao struct {
	BaseDao
	Schema *gorm.DB
}

func NewFavoritesFolderDao() *FavoritesFolderDao {
	schema := db.Table(model.NewFavoritesFolder().TableName()).Session(&gorm.Session{})
	return &FavoritesFolderDao{*baseDao, schema}
}

func (dao *FavoritesFolderDao) IsExist(userId int64, nickname, upper string) bool {
	var (
		condition string
		values    []interface{}
	)
	if nickname != "" {
		condition = "user_id = ? AND is_folder = 0 AND nickname = ? AND upper = ? AND deleted_at = 0"
		values = []interface{}{userId, nickname, upper}
	} else {
		condition = "user_id = ? AND is_folder = 0 AND upper = ? AND deleted_at = 0"
		values = []interface{}{userId, upper}
	}
	var count int64
	dao.Schema.Where(condition, values...).Count(&count)
	return count != 0
}

type FavoritesProfile struct {
	IsFolder bool
	Nickname string
	Path     string
}

func (dao *FavoritesFolderDao) GetRecordsByPath(userId int64, path string) []FavoritesProfile {
	var res []FavoritesProfile
	condition := "user_id = ? AND upper = ? AND deleted_at = 0"
	values := []interface{}{userId, path}
	dao.Schema.Where(condition, values...).Order("-is_folder").Find(&res)
	return res
}

func (dao *FavoritesFolderDao) GetRecordByNicknameAndUpper(userId int64, nickname, upper string) model.FavoritesFolder {
	var res model.FavoritesFolder
	condition := "user_id = ? AND is_folder = 0 AND upper = ? and nickname = ? AND deleted_at = 0"
	values := []interface{}{userId, upper, nickname}
	dao.Schema.Where(condition, values...).First(&res)
	return res
}

func (dao *FavoritesFolderDao) CreatRecord(userId int64, isFolder bool, nickName, upper, path string) bool {
	favoritesFolder := model.NewFavoritesFolder()
	favoritesFolder.UserId = userId
	favoritesFolder.IsFolder = isFolder
	favoritesFolder.Nickname = nickName
	favoritesFolder.Upper = upper
	favoritesFolder.Path = path
	err := dao.SingleInsert(&favoritesFolder)
	if err != nil {
		return false
	}
	return true
}

func (dao *FavoritesFolderDao) RemoveRecords(userId int64, nickname, upper string) {
	var (
		condition string
		values    []interface{}
	)
	if nickname == "" {
		condition = "user_id = ? AND upper LIKE ?"
		values = []interface{}{userId, upper}
	} else {
		condition = "user_id = ? AND nickname = ? AND upper LIKE ?"
		values = []interface{}{userId, nickname, upper}
	}
	var res []model.FavoritesFolder
	dao.Schema.Where(condition, values...).Find(&res)
	if len(res) != 0 {
		dao.Schema.Delete(&res)
	}
}
