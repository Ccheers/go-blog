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

type Link struct {
}

func NewLink() Console {
	return &Link{}
}

func (l *Link) Index(c *gin.Context) {
	appG := api.Gin{C: c}

	queryPage := c.DefaultQuery("page", "1")
	queryLimit := c.DefaultQuery("limit", conf.Cnf.DefaultLimit)
	queryPageInt, err := strconv.Atoi(queryPage)
	if err != nil {

		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	limit, offset := common.Offset(queryPage, queryLimit)

	links, cnt, err := service.LinkList(offset, limit)
	if err != nil {

		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	data := make(map[string]interface{})
	data["list"] = links
	data["page"] = common.MyPaginate(cnt, limit, queryPageInt)

	appG.Response(http.StatusOK, 0, data)
	return
}
func (l *Link) Create(c *gin.Context) {
}
func (l *Link) Store(c *gin.Context) {
	appG := api.Gin{C: c}
	requestJson, exists := c.Get("json")
	if !exists {

		appG.Response(http.StatusOK, 401000004, nil)
		return
	}
	ls, ok := requestJson.(common.LinkStore)
	if !ok {

		appG.Response(http.StatusOK, 400001001, nil)
		return
	}

	err := service.LinkSore(ls)
	if err != nil {

		appG.Response(http.StatusOK, 406000005, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}
func (l *Link) Edit(c *gin.Context) {
	linkIdStr := c.Param("id")
	linkIdInt, err := strconv.Atoi(linkIdStr)
	appG := api.Gin{C: c}

	if err != nil {

		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	link, err := service.LinkDetail(linkIdInt)
	if err != nil {

		appG.Response(http.StatusOK, 406000006, nil)
		return
	}
	appG.Response(http.StatusOK, 0, link)
	return
}
func (l *Link) Update(c *gin.Context) {
	linkIdStr := c.Param("id")
	linkIdInt, err := strconv.Atoi(linkIdStr)
	appG := api.Gin{C: c}

	if err != nil {

		appG.Response(http.StatusOK, 500000000, nil)
		return
	}

	requestJson, exists := c.Get("json")
	if !exists {

		appG.Response(http.StatusOK, 400001003, nil)
		return
	}
	ls, ok := requestJson.(common.LinkStore)
	if !ok {

		appG.Response(http.StatusOK, 400001001, nil)
		return
	}
	err = service.LinkUpdate(ls, linkIdInt)
	if err != nil {

		appG.Response(http.StatusOK, 406000007, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}

func (l *Link) Destroy(c *gin.Context) {
	linkIdStr := c.Param("id")
	linkIdInt, err := strconv.Atoi(linkIdStr)
	appG := api.Gin{C: c}

	if err != nil {

		appG.Response(http.StatusOK, 500000000, nil)
		return
	}

	err = service.LinkDestroy(linkIdInt)
	if err != nil {

		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}
