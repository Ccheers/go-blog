package service

import (
	"go-blog/conf"
	"go-blog/entity"
)

func GetUserById(userId int) (*entity.ZUsers, error) {
	user := new(entity.ZUsers)
	_, err := conf.SqlServer.ID(userId).Cols("name", "email").Get(user)
	if err != nil {

		return user, err
	}
	return user, nil
}

func UserCnt() (cnt int64, err error) {
	user := new(entity.ZUsers)
	cnt, err = conf.SqlServer.Count(user)
	return
}
