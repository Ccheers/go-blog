package validate

import (
	"github.com/gin-gonic/gin"
	"go-blog/common"
	"go-blog/common/api"
	"net/http"
)

type SystemUpdateV struct {
}

func (sv *SystemUpdateV) MyValidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := api.Gin{C: c}
		var json common.ConsoleSystem
		//接收各种参数
		if err := c.ShouldBindJSON(&json); err != nil {
			appG.Response(http.StatusOK, 400001000, nil)
			return
		}

		reqValidate := &SystemUpdate{
			Title:        json.Title,
			Keywords:     json.Keywords,
			Description:  json.Description,
			RecordNumber: json.RecordNumber,
			Theme:        json.Theme,
		}
		if b := appG.Validate(reqValidate); !b {
			return
		}
		c.Set("json", json)
		c.Next()
	}
}

type SystemUpdate struct {
	Title        string `valid:"required,maxstringlength(100)"`
	Keywords     string `valid:"required,maxstringlength(100)"`
	Description  string `valid:"required,maxstringlength(250)"`
	RecordNumber string `valid:"required,maxstringlength(50)"`
	Theme        int    `valid:"required"`
}

func (c *SystemUpdate) Message() map[string]int {
	return map[string]int{
		"Title.Required":        405000001,
		"Title.MaxSize":         405000002,
		"Keywords.Required":     405000003,
		"Keywords.MaxSize":      405000004,
		"Description.Required":  405000005,
		"Description.MaxSize":   405000006,
		"RecordNumber.Required": 405000007,
		"RecordNumber.MaxSize":  405000008,
		"Theme.Required":        405000009,
	}
}
