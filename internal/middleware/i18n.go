package middleware

import (
	"github.com/gin-gonic/gin"
)

func I18nMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.Param("lang")
		// 如果 URL 中没有语言参数，或者参数不在我们支持的列表里（可选扩展），默认为 en
		if lang == "" {
			lang = "en"
		}
		
		c.Set("lang", lang)
		c.Next()
	}
}