package middleware

import (
	"github.com/gin-gonic/gin"
	"go-blog/common"
	"go-blog/common/api"
	"go-blog/common/jwt"
	"net/http"
	"strconv"
)

func Permission(routerAsName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiG := api.Gin{C: c}
		res := common.CheckPermissions(routerAsName, c.Request.Method)
		if !res {
			return
		}

		token := c.GetHeader("x-auth-token")
		if routerAsName == "console.post.imgUpload" {
			token = c.PostForm("upload-token")
		}

		if token == "" {
			apiG.Response(http.StatusOK, 400001005, nil)
			return
		}

		userId, err := jwt.ParseToken(token)
		if err != nil {

			apiG.Response(http.StatusOK, 400001005, nil)
			return
		}

		userIdInt, err := strconv.Atoi(userId)
		if err != nil {

			apiG.Response(http.StatusOK, 400001005, nil)
			return
		}
		c.Set("userId", userIdInt)
		c.Set("token", token)
		//if routerAsName == "" {
		//	apiG.Response(http.StatusOK,0,nil)
		//	return
		//}
		c.Next()
	}
}
