package command_server

import (
	"WebHome/src/database/dao"
	"WebHome/src/database/model"
	"WebHome/src/server/middleware"
	"WebHome/src/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type RenameServer struct {
	BaseCommand
}

type AddAPIKeyServer struct {
	BaseCommand
}

type FindAPIKeyServer struct {
	BaseCommand
}

type UpdateAPIKeyServer struct {
	BaseCommand
}

type BanAPIKeyServer struct {
	BaseCommand
}

type AliasCommandServer struct {
	BaseCommand
}

var supportServices = []model.ThirdPartyServiceName{
	model.ChatGPT,
	model.DeepL,
	model.Judge0,
}

func (rs *RenameServer) ParseCommand(stdin string) {
	rs.Options = make(map[string]string)
	rawParts := strings.Split(stdin, " ")
	parts := utils.RemoveElements(rawParts, "").([]string)
	for i := 0; i < len(parts); i++ {
		arg := parts[i]
		rs.Options["rename"] += arg + " "
	}
	rs.Options["rename"] = strings.TrimRight(rs.Options["rename"], " ")
}

func (rs *RenameServer) ExecuteCommand(c *gin.Context) {
	userAuth := middleware.GetUserAuth(c)
	if userAuth.UserId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"response": "ERROR: This command needs to be used while logged in."})
		return
	}
	newUserName := rs.Options["rename"]
	oldUsername := userAuth.Username
	if newUserName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"response": "ERROR: Please enter the user name to be changed."})
		return
	}
	if newUserName == oldUsername {
		c.JSON(http.StatusOK, gin.H{"response": "Username changed successfully."})
		return
	}
	ok := update3rdPartyAPIKey(oldUsername, newUserName, userAuth.UserId)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"response": "Username change failed. Please try again later."})
		return
	}
	userAuth.Username = newUserName
	c.SetCookie("userAuthorization", utils.SerializationObj(userAuth), 3600, "/", "", false, true)
	loginInfo := map[string]string{
		"Username": newUserName,
		"IP":       c.ClientIP(),
	}
	c.SetCookie("__userInfo", utils.SerializationObj(loginInfo), 3600, "/", "", false, false)
	c.JSON(http.StatusOK, gin.H{"response": "Username changed successfully."})
}

func (adk *AddAPIKeyServer) ParseCommand(stdin string) {
	adk.Options = make(map[string]string)
	rawParts := strings.Split(stdin, " ")
	parts := utils.RemoveElements(rawParts, "").([]string)
	for i := 0; i < len(parts); i++ {
		arg := parts[i]
		if adk.Options["serviceName"] == "" {
			adk.Options["serviceName"] = arg
			continue
		}
		if adk.Options["APIKey"] == "" {
			adk.Options["APIKey"] = arg
			continue
		}
	}
}

func (adk *AddAPIKeyServer) ExecuteCommand(c *gin.Context) {
	userAuth := middleware.GetUserAuth(c)
	if userAuth.UserId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"response": "ERROR: This command needs to be used while logged in."})
		return
	}
	serviceName := adk.Options["serviceName"]
	APIKey := adk.Options["APIKey"]
	if serviceName == "" {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Please specify the API service name to be added."})
		return
	}
	isValid, thirdPartyServiceName := validationSupportService(serviceName)
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"response": fmt.Sprintf("ERROR: Unsupported services: %s", serviceName)})
		return
	}
	ok := saveAPIKey(userAuth.UserId, thirdPartyServiceName, APIKey, userAuth.Username)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"response": "ERROR: Failed to add, please try again later."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": fmt.Sprintf("%s Service API added successfully", thirdPartyServiceName)})
}

func (fdk *FindAPIKeyServer) ParseCommand(stdin string) {
	fdk.Options = make(map[string]string)
	rawParts := strings.Split(stdin, " ")
	parts := utils.RemoveElements(rawParts, "").([]string)
	for i := 0; i < len(parts); i++ {
		arg := parts[i]
		if fdk.Options["serviceName"] == "" {
			fdk.Options["serviceName"] = arg
			continue
		}
	}
}

