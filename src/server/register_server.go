package server

import (
	"WebHome/src/database/dao"
	"WebHome/src/database/model"
	"WebHome/src/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

var (
	registerWg    sync.WaitGroup
	userEntityDao dao.UserEntityDao
)

type registerRedis struct {
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	VerifyCode string `json:"verifyCode"`
}

func RegisterHandle(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		token, err := utils.GetToken(secretKey)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.SetCookie(cookieName, token, cookieMaxAge, "/", "", false, true)
		c.HTML(http.StatusOK, "register.html", gin.H{"title": "Register"})
	} else if c.Request.Method == http.MethodPost {
		postRegisterHandle(c)
	} else {
		c.AbortWithStatus(http.StatusMethodNotAllowed)
	}
}

func postRegisterHandle(c *gin.Context) {
	token, err := c.Cookie(cookieName)
	if err != nil {
		_ = c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	if ok, err := utils.VerifyToken(token, secretKey); !ok {
		_ = c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	requestID := c.Request.Header.Get("X-Request-ID")
	var registerForm registerForm
	if err := c.ShouldBind(&registerForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"response": err.Error()})
		return
	}

	r := checkRequest(requestID, registerForm)

	if r == nil {
		c.JSON(http.StatusBadRequest, gin.H{"response": "Exception request, please try again."})
		return
	}
	email := registerForm.Email
	userEntityDao = *dao.NewUserEntityDao()
	terminate := make(chan bool)
	var code string
	status := make(chan error)

	go func(email string) {
		user := userEntityDao.IsUserExists(email)
		if user.Email != email {
			terminate <- true
		} else {
			terminate <- false
		}
	}(email)

	if <-terminate {
		if r.VerifyCode != "" {
			code = r.VerifyCode
			c.JSON(http.StatusOK, gin.H{"response": "OK"})
			return
		} else {
			code = utils.GenerateVerificationCode()
			go sendRegisterVerificationEmail(email, code, status)
			emailStatus := <-status
			if emailStatus != nil {
				c.JSON(http.StatusBadRequest, gin.H{"response": "Mail delivery failure."})
				return
			} else {
				r.VerifyCode = code
				expiration := 10 * time.Minute
				updateRegisterRedis(requestID, r, expiration)
				c.JSON(http.StatusOK, gin.H{"response": "OK"})
				return
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"response": "This email address is already registered, please change your email address."})
		return
	}
}

func checkRequest(requestID string, registerForm registerForm) *registerRedis {
	var r registerRedis
	value, err := rdb.Get(ctx, requestID).Result()
	if err != nil {
		r.Email = registerForm.Email
		r.Username = registerForm.Username
		r.Password = registerForm.Password
		jsonData, err := json.Marshal(r)
		if err != nil {
			return nil
		}
		expiration := 10 * time.Minute
		err = rdb.Set(ctx, requestID, jsonData, expiration).Err()
		return &r
	} else {
		err = json.Unmarshal([]byte(value), &r)
	}
	if err != nil {
		return nil
	} else {
		return &r
	}
}

func updateRegisterRedis(requestID string, r *registerRedis, expiration time.Duration) bool {
	jsonData, err := json.Marshal(r)
	if err != nil {
		return false
	}
	err = rdb.Set(ctx, requestID, jsonData, expiration).Err()
	if err != nil {
		return false
	}
	return true
}

func sendRegisterVerificationEmail(email string, code string, status chan error) {
	topic := "Registration Verification"
	body := "Verify Code: " + code + ". The verification code is valid for ten minutes."
	status <- utils.SysSendEmail(email, topic, body)
}

func RegisterUser(c *gin.Context) {
	requestID := c.Request.Header.Get("X-Request-ID")

	registerWg.Add(1)
	go func() {
		defer registerWg.Done()

		var registerForm registerForm
		r := checkRequest(requestID, registerForm)
		if r == nil {
			c.JSON(419, gin.H{"response": "Request Expired."})
			return
		}

		username := r.Username
		email := r.Email
		password := r.Password
		userEntity := (*model.NewUserEntity()).CreateUser(username, email, password)
		userEntityDao.CreateClientUser(userEntity)
		_, _ = rdb.Del(ctx, requestID).Result()
		c.JSON(http.StatusFound, gin.H{"response": "Registration is successful, will jump to the login screen soon."})
	}()

	registerWg.Wait()
}
