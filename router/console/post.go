package console

import (
	"github.com/gin-gonic/gin"
	"go-blog/common"
	"go-blog/common/api"
	"go-blog/conf"
	"go-blog/service"
	"net/http"
	"path/filepath"
	"strconv"
)

type Post struct {
}

func NewPost() Console {
	return &Post{}
}

func NewPostImg() Img {
	return &Post{}
}

func NewTrash() Trash {
	return &Post{}
}

func (p *Post) Index(c *gin.Context) {
	appG := api.Gin{C: c}

	queryPage := c.DefaultQuery("page", "1")
	queryLimit := c.DefaultQuery("limit", conf.Cnf.DefaultLimit)

	limit, offset := common.Offset(queryPage, queryLimit)
	postList, err := service.ConsolePostIndex(limit, offset, false)
	if err != nil {
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	queryPageInt, err := strconv.Atoi(queryPage)
	if err != nil {
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	postCount, err := service.ConsolePostCount(limit, offset, false)

	data := make(map[string]interface{})
	data["list"] = postList
	data["page"] = common.MyPaginate(postCount, limit, queryPageInt)

	appG.Response(http.StatusOK, 0, data)
	return
}

func (p *Post) Create(c *gin.Context) {
	cates, err := service.CateListBySort()
	appG := api.Gin{C: c}
	if err != nil {
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	tags, err := service.AllTags()
	if err != nil {
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	data := make(map[string]interface{})
	data["cates"] = cates
	data["tags"] = tags
	data["imgUploadUrl"] = conf.Cnf.ImgUploadUrl
	appG.Response(http.StatusOK, 0, data)
	return
}

func (p *Post) Store(c *gin.Context) {
	appG := api.Gin{C: c}
	requestJson, exists := c.Get("json")
	if !exists {
		appG.Response(http.StatusOK, 401000004, nil)
		return
	}
	var ps common.PostStore
	ps, ok := requestJson.(common.PostStore)
	if !ok {
		appG.Response(http.StatusOK, 400001001, nil)
		return
	}

	userId, exist := c.Get("userId")
	if !exist || userId.(int) == 0 {
		appG.Response(http.StatusOK, 400001004, nil)
		return
	}

	service.PostStore(ps, userId.(int))
	appG.Response(http.StatusOK, 0, nil)
	return
}

func (p *Post) Edit(c *gin.Context) {
	postIdStr := c.Param("id")
	postIdInt, err := strconv.Atoi(postIdStr)
	appG := api.Gin{C: c}

	if err != nil {
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	post, err := service.PostDetail(postIdInt)
	if err != nil {
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	postTags, err := service.PostIdTag(postIdInt)
	if err != nil {
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	postCate, err := service.PostCate(postIdInt)
	if err != nil {
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	data := make(map[string]interface{})
	posts := make(map[string]interface{})
	posts["post"] = post
	posts["postCate"] = postCate
	posts["postTag"] = postTags
	data["post"] = posts
	cates, err := service.CateListBySort()
	if err != nil {
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	tags, err := service.AllTags()
	if err != nil {
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	data["cates"] = cates
	data["tags"] = tags
	data["imgUploadUrl"] = conf.Cnf.ImgUploadUrl
	appG.Response(http.StatusOK, 0, data)
	return
}

func (p *Post) Update(c *gin.Context) {
	postIdStr := c.Param("id")
	postIdInt, err := strconv.Atoi(postIdStr)
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
	var ps common.PostStore
	ps, ok := requestJson.(common.PostStore)
	if !ok {
		appG.Response(http.StatusOK, 400001001, nil)
		return
	}
	service.PostUpdate(postIdInt, ps)
	appG.Response(http.StatusOK, 0, nil)
	return
}

func (p *Post) Destroy(c *gin.Context) {
	postIdStr := c.Param("id")
	postIdInt, err := strconv.Atoi(postIdStr)
	appG := api.Gin{C: c}

	if err != nil {
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}

	_, err = service.PostDestroy(postIdInt)
	if err != nil {
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}

func (p *Post) TrashIndex(c *gin.Context) {
	appG := api.Gin{C: c}

	queryPage := c.DefaultQuery("page", "1")
	queryLimit := c.DefaultQuery("limit", conf.Cnf.DefaultLimit)

	limit, offset := common.Offset(queryPage, queryLimit)
	postList, err := service.ConsolePostIndex(limit, offset, true)
	if err != nil {
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	queryPageInt, err := strconv.Atoi(queryPage)
	if err != nil {
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	postCount, err := service.ConsolePostCount(limit, offset, true)

	data := make(map[string]interface{})
	data["list"] = postList
	data["page"] = common.MyPaginate(postCount, limit, queryPageInt)

	appG.Response(http.StatusOK, 0, data)
	return
}

func (p *Post) UnTrash(c *gin.Context) {
	postIdStr := c.Param("id")
	postIdInt, err := strconv.Atoi(postIdStr)
	appG := api.Gin{C: c}

	if err != nil {
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	_, err = service.PostUnTrash(postIdInt)
	if err != nil {
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}

func (p *Post) ImgUpload(c *gin.Context) {
	appG := api.Gin{C: c}

	file, err := c.FormFile("file")
	if err != nil {
		appG.Response(http.StatusOK, 401000004, nil)
		return
	}

	filename := filepath.Base(file.Filename)
	dst := conf.Dirname + "/" + conf.Cnf.ImgUploadDst + filename
	if err := c.SaveUploadedFile(file, dst); err != nil {
		appG.Response(http.StatusOK, 401000005, nil)
		return
	}

	// Default upload both
	data := make(map[string]interface{})
	if conf.Cnf.ImgUploadBoth {
		go service.Qiniu(dst, filename)
		data["path"] = conf.Cnf.AppImgUrl + filename
	} else {
		if conf.Cnf.QiNiuUploadImg {
			go service.Qiniu(dst, filename)
			data["path"] = conf.Cnf.QiNiuHostName + filename
		} else {
			data["path"] = conf.Cnf.AppImgUrl + filename
		}
	}

	appG.Response(http.StatusOK, 0, data)
	return
}