func (fdk *FindAPIKeyServer) ExecuteCommand(c *gin.Context) {
	userAuth := middleware.GetUserAuth(c)
	if userAuth.UserId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"response": "ERROR: This command needs to be used while logged in."})
		return
	}
	serviceName := fdk.Options["serviceName"]
	if serviceName == "" {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Please specify the API service name to be queried."})
		return
	}
	isValid, thirdPartyServiceName := validationSupportService(serviceName)
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"response": fmt.Sprintf("ERROR: Unsupported services: %s", serviceName)})
		return
	}
	APIKey := findAPIKey(userAuth.UserId, thirdPartyServiceName)
	if APIKey == "" {
		c.JSON(http.StatusOK, gin.H{"response": fmt.Sprintf("ERROR: Failed to query the available %s service corresponding API key.\nPlease use the \"adk\" command to add or \"upk\" command to update to available.", serviceName)})
		return
	}
	decryptValue := utils.DecryptCipherText(APIKey, userAuth.Username)
	c.JSON(http.StatusOK, gin.H{"response": fmt.Sprintf("The API key for the %s service is: %s", serviceName, decryptValue)})
}

func (upk *UpdateAPIKeyServer) ParseCommand(stdin string) {
	upk.Options = make(map[string]string)
	rawParts := strings.Split(stdin, " ")
	parts := utils.RemoveElements(rawParts, "").([]string)
	for i := 0; i < len(parts); i++ {
		arg := parts[i]
		if upk.Options["serviceName"] == "" {
			upk.Options["serviceName"] = arg
			continue
		}
		if upk.Options["APIKey"] == "" {
			upk.Options["APIKey"] = arg
			continue
		}
	}
}

func (upk *UpdateAPIKeyServer) ExecuteCommand(c *gin.Context) {
	userAuth := middleware.GetUserAuth(c)
	if userAuth.UserId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"response": "ERROR: This command needs to be used while logged in."})
		return
	}
	serviceName := upk.Options["serviceName"]
	APIKey := upk.Options["APIKey"]
	if serviceName == "" {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Please specify the API service name to be updated."})
		return
	}
	isValid, thirdPartyServiceName := validationSupportService(serviceName)
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"response": fmt.Sprintf("ERROR: Unsupported 3rd-party services: %s", serviceName)})
		return
	}
	ok := updateAPIKey(userAuth.UserId, thirdPartyServiceName, APIKey, userAuth.Username)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Cannot find the service, please use the \"adk\" command to add the service."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": fmt.Sprintf("%s Service API updated successfully", thirdPartyServiceName)})
}

func (ban *BanAPIKeyServer) ParseCommand(stdin string) {
	ban.Options = make(map[string]string)
	rawParts := strings.Split(stdin, " ")
	parts := utils.RemoveElements(rawParts, "").([]string)
	for i := 0; i < len(parts); i++ {
		arg := parts[i]
		if ban.Options["serviceName"] == "" {
			ban.Options["serviceName"] = arg
			continue
		}
		if ban.Options["APIKey"] == "" {
			ban.Options["APIKey"] = arg
			continue
		}
	}
}

func (ban *BanAPIKeyServer) ExecuteCommand(c *gin.Context) {
	userAuth := middleware.GetUserAuth(c)
	if userAuth.UserId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"response": "ERROR: This command needs to be used while logged in."})
		return
	}
	serviceName := ban.Options["serviceName"]
	if serviceName == "" {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Please specify the API service name to be updated."})
		return
	}
	isValid, thirdPartyServiceName := validationSupportService(serviceName)
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"response": fmt.Sprintf("ERROR: Unsupported services: %s", serviceName)})
		return
	}
	ok := disableAPIKey(userAuth.UserId, thirdPartyServiceName)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Cannot find the service, please use the adk command to add the service."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": fmt.Sprintf("%s Service API disabled successfully", thirdPartyServiceName)})
}

func (acs *AliasCommandServer) ParseCommand(stdin string) {

}

func (acs *AliasCommandServer) ExecuteCommand(c *gin.Context) {

}

