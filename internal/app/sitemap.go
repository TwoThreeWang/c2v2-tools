package app

import (
	"c2v2/internal/tools"
	"encoding/xml"
	"time"

	"github.com/gin-gonic/gin"
)

// URLSet 是 sitemap 的根节点
type URLSet struct {
	XMLName xml.Name `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 urlset"`
	URLs    []URL    `xml:"url"`
}

// URL 是单个链接节点
type URL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

// SitemapHandler 动态生成 sitemap.xml
func SitemapHandler(domain string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var urls []URL
		supportedLangs := []string{"en", "zh"}

		// 从统一注册中心获取所有路由
		for _, route := range tools.AllRoutes() {
			// 遍历所有语言
			for _, lang := range supportedLangs {
				// 构建完整 URL
				var fullURL string

				// 处理路径逻辑
				// route == "" 是首页
				path := route
				if path != "" {
					path = "/" + path
				}

				if lang == "en" {
					// 英文（默认）： domain + /path
					fullURL = domain + path
					if path == "" {
						fullURL = domain + "/" // 首页保留尾部斜杠
					}
				} else {
					// 其他语言： domain + /zh + /path
					fullURL = domain + "/" + lang + path
				}

				// 设置优先级和更新频率
				priority := "0.8"
				if route == "" {
					priority = "1.0" // 首页权重最高
				}

				urls = append(urls, URL{
					Loc:        fullURL,
					LastMod:    time.Now().Format("2006-01-02"), // 简化处理，实际可以使用工具的更新时间
					ChangeFreq: "weekly",
					Priority:   priority,
				})
			}
		}

		// 生成 XML
		urlSet := URLSet{URLs: urls}

		c.Header("Content-Type", "application/xml")
		output, _ := xml.MarshalIndent(urlSet, "", "  ")
		c.Writer.Write([]byte(xml.Header + string(output)))
	}
}
