package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// CacheControl 返回一个为静态资源设置缓存头的中间件
func CacheControl() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 为静态资源设置缓存头
		if strings.HasPrefix(path, "/static/") {
			// 静态资源缓存 7 天
			c.Header("Cache-Control", "public, max-age=604800, immutable")
		}

		c.Next()
	}
}
