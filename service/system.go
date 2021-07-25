package service

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"go-blog/common"
	"go-blog/conf"
	"go-blog/entity"
	"time"
)

func GetSystemList() (system *entity.ZSystems, err error) {
	system = new(entity.ZSystems)
	_, err = conf.SqlServer.Get(system)
	if err != nil {

		return
	}
	if system.Id <= 0 {
		systemInsert := entity.ZSystems{
			Theme:        conf.Cnf.Theme,
			Title:        conf.Cnf.Title,
			Keywords:     conf.Cnf.Keywords,
			Description:  conf.Cnf.Description,
			RecordNumber: conf.Cnf.RecordNumber,
		}
		_, err = conf.SqlServer.Insert(systemInsert)
		if err != nil {

			return
		}
		_, err = conf.SqlServer.Get(system)
		if err != nil {

			return
		}
	}
	return
}

func SystemUpdate(sId int, ss common.ConsoleSystem) error {
	systemUpdate := entity.ZSystems{
		Title:        ss.Title,
		Keywords:     ss.Keywords,
		Description:  ss.Description,
		RecordNumber: ss.RecordNumber,
		Theme:        ss.Theme,
	}
	_, err := conf.SqlServer.ID(sId).Update(&systemUpdate)
	return err
}

func IndexSystem() (system *entity.ZSystems, err error) {
	cacheKey := conf.Cnf.SystemIndexKey
	cacheRes, err := conf.CacheClient.Get(cacheKey).Result()
	if err == redis.Nil {
		system, err := doCacheIndexSystem(cacheKey)
		if err != nil {

			return system, err
		}
		return system, nil
	} else if err != nil {

		return system, err
	}

	err = json.Unmarshal([]byte(cacheRes), &system)
	if err != nil {

		system, err = doCacheIndexSystem(cacheKey)
		if err != nil {

			return nil, err
		}
		return system, nil
	}
	return system, nil
}

func doCacheIndexSystem(cacheKey string) (system *entity.ZSystems, err error) {
	system = new(entity.ZSystems)
	_, err = conf.SqlServer.Get(system)
	if err != nil {

		return system, err
	}
	jsonRes, err := json.Marshal(&system)
	if err != nil {

		return system, err
	}
	err = conf.CacheClient.Set(cacheKey, jsonRes, time.Duration(conf.Cnf.DataCacheTimeDuration)*time.Hour).Err()
	if err != nil {

		return system, err
	}
	return system, nil
}
