package validate

import (
	"github.com/gin-gonic/gin"
	"go-blog/common"
	"go-blog/common/api"
	"net/http"
)

type PostStoreV struct {
}

func (pv *PostStoreV) MyValidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := api.Gin{C: c}
		var json common.PostStore
		//接收各种参数
		if err := c.ShouldBindJSON(&json); err != nil {
			appG.Response(http.StatusOK, 400001000, nil)
			return
		}

		reqValidate := &PostStore{
			Title:   json.Title,
			Tags:    json.Tags,
			Summary: json.Summary,
		}
		if b := appG.Validate(reqValidate); !b {
			return
		}
		c.Set("json", json)
		c.Next()
	}
}

type PostStore struct {
	Title string `valid:"required"`
	Tags  []int
	//Category int `valid:Required`
	Summary string `valid:"required"`
}

func (c *PostStore) Message() map[string]int {
	return map[string]int{
		"Title.Required":   401000000,
		"Summary.Required": 401000003,
	}
}
