package app

import (
	"c2v2/internal/middleware"
	"c2v2/internal/pkg/i18n"
	"c2v2/internal/pkg/render"
	"c2v2/internal/tools"
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func SetupRouter(i18nMgr *i18n.Manager, cfg *Config) *gin.Engine {
	r := gin.Default()

	// 全局中间件
	r.Use(gzip.Gzip(gzip.DefaultCompression)) // Gzip 压缩
	r.Use(middleware.SecurityHeaders())       // 安全响应头
	r.Use(middleware.CacheControl())          // 缓存控制

	// 从配置读取域名
	domain := cfg.Domain
	renderHelper := render.NewHelper(i18nMgr, domain)

	r.Static("/static", "./static")

	// 注册自定义模板函数
	r.SetFuncMap(template.FuncMap{
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, nil // should return error in real world
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, nil
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"list": func(values ...interface{}) []interface{} {
			return values
		},
		"split": func(s, sep string) []string {
			return strings.Split(s, sep)
		},
		"len": func(v interface{}) int {
			switch val := v.(type) {
			case string:
				return len(val)
			case []string:
				return len(val)
			case []interface{}:
				return len(val)
			default:
				return 0
			}
		},
		"loop": func(n int) []int {
			result := make([]int, n)
			for i := 0; i < n; i++ {
				result[i] = i
			}
			return result
		},
		"add": func(a, b int) int {
			return a + b
		},
		"eq": func(a, b interface{}) bool {
			return a == b
		},
		"not": func(b bool) bool {
			return !b
		},
	})

	r.LoadHTMLGlob("templates/**/*")

	// 注册 Sitemap
	r.GET("/sitemap.xml", SitemapHandler(domain))

	// 注册 robots.txt
	r.GET("/robots.txt", func(c *gin.Context) {
		robotsTxt := `User-agent: *
Allow: /

Sitemap: ` + domain + `/sitemap.xml
`
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, robotsTxt)
	})

	base64Tool := tools.NewBase64Tool(renderHelper)
	jsonTool := tools.NewJsonFmtTool(renderHelper)
	htmlTool := tools.NewHTMLFmtTool(renderHelper)
	cssTool := tools.NewCSSFmtTool(renderHelper)
	heicTool := tools.NewHeicTool(renderHelper)

	// 从统一注册中心获取工具数据
	categories := tools.Categories()
	allTools := tools.AllTools()

	// 搜索处理函数
	searchHandler := func(c *gin.Context) {
		query := c.Query("q")
		lang := c.GetString("lang")

		// 如果没有搜索词，返回分类列表视图（还原首页状态）
		if query == "" {
			renderHelper.HTML(c, http.StatusOK, "category_list", gin.H{
				"Categories": categories,
			})
			return
		}

		// 简单的模糊匹配
		var results []tools.Tool
		query = strings.ToLower(query)

		for _, tool := range allTools {
			name := i18nMgr.Translate(lang, tool.NameKey)
			desc := i18nMgr.Translate(lang, tool.DescKey)

			if strings.Contains(strings.ToLower(name), query) ||
				strings.Contains(strings.ToLower(desc), query) {
				results = append(results, tool)
			}
		}

		renderHelper.HTML(c, http.StatusOK, "tool_list", gin.H{
			"Tools": results,
		})
	}

	// 1. 默认语言路由 (英语)
	defaultGroup := r.Group("")
	defaultGroup.Use(middleware.I18nMiddleware())
	{
		// 首页 /
		defaultGroup.GET("/", func(c *gin.Context) {
			renderHelper.HTML(c, http.StatusOK, "index.html", gin.H{
				"title":       "meta_title_index",
				"description": "meta_desc_index",
				"Categories":  categories, // 传递分类列表
			})
		})

		// 搜索 API
		defaultGroup.GET("/api/search", searchHandler)

		// 工具路由
		defaultGroup.GET("/base64", base64Tool.Handler)
		defaultGroup.POST("/base64", base64Tool.Handler)
		defaultGroup.GET("/json-fmt", jsonTool.Handler)
		defaultGroup.POST("/json-fmt", jsonTool.Handler)
		defaultGroup.GET("/html-fmt", htmlTool.Handler)
		defaultGroup.POST("/html-fmt", htmlTool.Handler)
		defaultGroup.GET("/css-fmt", cssTool.Handler)
		defaultGroup.POST("/css-fmt", cssTool.Handler)
		defaultGroup.GET("/heic-to-jpg", heicTool.Handler)
		defaultGroup.POST("/heic-to-jpg", heicTool.Handler)

		// 静态页面路由
		defaultGroup.GET("/about", func(c *gin.Context) {
			renderHelper.HTML(c, http.StatusOK, "about.html", gin.H{"title": "nav_about", "description": "about_desc"})
		})
		defaultGroup.GET("/privacy", func(c *gin.Context) {
			renderHelper.HTML(c, http.StatusOK, "privacy.html", gin.H{"title": "nav_privacy", "description": "privacy_desc"})
		})
		defaultGroup.GET("/terms", func(c *gin.Context) {
			renderHelper.HTML(c, http.StatusOK, "terms.html", gin.H{"title": "nav_terms", "description": "terms_desc"})
		})
		defaultGroup.GET("/contact", func(c *gin.Context) {
			renderHelper.HTML(c, http.StatusOK, "contact.html", gin.H{"title": "nav_contact", "description": "contact_desc"})
		})
	}

	// 2. 其他语言路由 (例如 /zh/...)
	langGroup := r.Group("/:lang")
	langGroup.Use(middleware.I18nMiddleware())
	{
		// 语言首页
		langGroup.GET("/", func(c *gin.Context) {
			renderHelper.HTML(c, http.StatusOK, "index.html", gin.H{
				"title":       "meta_title_index",
				"description": "meta_desc_index",
				"Categories":  categories,
			})
		})

		// 搜索 API
		langGroup.GET("/api/search", searchHandler)

		// 工具路由
		langGroup.GET("/base64", base64Tool.Handler)
		langGroup.POST("/base64", base64Tool.Handler)
		langGroup.GET("/json-fmt", jsonTool.Handler)
		langGroup.POST("/json-fmt", jsonTool.Handler)
		langGroup.GET("/html-fmt", htmlTool.Handler)
		langGroup.POST("/html-fmt", htmlTool.Handler)
		langGroup.GET("/css-fmt", cssTool.Handler)
		langGroup.POST("/css-fmt", cssTool.Handler)
		langGroup.GET("/heic-to-jpg", heicTool.Handler)
		langGroup.POST("/heic-to-jpg", heicTool.Handler)

		// 静态页面路由
		langGroup.GET("/about", func(c *gin.Context) {
			renderHelper.HTML(c, http.StatusOK, "about.html", gin.H{"title": "nav_about", "description": "about_desc"})
		})
		langGroup.GET("/privacy", func(c *gin.Context) {
			renderHelper.HTML(c, http.StatusOK, "privacy.html", gin.H{"title": "nav_privacy", "description": "privacy_desc"})
		})
		langGroup.GET("/terms", func(c *gin.Context) {
			renderHelper.HTML(c, http.StatusOK, "terms.html", gin.H{"title": "nav_terms", "description": "terms_desc"})
		})
		langGroup.GET("/contact", func(c *gin.Context) {
			renderHelper.HTML(c, http.StatusOK, "contact.html", gin.H{"title": "nav_contact", "description": "contact_desc"})
		})
	}

	return r
}
