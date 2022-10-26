package index

import (
	"bytes"
	"go-blog/common"
	"go-blog/conf"
	"go-blog/entity"
	"go-blog/service"
	"html/template"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

type Web struct {
	ApiController
}

func NewIndex() Home {
	return &Web{}
}

func (w *Web) Index(c *gin.Context) {

	queryPage := c.DefaultQuery("page", "1")
	queryLimit := c.DefaultQuery("limit", conf.Cnf.DefaultIndexLimit)

	h, err := service.CommonData()
	if err != nil {

		w.Response(c, http.StatusOK, 404, h)
		return
	}

	postData, err := service.IndexPost(queryPage, queryLimit, "default", "")
	if err != nil {

		w.Response(c, http.StatusOK, 404, h)
		return
	}

	h["post"] = postData.PostListArr
	h["paginate"] = postData.Paginate
	h["title"] = h["system"].(*entity.ZSystems).Title
	w.Response(c, http.StatusOK, 0, h)
	return
}

func (w *Web) IndexTag(c *gin.Context) {

	queryPage := c.DefaultQuery("page", "1")
	queryLimit := c.DefaultQuery("limit", conf.Cnf.DefaultIndexLimit)
	name := c.Param("name")
	h, err := service.CommonData()
	if err != nil {

		w.Response(c, http.StatusOK, 404, h)
		return
	}

	postData, err := service.IndexPost(queryPage, queryLimit, "tag", name)
	if err != nil {

		w.Response(c, http.StatusOK, 404, h)
		return
	}

	h["post"] = postData.PostListArr
	h["paginate"] = postData.Paginate
	h["tagName"] = name
	h["tem"] = "tagList"
	h["title"] = template.HTML(name + " --  tags &nbsp;&nbsp;-&nbsp;&nbsp;" + h["system"].(*entity.ZSystems).Title)

	c.HTML(http.StatusOK, "master.tmpl", h)
	return
}

func (w *Web) IndexCate(c *gin.Context) {

	queryPage := c.DefaultQuery("page", "1")
	queryLimit := c.DefaultQuery("limit", conf.Cnf.DefaultIndexLimit)
	name := c.Param("name")

	h, err := service.CommonData()
	if err != nil {

		w.Response(c, http.StatusOK, 404, h)
		return
	}

	postData, err := service.IndexPost(queryPage, queryLimit, "cate", name)
	if err != nil {

		w.Response(c, http.StatusOK, 404, h)
		return
	}

	h["post"] = postData.PostListArr
	h["paginate"] = postData.Paginate
	h["cateName"] = name
	h["tem"] = "cateList"
	h["title"] = template.HTML(name + " --  category &nbsp;&nbsp;-&nbsp;&nbsp;" + h["system"].(*entity.ZSystems).Title)

	w.Response(c, http.StatusOK, 0, h)
	return

}

func (w *Web) Detail(c *gin.Context) {

	postIdStr := c.Param("id")

	h, err := service.CommonData()
	if err != nil {

		w.Response(c, http.StatusOK, 404, h)
		return
	}

	postDetail, err := service.IndexPostDetail(postIdStr)
	if err != nil {

		w.Response(c, http.StatusOK, 404, h)
		return
	}

	go service.PostViewAdd(postIdStr)

	github := common.IndexGithubParam{
		GithubName:         conf.Cnf.GithubName,
		GithubRepo:         conf.Cnf.GithubRepo,
		GithubClientId:     conf.Cnf.GithubClientId,
		GithubClientSecret: conf.Cnf.GithubClientSecret,
		GithubLabels:       conf.Cnf.GithubLabels,
	}

	h["post"] = postDetail
	h["github"] = github
	h["tem"] = "detail"
	h["title"] = template.HTML(postDetail.Post.Title + " &nbsp;&nbsp;-&nbsp;&nbsp;" + h["system"].(*entity.ZSystems).Title)

	w.Response(c, http.StatusOK, 0, h)
	return
}

func (w *Web) Archives(c *gin.Context) {

	h, err := service.CommonData()
	if err != nil {

		w.Response(c, http.StatusOK, 404, h)
		return
	}

	res, err := service.PostArchives()
	if err != nil {

		w.Response(c, http.StatusOK, 404, h)
		return
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")

	var dateIndexs []int
	for k, _ := range res {
		tt, _ := time.ParseInLocation("2006-01-02 15:04:05", k+"-01 00:00:00", loc)
		dateIndexs = append(dateIndexs, int(tt.Unix()))
	}

	sort.Sort(sort.Reverse(sort.IntSlice(dateIndexs)))

	var newData []interface{}
	for _, j := range dateIndexs {
		dds := make(map[string]interface{})
		tm := time.Unix(int64(j), 0)
		dateIndex := tm.Format("2006-01")
		dds["dates"] = dateIndex
		dds["lists"] = res[dateIndex]
		newData = append(newData, dds)
	}

	h["tem"] = "archives"
	h["archives"] = newData
	h["title"] = template.HTML("归档 &nbsp;&nbsp;-&nbsp;&nbsp;" + h["system"].(*entity.ZSystems).Title)

	w.Response(c, http.StatusOK, 0, h)
	return
}

func (w *Web) NoFound(c *gin.Context) {

	w.Response(c, http.StatusOK, 404, gin.H{
		"themeJs":  "/static/home/assets/js",
		"themeCss": "/static/home/assets/css",
	})
	return
}

func (w *Web) SiteMap(c *gin.Context) {
	siteMap := service.GetSiteMap()
	buf := &bytes.Buffer{}
	buf.WriteString(siteMap.String())
	c.DataFromReader(http.StatusOK, int64(buf.Len()), "text/xml; charset=UTF-8", buf, make(map[string]string))
	return
}
