package console

import (
	"github.com/gin-gonic/gin"
	"go-blog/common"
	"go-blog/common/api"
	"go-blog/service"
	"net/http"
	"strconv"
)

type Home struct {
}

func NewHome() System {
	return &Home{}
}

func (s *Home) Index(c *gin.Context) {
	appG := api.Gin{C: c}
	themes := make(map[int]interface{})
	themes[1] = 1
	system, err := service.GetSystemList()
	if err != nil {

		return
	}
	data := make(map[string]interface{})
	data["themes"] = themes
	data["system"] = system

	appG.Response(http.StatusOK, 0, data)
	return
}

func (s *Home) Update(c *gin.Context) {
	systemIdStr := c.Param("id")
	systemIdInt, err := strconv.Atoi(systemIdStr)
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
	//var ss common.ConsoleSystem
	ss, ok := requestJson.(common.ConsoleSystem)
	if !ok {

		appG.Response(http.StatusOK, 400001001, nil)
		return
	}
	err = service.SystemUpdate(systemIdInt, ss)
	if err != nil {

		appG.Response(http.StatusOK, 405000000, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}
