package app

import (
	"c2v2/internal/pkg/i18n"

	"github.com/gin-gonic/gin"
)

// RenderHelper 封装了渲染逻辑，自动注入 T 函数和常用变量
type RenderHelper struct {
	i18n *i18n.Manager
}

func NewRenderHelper(i *i18n.Manager) *RenderHelper {
	return &RenderHelper{i18n: i}
}

// HTML 替代 c.HTML，自动注入 T, lang 等
func (h *RenderHelper) HTML(c *gin.Context, code int, name string, data gin.H) {
	lang := c.GetString("lang")
	if lang == "" {
		lang = "en"
	}

	// 注入通用数据
	data["lang"] = lang
	data["T"] = func(key string) string {
		return h.i18n.Translate(lang, key)
	}

	// 注入当前 URL 用于语言切换 (简单处理)
	// 实际项目中可能需要更复杂的 URL 处理
	data["CurrentPath"] = c.Request.URL.Path

	c.HTML(code, name, data)
}
