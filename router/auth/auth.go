package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/mojocn/base64Captcha"
	"go-blog/common"
	"go-blog/common/api"
	"go-blog/common/jwt"
	"go-blog/conf"
	"go-blog/service"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
)

type ConsoleAuth interface {
	Register(*gin.Context)
	AuthRegister(*gin.Context)
	Login(*gin.Context)
	AuthLogin(*gin.Context)
	Logout(*gin.Context)
	DelCache(*gin.Context)
}

type Auth struct {
}

func NewAuth() ConsoleAuth {
	return &Auth{}
}

// customizeRdsStore An object implementing Store interface
type customizeRdsStore struct {
	redisClient *redis.Client
}

// customizeRdsStore implementing Set method of  Store interface
func (s *customizeRdsStore) Set(id string, value string) {
	err := s.redisClient.Set(id, value, time.Minute*10).Err()
	if err != nil {

	}
}

// customizeRdsStore implementing Get method of  Store interface
func (s *customizeRdsStore) Get(id string, clear bool) (value string) {
	val, err := s.redisClient.Get(id).Result()
	if err != nil {

		return
	}
	if clear {
		err := s.redisClient.Del(id).Err()
		if err != nil {

			return
		}
	}
	return val
}

func (c *Auth) Register(ctx *gin.Context) {
	appG := api.Gin{C: ctx}
	cnt, err := service.GetUserCnt()
	if err != nil {

		appG.Response(http.StatusOK, 400001004, nil)
		return
	}
	if cnt >= int64(conf.Cnf.UserCnt) {

		appG.Response(http.StatusOK, 407000015, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}
func (c *Auth) AuthRegister(ctx *gin.Context) {
	appG := api.Gin{C: ctx}
	requestJson, exists := ctx.Get("json")
	if !exists {

		appG.Response(http.StatusOK, 401000004, nil)
		return
	}
	ar, ok := requestJson.(common.AuthRegister)
	if !ok {

		appG.Response(http.StatusOK, 400001001, nil)
		return
	}
	cnt, err := service.GetUserCnt()
	if err != nil {

		appG.Response(http.StatusOK, 400001004, nil)
		return
	}
	if cnt >= int64(conf.Cnf.UserCnt) {

		appG.Response(http.StatusOK, 400001004, nil)
		return
	}
	service.UserStore(ar)
	appG.Response(http.StatusOK, 0, nil)
	return
}
func (c *Auth) Login(ctx *gin.Context) {
	appG := api.Gin{C: ctx}
	customStore := customizeRdsStore{conf.CacheClient}
	base64Captcha.SetCustomStore(&customStore)
	var configD = base64Captcha.ConfigDigit{
		Height:     80,
		Width:      240,
		MaxSkew:    0.7,
		DotCount:   80,
		CaptchaLen: 5,
	}
	idKeyD, capD := base64Captcha.GenerateCaptcha("", configD)
	base64stringD := base64Captcha.CaptchaWriteToBase64Encoding(capD)
	data := make(map[string]interface{})
	data["key"] = idKeyD
	data["png"] = base64stringD
	appG.Response(http.StatusOK, 0, data)
	return
}
func (c *Auth) AuthLogin(ctx *gin.Context) {
	appG := api.Gin{C: ctx}
	requestJson, exists := ctx.Get("json")
	if !exists {

		appG.Response(http.StatusOK, 401000004, nil)
		return
	}
	al, ok := requestJson.(common.AuthLogin)
	if !ok {

		appG.Response(http.StatusOK, 400001001, nil)
		return
	}
	verifyResult := base64Captcha.VerifyCaptcha(al.CaptchaKey, al.Captcha)
	if !verifyResult {

		appG.Response(http.StatusOK, 407000008, nil)
		return
	}

	user, err := service.GetUserByEmail(al.Email)
	if err != nil {

		appG.Response(http.StatusOK, 407000010, nil)
		return
	}
	if user.Id <= 0 {

		appG.Response(http.StatusOK, 407000010, nil)
		return
	}

	password := []byte(al.Password)
	hashedPassword := []byte(user.Password)
	err = bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {

		appG.Response(http.StatusOK, 407000010, nil)
		return
	}

	userIdStr := strconv.Itoa(user.Id)
	token, err := jwt.CreateToken(userIdStr)
	if err != nil {

		appG.Response(http.StatusOK, 407000011, nil)
		return
	}
	appG.Response(http.StatusOK, 0, token)
	return
}

func (c *Auth) Logout(ctx *gin.Context) {
	appG := api.Gin{C: ctx}
	token, exist := ctx.Get("token")
	if !exist || token == "" {

		appG.Response(http.StatusOK, 400001004, nil)
		return
	}
	_, err := jwt.UnsetToken(token.(string))
	if err != nil {

		appG.Response(http.StatusOK, 407000014, nil)
		return
	}
	appG.Response(http.StatusOK, 0, token)
	return
}

func (c *Auth) DelCache(ctx *gin.Context) {
	appG := api.Gin{C: ctx}
	service.DelAllCache()
	appG.Response(http.StatusOK, 0, nil)
	return
}
