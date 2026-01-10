package middleware

import (
	"github.com/gin-gonic/gin"
)

// 支持的语言列表
var supportedLangs = map[string]bool{
	"en": true,
	"zh": true,
}

// I18nMiddleware 处理多语言支持
// 逻辑：URL 有语言前缀（如 /zh/）则使用该语言，否则使用默认语言 (en)
func I18nMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.Param("lang")

		// 如果 URL 中有有效的语言参数，使用该语言
		if lang != "" && supportedLangs[lang] {
			c.Set("lang", lang)
			c.Next()
			return
		}

		// 否则使用默认语言（英文）
		c.Set("lang", "en")
		c.Next()
	}
}
