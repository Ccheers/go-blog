package console

import (
	"github.com/gin-gonic/gin"
	"go-blog/common"
	"go-blog/common/api"
	"go-blog/conf"
	"go-blog/service"
	"net/http"
	"strconv"
)

type Tag struct {
}

func NewTag() Console {
	return &Tag{}
}

func (t *Tag) Index(c *gin.Context) {
	appG := api.Gin{C: c}

	queryPage := c.DefaultQuery("page", "1")
	queryLimit := c.DefaultQuery("limit", conf.Cnf.DefaultLimit)

	limit, offset := common.Offset(queryPage, queryLimit)
	count, tags, err := service.TagsIndex(limit, offset)
	if err != nil {

		appG.Response(http.StatusOK, 402000001, nil)
		return
	}
	queryPageInt, err := strconv.Atoi(queryPage)
	if err != nil {

		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	data := make(map[string]interface{})
	data["list"] = tags
	data["page"] = common.MyPaginate(count, limit, queryPageInt)

	appG.Response(http.StatusOK, 0, data)
	return
}

func (t *Tag) Create(c *gin.Context) {

}

func (t *Tag) Store(c *gin.Context) {
	appG := api.Gin{C: c}
	requestJson, exists := c.Get("json")
	if !exists {

		appG.Response(http.StatusOK, 400001003, nil)
		return
	}
	var ts common.TagStore
	ts, ok := requestJson.(common.TagStore)
	if !ok {

		appG.Response(http.StatusOK, 400001001, nil)
		return
	}
	err := service.TagStore(ts)
	if err != nil {

		appG.Response(http.StatusOK, 403000006, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}

func (t *Tag) Edit(c *gin.Context) {
	tagIdStr := c.Param("id")
	tagIdInt, err := strconv.Atoi(tagIdStr)
	appG := api.Gin{C: c}

	if err != nil {

		appG.Response(http.StatusOK, 400001002, nil)
		return
	}
	tagData, err := service.GetTagById(tagIdInt)
	if err != nil {

		appG.Response(http.StatusOK, 403000008, nil)
		return
	}
	appG.Response(http.StatusOK, 0, tagData)
	return
}

func (t *Tag) Update(c *gin.Context) {
	tagIdStr := c.Param("id")
	tagIdInt, err := strconv.Atoi(tagIdStr)
	appG := api.Gin{C: c}

	if err != nil {

		appG.Response(http.StatusOK, 400001002, nil)
		return
	}
	requestJson, exists := c.Get("json")
	if !exists {

		appG.Response(http.StatusOK, 400001003, nil)
		return
	}
	ts, ok := requestJson.(common.TagStore)
	if !ok {

		appG.Response(http.StatusOK, 400001001, nil)
		return
	}
	err = service.TagUpdate(tagIdInt, ts)
	if err != nil {

		appG.Response(http.StatusOK, 403000007, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}

func (t *Tag) Destroy(c *gin.Context) {
	tagIdStr := c.Param("id")
	tagIdInt, err := strconv.Atoi(tagIdStr)
	appG := api.Gin{C: c}

	if err != nil {

		appG.Response(http.StatusOK, 400001002, nil)
		return
	}

	_, err = service.GetTagById(tagIdInt)
	if err != nil {

		appG.Response(http.StatusOK, 403000008, nil)
		return
	}
	service.DelTagRel(tagIdInt)
	appG.Response(http.StatusOK, 0, nil)
	return
}
