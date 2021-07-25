package service

import (
	"fmt"
	"go-blog/common"
	"go-blog/conf"
	"go-blog/entity"
	"golang.org/x/crypto/bcrypt"
)

func GetUserByEmail(email string) (user *entity.ZUsers, err error) {
	user = new(entity.ZUsers)
	_, err = conf.SqlServer.Where("email = ?", email).Get(user)
	return
}

func GetUserCnt() (cnt int64, err error) {
	user := new(entity.ZUsers)
	cnt, err = conf.SqlServer.Count(user)
	return
}

func UserStore(ar common.AuthRegister) (user *entity.ZUsers, err error) {
	password := []byte(ar.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {

		return
	}
	userInsert := entity.ZUsers{
		Name:     ar.UserName,
		Email:    ar.Email,
		Password: string(hashedPassword),
		Status:   1,
	}
	_, err = conf.SqlServer.Insert(&userInsert)
	if err != nil {

		return
	}
	fmt.Println(userInsert.Id)
	return
}

func DelAllCache() {
	conf.CacheClient.Del(
		conf.Cnf.TagListKey,
		conf.Cnf.CateListKey,
		conf.Cnf.ArchivesKey,
		conf.Cnf.LinkIndexKey,
		conf.Cnf.PostIndexKey,
		conf.Cnf.SystemIndexKey,
		conf.Cnf.TagPostIndexKey,
		conf.Cnf.CatePostIndexKey,
		conf.Cnf.PostDetailIndexKey,
	)
}
