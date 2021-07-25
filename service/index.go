package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go-blog/common"
	"go-blog/conf"
	"go-blog/entity"
	"html/template"
	"strconv"
)

type IndexType string

const (
	IndexTypeOne   IndexType = "tag"
	IndexTypeTwo   IndexType = "cate"
	IndexTypeThree IndexType = "default"
)

func CommonData() (h gin.H, err error) {
	h = gin.H{
		"themeJs":          conf.Cnf.ThemeJs,
		"themeCss":         conf.Cnf.ThemeCss,
		"themeImg":         conf.Cnf.ThemeImg,
		"themeFancyboxCss": conf.Cnf.ThemeFancyboxCss,
		"themeFancyboxJs":  conf.Cnf.ThemeFancyboxJs,
		"themeHLightCss":   conf.Cnf.ThemeHLightCss,
		"themeHLightJs":    conf.Cnf.ThemeHLightJs,
		"themeShareCss":    conf.Cnf.ThemeShareCss,
		"themeShareJs":     conf.Cnf.ThemeShareJs,
		"themeArchivesJs":  conf.Cnf.ThemeArchivesJs,
		"themeArchivesCss": conf.Cnf.ThemeArchivesCss,
		"themeNiceImg":     conf.Cnf.ThemeNiceImg,
		"themeAllCss":      conf.Cnf.ThemeAllCss,
		"themeIndexImg":    conf.Cnf.ThemeIndexImg,
		"themeCateImg":     conf.Cnf.ThemeCateImg,
		"themeTagImg":      conf.Cnf.ThemeTagImg,
		"title":            "",

		"tem": "defaultList",
	}
	h["script"] = template.HTML(conf.Cnf.OtherScript)
	cates, err := CateListBySort()
	if err != nil {

		return
	}
	var catess []common.IndexCategory
	for _, v := range cates {
		c := common.IndexCategory{
			Cates: v.Cates,
			Html:  template.HTML(v.Html),
		}
		catess = append(catess, c)
	}

	tags, err := AllTags()
	if err != nil {

		return
	}

	links, err := AllLink()
	if err != nil {

		return
	}

	system, err := IndexSystem()
	if err != nil {

		return
	}
	h["cates"] = catess
	h["system"] = system
	h["links"] = links
	h["tags"] = tags
	return
}

func IndexPost(page string, limit string, indexType IndexType, name string) (indexPostIndex common.IndexPostList, err error) {
	var postKey string
	switch indexType {
	case IndexTypeOne:
		postKey = conf.Cnf.TagPostIndexKey
	case IndexTypeTwo:
		postKey = conf.Cnf.CatePostIndexKey
	case IndexTypeThree:
		postKey = conf.Cnf.PostIndexKey
		name = "default"
	default:
		postKey = conf.Cnf.PostIndexKey
	}

	field := ":name:" + name + ":page:" + page + ":limit:" + limit
	cacheRes, err := conf.CacheClient.HGet(postKey, field).Result()
	if err == redis.Nil {
		// cache key does not exist
		// set data to the cache what use the cache key
		indexPostIndex, err := doCacheIndexPostList(postKey, field, indexType, name, page, limit)
		if err != nil {

			return indexPostIndex, err
		}
		return indexPostIndex, nil
	} else if err != nil {

		return indexPostIndex, err
	}

	if cacheRes == "" {
		// Data is  null
		// set data to the cache what use the cache key
		indexPostIndex, err := doCacheIndexPostList(postKey, field, indexType, name, page, limit)
		if err != nil {

			return indexPostIndex, err
		}
		return indexPostIndex, nil
	}

	err = json.Unmarshal([]byte(cacheRes), &indexPostIndex)
	if err != nil {

		indexPostIndex, err := doCacheIndexPostList(postKey, field, indexType, name, page, limit)
		if err != nil {

			return indexPostIndex, err
		}
		return indexPostIndex, nil
	}
	return
}

