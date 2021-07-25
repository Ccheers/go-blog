package console

import (
	"github.com/gin-gonic/gin"
	"go-blog/common"
	"go-blog/common/api"
	"go-blog/service"
	"net/http"
	"strconv"
)

type Category struct {
}

func NewCategory() Console {
	return &Category{}
}

func (cate *Category) Index(c *gin.Context) {
	appG := api.Gin{C: c}
	cates, err := service.CateListBySort()
	if err != nil {

		appG.Response(http.StatusOK, 402000001, nil)
		return
	}
	appG.Response(http.StatusOK, 0, cates)
	return
}

func (cate *Category) Create(c *gin.Context) {

}

func (cate *Category) Store(c *gin.Context) {
	appG := api.Gin{C: c}
	requestJson, exists := c.Get("json")
	if !exists {

		appG.Response(http.StatusOK, 400001003, nil)
		return
	}
	var cs common.CateStore
	cs, ok := requestJson.(common.CateStore)
	if !ok {

		appG.Response(http.StatusOK, 400001001, nil)
		return
	}

	_, err := service.CateStore(cs)
	if err != nil {

		appG.Response(http.StatusOK, 402000010, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}

func (cate *Category) Edit(c *gin.Context) {
	cateIdStr := c.Param("id")
	cateIdInt, err := strconv.Atoi(cateIdStr)
	appG := api.Gin{C: c}

	if err != nil {

		appG.Response(http.StatusOK, 400001002, nil)
		return
	}
	cateData, err := service.GetCateById(cateIdInt)
	if err != nil {

		appG.Response(http.StatusOK, 402000000, nil)
		return
	}
	appG.Response(http.StatusOK, 0, cateData)
	return
}

func (cate *Category) Update(c *gin.Context) {
	cateIdStr := c.Param("id")
	cateIdInt, err := strconv.Atoi(cateIdStr)
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
	var cs common.CateStore
	cs, ok := requestJson.(common.CateStore)
	if !ok {

		appG.Response(http.StatusOK, 400001001, nil)
		return
	}
	_, err = service.CateUpdate(cateIdInt, cs)
	if err != nil {

		appG.Response(http.StatusOK, 402000009, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}

func (cate *Category) Destroy(c *gin.Context) {
	cateIdStr := c.Param("id")
	cateIdInt, err := strconv.Atoi(cateIdStr)
	appG := api.Gin{C: c}

	if err != nil {

		appG.Response(http.StatusOK, 400001002, nil)
		return
	}

	_, err = service.GetCateById(cateIdInt)
	if err != nil {

		appG.Response(http.StatusOK, 402000000, nil)
		return
	}

	pd, err := service.GetCateByParentId(cateIdInt)
	if err != nil {

		appG.Response(http.StatusOK, 402000000, nil)
		return
	}
	if pd.Id > 0 {

		appG.Response(http.StatusOK, 402000011, nil)
		return
	}

	service.DelCateRel(cateIdInt)
	appG.Response(http.StatusOK, 0, nil)
	return
}
