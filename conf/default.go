package conf

import (
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/speps/go-hashids"
	"go-blog/common/alarm"
	"go-blog/common/hashid"
	"go-blog/common/jwt"
	"go-blog/common/mail"
	"go-blog/common/qqcaptcha"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
	"xorm.io/xorm"
)

var (
	SqlServer   *xorm.Engine
	ZHashId     *hashids.HashID
	CacheClient *redis.Client
	MailClient  *mail.EmailParam
	Cnf         *Conf
	Env         string
	Dirname     string
)

func init() {
	flag.StringVar(&Env, "env", "dev", "environment")
	flag.StringVar(&Dirname, "dir", "./config", "environment")
	flag.Parse()
}

func DefaultInit() {
	CnfInit()
	DbInit()
	AlarmInit()
	MailInit()
	ZHashIdInit()
	RedisInit()
	JwtInit()
	QCaptchaInit()
	// the customer error code init
	SetMsg(Msg)
	//BackUpInit()
}

func DbInit() {
	sqlServer, err := xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", Cnf.DbUser, Cnf.DbPassword, Cnf.DbHost, Cnf.DbPort, Cnf.DbDataBase))
	SqlServer = sqlServer
	if err != nil {
		panic(err.Error())
	}
	err = SqlServer.Ping()
	if err != nil {
		panic(err)
	}
	return
}

func AlarmInit() {
	a := new(alarm.Param)
	alarmT := a.SetType(alarm.Type(Cnf.AlarmType))
	mailTo := a.SetMailTo("xzghua@gmail.com")
	err := a.AlarmInit(alarmT, mailTo)
	if err != nil {
	}
	return
}

func MailInit() {
	m := new(mail.EmailParam)
	mailUser := m.SetMailUser(mail.EmailType(Cnf.MailUser))
	mailPwd := m.SetMailPwd(mail.EmailType(Cnf.MailPwd))
	mailHost := m.SetMailHost(mail.EmailType(Cnf.MailHost))
	mails, err := m.MailInit(mailPwd, mailHost, mailUser)
	if err != nil {
		panic(err.Error())
	}
	MailClient = mails
	return
}

func ZHashIdInit() {
	hd := new(hashid.Params)
	salt := hd.SetHashIdSalt(Cnf.HashIdSalt)
	hdLength := hd.SetHashIdLength(Cnf.HashIdLength)
	zHashId, err := hd.HashIdInit(hdLength, salt)
	if err != nil {
	}
	ZHashId = zHashId
	return
}

func RedisInit() {

	options := &redis.Options{
		Network:  "tcp",
		Addr:     Cnf.RedisAddr,
		Password: Cnf.RedisPwd,
		DB:       Cnf.RedisDb,
	}
	client := redis.NewClient(options)
	err := client.Ping().Err()
	if err != nil {
		panic(err.Error())
	}
	CacheClient = client
	return
}

func JwtInit() {
	jt := new(jwt.Param)
	ad := jt.SetDefaultAudience(Cnf.JwtAudience)
	jti := jt.SetDefaultJti(Cnf.JwtJti)
	iss := jt.SetDefaultIss(Cnf.JwtIss)
	sk := jt.SetDefaultSecretKey(Cnf.JwtSecretKey)
	rc := jt.SetRedisCache(CacheClient)
	tl := jt.SetTokenLife(time.Hour * time.Duration(Cnf.JwtTokenLife))
	_ = jt.JwtInit(ad, jti, iss, sk, rc, tl)
	return
}

func QCaptchaInit() {
	qc := new(qqcaptcha.QQCaptcha)
	aid := qc.SetAid(Cnf.QCaptchaAid)
	sk := qc.SetSecretKey(Cnf.QCaptchaSecretKey)
	_ = qc.QQCaptchaInit(aid, sk)
	return
}

