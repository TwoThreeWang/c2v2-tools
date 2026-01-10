package render

import (
	"c2v2/internal/pkg/i18n"
	"encoding/json"
	"html/template"
	"strings"

	"github.com/gin-gonic/gin"
)

type Helper struct {
	i18n   *i18n.Manager
	Domain string // 例如: "https://www.example.com"
}

func NewHelper(i *i18n.Manager, domain string) *Helper {
	return &Helper{
		i18n:   i,
		Domain: strings.TrimRight(domain, "/"),
	}
}

// Translate exposes the i18n translate function to Go code
func (h *Helper) Translate(lang, key string) string {
	return h.i18n.Translate(lang, key)
}

// HTML 替代 c.HTML，自动注入 T, lang, SEO tags 等
func (h *Helper) HTML(c *gin.Context, code int, name string, data gin.H) {
	lang := c.GetString("lang")
	if lang == "" {
		lang = "en"
	}

	if data == nil {
		data = gin.H{}
	}
	data["lang"] = lang
	data["T"] = func(key string) string {
		return h.i18n.Translate(lang, key)
	}

	// 获取当前路径（去除语言前缀的纯净路径）
	currentPath := c.Request.URL.Path
	cleanPath := currentPath
	if strings.HasPrefix(currentPath, "/zh/") {
		cleanPath = "/" + strings.TrimPrefix(currentPath, "/zh/")
	} else if currentPath == "/zh" {
		cleanPath = "/"
	}

	// 1. 生成 Canonical URL (当前页面的绝对路径)
	// 如果是 en (默认)，不带前缀；如果是其他语言，带前缀
	canonicalPath := cleanPath
	if lang != "en" {
		if cleanPath == "/" {
			canonicalPath = "/" + lang + "/"
		} else {
			canonicalPath = "/" + lang + cleanPath
		}
	}
	data["Canonical"] = h.Domain + canonicalPath

	// 2. 生成 Hreflang Map (所有语言版本的绝对路径)
	// 这里硬编码支持的语言，实际项目中可从 i18n 配置获取
	supportedLangs := []string{"en", "zh"}
	hreflangs := make(map[string]string)

	for _, l := range supportedLangs {
		path := cleanPath
		if l != "en" {
			if cleanPath == "/" {
				path = "/" + l + "/"
			} else {
				path = "/" + l + cleanPath
			}
		}
		hreflangs[l] = h.Domain + path
	}
	data["Hreflangs"] = hreflangs

	// 3. 处理 JSON-LD Schema (如果有)
	// 如果传入了 "SchemaData" (map/struct)，我们将其序列化为 JSON 字符串
	if v, ok := data["SchemaData"]; ok {
		jsonBytes, err := json.Marshal(v)
		if err == nil {
			data["SchemaJSON"] = template.JS(jsonBytes)
		}
	}

	// 链接生成助手
	data["SwitchLangLink"] = func(targetLang string) string {
		if targetLang == "en" {
			return cleanPath
		}
		if cleanPath == "/" {
			return "/" + targetLang + "/"
		}
		return "/" + targetLang + cleanPath
	}

	data["L"] = func(path string) string {
		if lang == "en" {
			return path
		}
		if path == "/" {
			return "/" + lang + "/"
		}
		return "/" + lang + path
	}

	c.HTML(code, name, data)
}