func doCacheIndexPostList(cacheKey string, field string, indexType IndexType, name string, queryPage string, queryLimit string) (res common.IndexPostList, err error) {
	limit, offset := common.Offset(queryPage, queryLimit)
	queryPageInt, err := strconv.Atoi(queryPage)
	if err != nil {

		return
	}
	var postList []*common.ConsolePostList
	var postCount int64
	switch indexType {
	case IndexTypeOne:
		tag := new(entity.ZTags)
		_, err = conf.SqlServer.Where("Name = ?", name).Get(tag)
		if err != nil {

			return
		}
		postList, err = PostTagList(tag.Id, limit, offset)
		if err != nil {

			return
		}
		postCount, err = PostTagListCount(tag.Id, limit, offset)
		if err != nil {

			return
		}
	case IndexTypeTwo:
		cate := new(entity.ZCategories)
		_, err = conf.SqlServer.Where("Name = ?", name).Get(cate)
		if err != nil {

			return
		}
		postList, err = PostCateList(cate.Id, limit, offset)
		if err != nil {

			return
		}
		postCount, err = PostCateListCount(cate.Id, limit, offset)
		if err != nil {

			return
		}
	case IndexTypeThree:
		postList, err = ConsolePostIndex(limit, offset, false)
		if err != nil {

			return
		}
		postCount, err = ConsolePostCount(limit, offset, false)
		if err != nil {

			return
		}
	default:
		postList, err = ConsolePostIndex(limit, offset, false)
		if err != nil {

			return
		}

		postCount, err = ConsolePostCount(limit, offset, false)
		if err != nil {

			return
		}
	}

	paginate := common.MyPaginate(postCount, limit, queryPageInt)

	res = common.IndexPostList{
		PostListArr: postList,
		Paginate:    paginate,
	}

	jsonRes, err := json.Marshal(&res)
	if err != nil {

		return
	}
	err = conf.CacheClient.HSet(cacheKey, field, jsonRes).Err()
	if err != nil {

		return
	}
	return
}

func IndexPostDetail(postIdStr string) (postDetail common.IndexPostDetail, err error) {
	cacheKey := conf.Cnf.PostDetailIndexKey
	field := ":post:id:" + postIdStr

	postIdInt, err := strconv.Atoi(postIdStr)
	if err != nil {

		return
	}
	cacheRes, err := conf.CacheClient.HGet(cacheKey, field).Result()
	if err == redis.Nil {
		// cache key does not exist
		// set data to the cache what use the cache key
		postDetail, err := doCacheIndexPostDetail(cacheKey, field, postIdInt)
		if err != nil {

			return postDetail, err
		}
		return postDetail, nil
	} else if err != nil {

		return postDetail, err
	}

	if cacheRes == "" {
		// Data is  null
		// set data to the cache what use the cache key
		postDetail, err = doCacheIndexPostDetail(cacheKey, field, postIdInt)
		if err != nil {

			return postDetail, err
		}
		return postDetail, nil
	}

	err = json.Unmarshal([]byte(cacheRes), &postDetail)
	if err != nil {

		postDetail, err = doCacheIndexPostDetail(cacheKey, field, postIdInt)
		if err != nil {

			return postDetail, err
		}
		return postDetail, nil
	}
	return

}

func doCacheIndexPostDetail(cacheKey string, field string, postIdInt int) (postDetail common.IndexPostDetail, err error) {
	postDetail, err = IndexPostDetailDao(postIdInt)
	if err != nil {

		return
	}
	jsonRes, err := json.Marshal(&postDetail)
	if err != nil {

		return
	}
	err = conf.CacheClient.HSet(cacheKey, field, jsonRes).Err()
	if err != nil {

		return
	}
	return
}

func PostViewAdd(postIdStr string) {
	postIdInt, err := strconv.Atoi(postIdStr)
	if err != nil {

		return
	}
	_, err = conf.SqlServer.Where("post_id = ?", postIdInt).Incr("num").Update(entity.ZPostViews{})
	if err != nil {

		return
	}
	return
}

func PostArchives() (archivesList map[string][]*entity.ZPosts, err error) {
	cacheKey := conf.Cnf.ArchivesKey
	field := ":all:"

	cacheRes, err := conf.CacheClient.HGet(cacheKey, field).Result()
	if err == redis.Nil {
		// cache key does not exist
		// set data to the cache what use the cache key
		archivesList, err := doCacheArchives(cacheKey, field)
		if err != nil {

			return archivesList, err
		}
		return archivesList, nil
	} else if err != nil {

		return archivesList, err
	}

	if cacheRes == "" {
		// Data is  null
		// set data to the cache what use the cache key
		archivesList, err := doCacheArchives(cacheKey, field)
		if err != nil {

			return archivesList, err
		}
		return archivesList, nil
	}

	archivesList = make(map[string][]*entity.ZPosts)
	err = json.Unmarshal([]byte(cacheRes), &archivesList)
	if err != nil {

		archivesList, err := doCacheArchives(cacheKey, field)
		if err != nil {

			return archivesList, err
		}
		return archivesList, nil
	}
	return
}

func doCacheArchives(cacheKey string, field string) (archivesList map[string][]*entity.ZPosts, err error) {
	posts := make([]*entity.ZPosts, 0)
	err = conf.SqlServer.Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").Desc("created_at").Find(&posts)
	if err != nil {

		return
	}
	archivesList = make(map[string][]*entity.ZPosts)
	for _, v := range posts {
		date := v.CreatedAt.Format("2006-01")
		archivesList[date] = append(archivesList[date], v)
	}

	jsonRes, err := json.Marshal(&archivesList)
	if err != nil {

		return
	}
	err = conf.CacheClient.HSet(cacheKey, field, jsonRes).Err()
	if err != nil {

		return
	}
	return
}
