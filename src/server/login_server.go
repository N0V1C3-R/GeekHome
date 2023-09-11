package server

import (
	"WebHome/src/database/dao"
	"WebHome/src/database/model"
	"WebHome/src/utils"
	"encoding/binary"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

func LoginHandle(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		token, err := utils.GetToken(secretKey)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.SetCookie(cookieName, token, cookieMaxAge, "/", "", false, true)
		c.HTML(http.StatusOK, "login.html", gin.H{"title": "Login"})
	} else if c.Request.Method == http.MethodPost {
		postLoginHandle(c)
	} else {
		c.AbortWithStatus(http.StatusMethodNotAllowed)
	}
}

func postLoginHandle(c *gin.Context) {
	token, err := c.Cookie(cookieName)
	if err != nil {
		fmt.Println(err)
		_ = c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	if ok, err := utils.VerifyToken(token, secretKey); !ok {
		fmt.Println(err)
		_ = c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	var loginForm loginForm
	if err := c.ShouldBind(&loginForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	email := loginForm.Email
	password := loginForm.Password

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		userEntity := verifyLogin(email, password)
		if userEntity.Email != email {
			c.JSON(http.StatusBadRequest, gin.H{"response": "Incorrect username or password"})
			return
		}
		flushLastLoginTime(userEntity.Id)
		role := make([]byte, 4)
		binary.LittleEndian.PutUint32(role, uint32(userEntity.Role))
		userAuth := userCookie{
			UserId:      userEntity.Id,
			Username:    userEntity.Username,
			Role:        utils.EncryptPlainText(role, userEntity.Username),
			WorkingPath: "/",
		}
		loginInfo := loginInfo{
			Username: userEntity.Username,
			IP:       c.ClientIP(),
		}
		c.SetCookie("userAuthorization", utils.SerializationObj(userAuth), 3600, "/", "", true, true)
		c.SetCookie("__userInfo", utils.SerializationObj(loginInfo), 3600, "/", "", false, false)
		c.SetSameSite(http.SameSiteStrictMode)
		c.JSON(http.StatusFound, gin.H{"response": "Login successful!"})
	}()

	wg.Wait()
}

func verifyLogin(email, password string) model.UserEntity {
	encryptedPassword := utils.EncryptString(password)
	userDao := dao.NewUserEntityDao()
	result := userDao.GetUser(email, encryptedPassword)
	return result
}

func flushLastLoginTime(id int64) {
	userDao := dao.NewUserEntityDao()
	userDao.UpdateUserLoginTime(id)
}
