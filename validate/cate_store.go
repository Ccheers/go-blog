package validate

import (
	"github.com/gin-gonic/gin"
	"go-blog/common"
	"go-blog/common/api"
	"net/http"
)

type CateStoreV struct {
}

func (cv *CateStoreV) MyValidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := api.Gin{C: c}
		var json common.CateStore
		//接收各种参数
		if err := c.ShouldBindJSON(&json); err != nil {
			appG.Response(http.StatusOK, 400001000, nil)
			return
		}

		reqValidate := &CateStore{
			Name:        json.Name,
			DisplayName: json.DisplayName,
			ParentId:    json.ParentId,
			SeoDesc:     json.SeoDesc,
		}
		if b := appG.Validate(reqValidate); !b {
			return
		}
		c.Set("json", json)
		c.Next()
	}
}

type CateStore struct {
	Name        string `valid:"required,maxstringlength(100)"`
	DisplayName string `valid:"required,maxstringlength(100)"`
	ParentId    int    `valid:"int"`
	SeoDesc     string `valid:"required,maxstringlength(250)"`
}

func (c *CateStore) Message() map[string]int {
	return map[string]int{
		"Name.Required":        402000002,
		"Name.MaxSize":         402000006,
		"DisplayName.Required": 402000003,
		"DisplayName.MaxSize":  402000007,
		"ParentId.Min":         402000004,
		"SeoDesc.Required":     402000005,
		"SeoDesc.MaxSize":      402000008,
	}
}