func update3rdPartyAPIKey(oldName, newName string, userId int64) bool {
	var flag bool
	userEntityDao := dao.NewUserEntityDao()
	users := userEntityDao.FindUserListByUserId([]int64{userId})
	if len(users) == 0 {
		return false
	}
	user := users[0]
	userAPIKeyDao := dao.NewUserAPIKeyDao()
	userAPIKeyModels := userAPIKeyDao.GetUserAllAPIKeys(userId)
	tx := userAPIKeyDao.DB.Begin()
	for i := 0; i < 3; i++ {
		for _, userAPIKeyModel := range userAPIKeyModels {
			APIKey := userAPIKeyModel.APIKey
			updatedAt := userAPIKeyModel.UpdatedAt
			decryptValue := utils.DecryptCipherText(APIKey, oldName)
			userAPIKeyModel.APIKey = utils.EncryptPlainText([]byte(decryptValue), newName)
			userAPIKeyModel.UpdatedAt = utils.ConvertToMilliTime(utils.GetCurrentTime())
			if err := tx.Updates(userAPIKeyModel).Error; err != nil {
				flag = false
				tx.Rollback()
				userAPIKeyModel.APIKey = APIKey
				userAPIKeyModel.UpdatedAt = updatedAt
				break
			}
		}
		user.Username = newName
		updatedAt := user.UpdatedAt
		user.UpdatedAt = utils.ConvertToMilliTime(utils.GetCurrentTime())
		if err := tx.Updates(user).Error; err != nil {
			flag = false
			tx.Rollback()
			user.Username = oldName
			user.UpdatedAt = updatedAt
		} else {
			flag = true
			tx.Commit()
			break
		}
	}
	return flag
}

func validationSupportService(serviceName string) (bool, model.ThirdPartyServiceName) {
	for _, supportService := range supportServices {
		if strings.EqualFold(string(supportService), serviceName) {
			return true, supportService
		}
	}
	return false, ""
}

func saveAPIKey(userId int64, serviceName model.ThirdPartyServiceName, APIKey, username string) bool {
	var ok bool
	userAPIKeyDao := dao.NewUserAPIKeyDao()
	entity := userAPIKeyDao.GetAPIKey(userId, serviceName, false)
	if entity.Id != 0 {
		entity.APIKey = utils.EncryptPlainText([]byte(APIKey), username)
		entity.IsEnabled = true
		entity.UpdatedAt = utils.ConvertToMilliTime(utils.GetCurrentTime())
		userAPIKeyDao.Updates(entity)
		return true
	}
	for i := 0; i < 3; i++ {
		ok = userAPIKeyDao.CreateUserAPIKey(userId, serviceName, APIKey, username)
		if ok {
			return ok
		}
	}
	return ok
}

func findAPIKey(userId int64, serviceName model.ThirdPartyServiceName) string {
	userAPIKeyDao := dao.NewUserAPIKeyDao()
	entity := userAPIKeyDao.GetAPIKey(userId, serviceName, true)
	return entity.APIKey
}

func updateAPIKey(userId int64, serviceName model.ThirdPartyServiceName, APIKey, username string) bool {
	userAPIKeyDao := dao.NewUserAPIKeyDao()
	entity := userAPIKeyDao.GetAPIKey(userId, serviceName, false)
	if entity.Id != 0 {
		entity.APIKey = utils.EncryptPlainText([]byte(APIKey), username)
		entity.IsEnabled = true
		entity.UpdatedAt = utils.ConvertToMilliTime(utils.GetCurrentTime())
		err := userAPIKeyDao.SingleUpdate(entity)
		if err != nil {
			return false
		}
		return true
	}
	return false
}

func disableAPIKey(userId int64, serviceName model.ThirdPartyServiceName) bool {
	userAPIKeyDao := dao.NewUserAPIKeyDao()
	entity := userAPIKeyDao.GetAPIKey(userId, serviceName, true)
	if entity.Id != 0 {
		entity.IsEnabled = false
		entity.UpdatedAt = utils.ConvertToMilliTime(utils.GetCurrentTime())
		err := userAPIKeyDao.SingleUpdate(entity)
		if err != nil {
			return false
		}
		return true
	}
	return false
}
