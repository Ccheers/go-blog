package index

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ApiController struct {
}

func (a *ApiController) Response(c *gin.Context, httpCode, errCode int, data gin.H) {
	if data == nil {
		panic("常规信息应该设置")
	}
	//msg := conf.GetMsg(errCode)
	beginTime, _ := strconv.ParseInt(c.Writer.Header().Get("X-Begin-Time"), 10, 64)

	duration := time.Now().Sub(time.Unix(0, beginTime))
	milliseconds := float64(duration) / float64(time.Millisecond)
	rounded := float64(int(milliseconds*100+.5)) / 100
	roundedStr := fmt.Sprintf("%.3fms", rounded)
	c.Writer.Header().Set("X-Response-time", roundedStr)
	//requestUrl := c.Request.URL.String()
	//requestMethod := c.Request.Method

	if errCode == 500 {
		c.HTML(http.StatusOK, "5xx.tmpl", data)
	} else if errCode == 404 {
		c.HTML(http.StatusOK, "4xx.tmpl", data)
	} else if errCode == 0 {
		c.HTML(http.StatusOK, "master.tmpl", data)
	} else {
		c.HTML(http.StatusOK, "5xx.tmpl", nil)
	}

	c.Abort()
	return
}
