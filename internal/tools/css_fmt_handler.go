package tools

import (
	"c2v2/internal/pkg/render"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CSSFmtTool 处理 CSS 格式化操作
type CSSFmtTool struct {
	renderHelper *render.Helper
}

// NewCSSFmtTool 创建新的 CSS 格式化工具
func NewCSSFmtTool(renderHelper *render.Helper) *CSSFmtTool {
	return &CSSFmtTool{
		renderHelper: renderHelper,
	}
}

// Handler 处理 CSS 格式化工具的 HTTP 请求（仅渲染页面，处理在前端完成）
func (t *CSSFmtTool) Handler(c *gin.Context) {
	lang := c.GetString("lang")
	if lang == "" {
		lang = "en"
	}

	// 创建 SEO schema 数据
	appSchema := map[string]any{
		"@type":               "SoftwareApplication",
		"name":                t.renderHelper.Translate(lang, "tool_css_title"),
		"applicationCategory": "DeveloperApplication",
		"operatingSystem":     "Web",
		"offers": map[string]string{
			"@type": "Offer",
			"price": "0",
		},
		"description": t.renderHelper.Translate(lang, "tool_css_desc"),
	}

	faqSchema := map[string]any{
		"@type": "FAQPage",
		"mainEntity": []map[string]any{
			{
				"@type": "Question",
				"name":  t.renderHelper.Translate(lang, "css_seo_faq_1_q"),
				"acceptedAnswer": map[string]any{
					"@type": "Answer",
					"text":  t.renderHelper.Translate(lang, "css_seo_faq_1_a"),
				},
			},
			{
				"@type": "Question",
				"name":  t.renderHelper.Translate(lang, "css_seo_faq_2_q"),
				"acceptedAnswer": map[string]any{
					"@type": "Answer",
					"text":  t.renderHelper.Translate(lang, "css_seo_faq_2_a"),
				},
			},
		},
	}

	graphSchema := map[string]any{
		"@context": "https://schema.org",
		"@graph":   []any{appSchema, faqSchema},
	}

	t.renderHelper.HTML(c, http.StatusOK, "css_fmt.html", gin.H{
		"title":       "tool_css_page_title",
		"description": "tool_css_page_desc",
		"SchemaData":  graphSchema,
	})
}