func CnfInit() {
	cf := &Conf{
		AppUrl:                "http://localhost:8081",
		AppImgUrl:             "http://localhost:8081/static/uploads/images/",
		DefaultLimit:          "20",
		DefaultIndexLimit:     "3",
		DbUser:                "root",
		DbPassword:            "",
		DbPort:                "3306",
		DbDataBase:            "go-blog",
		DbHost:                "127.0.0.1",
		AlarmType:             "mail,wechat",
		MailUser:              "test@test.com",
		MailPwd:               "",
		MailHost:              "smtp.mxhichina.com:25",
		HashIdSalt:            "i must add a salt what is only for me",
		HashIdLength:          8,
		JwtIss:                "go-blog",
		JwtAudience:           "blog",
		JwtJti:                "go-blog",
		JwtSecretKey:          "go-blog",
		JwtTokenKey:           "login:token:",
		JwtTokenLife:          3,
		RedisAddr:             "localhost:6379",
		RedisPwd:              "",
		RedisDb:               0,
		QCaptchaAid:           "",
		QCaptchaSecretKey:     "**",
		BackUpFilePath:        "./backup/",
		BackUpDuration:        "* * */1 * *",
		BackUpSentTo:          "xzghua@gmail.com",
		DataCacheTimeDuration: 720,
		ImgUploadUrl:          "http://localhost:8081/console/post/imgUpload",
		ImgUploadDst:          "./static/uploads/images/",
		ImgUploadBoth:         true, // img will upload to qiniu and your server local
		QiNiuUploadImg:        true,
		QiNiuHostName:         "",
		QiNiuAccessKey:        "",
		QiNiuSecretKey:        "",
		QiNiuBucket:           "go-blog",
		QiNiuZone:             "HUABEI",
		CateListKey:           "all:cate:sort",
		TagListKey:            "all:tag",
		Theme:                 0,
		Title:                 "默认Title",
		Keywords:              "默认关键词,叶落山城秋",
		Description:           "个人网站,https://github.com/izghua/go-blog",
		RecordNumber:          "000-0000",
		UserCnt:               2,
		PostIndexKey:          "index:all:post:list",
		TagPostIndexKey:       "index:all:tag:post:list",
		CatePostIndexKey:      "index:all:cate:post:list",
		LinkIndexKey:          "index:all:link:list",
		SystemIndexKey:        "index:all:system:list",
		PostDetailIndexKey:    "index:post:detail",
		ArchivesKey:           "index:archives:list",
		GithubName:            "",
		GithubRepo:            "",
		GithubClientId:        "",
		GithubClientSecret:    "",
		GithubLabels:          "Gitalk",
		ThemeJs:               "/static/home/assets/js",
		ThemeCss:              "/static/home/assets/css",
		ThemeImg:              "/static/home/assets/img",
		ThemeFancyboxCss:      "/static/home/assets/fancybox",
		ThemeFancyboxJs:       "/static/home/assets/fancybox",
		ThemeHLightCss:        "/static/home/assets/highlightjs",
		ThemeHLightJs:         "/static/home/assets/highlightjs",
		ThemeShareCss:         "/static/home/assets/css",
		ThemeShareJs:          "/static/home/assets/js",
		ThemeArchivesJs:       "/static/home/assets/js",
		ThemeArchivesCss:      "/static/home/assets/css",
		ThemeNiceImg:          "/static/home/assets/img",
		ThemeAllCss:           "/static/home/assets/css",
		ThemeIndexImg:         "/static/home/assets/img",
		ThemeCateImg:          "/static/home/assets/img",
		ThemeTagImg:           "/static/home/assets/img",
	}
	var fileName string
	log.Println(fmt.Sprintf("Env : %s", Env))
	switch Env {
	case "dev":
		fileName = "env.dev.yaml"
	case "prod":
		fileName = "env.prod.yaml"
	default:
		fileName = "default"
	}
	if fileName == "default" {
		Cnf = cf
		return
	}
	log.Println(fmt.Sprintf("local directory : %s", Dirname))
	log.Printf(fmt.Sprintf("%s/%s", Dirname, fileName))
	//读取yaml配置文件
	yamlFile, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", Dirname, fileName))
	if err != nil {
	}
	err = yaml.Unmarshal(yamlFile, &cf)
	if err != nil {
	}
	Cnf = cf
	return
}
