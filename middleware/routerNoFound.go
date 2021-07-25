package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func CheckExist() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		s := strings.Contains(path, "/backend/")
		c.Next()
		status := c.Writer.Status()
		if status == 404 {
			if s {
				c.Redirect(http.StatusMovedPermanently, "/backend/")
			} else {
				c.Redirect(http.StatusMovedPermanently, "/404")
			}
		}
	}
}
