package tools

import (
	"c2v2/internal/pkg/render"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Base64Tool struct {
	Render *render.Helper
}

func NewBase64Tool(r *render.Helper) *Base64Tool {
	return &Base64Tool{Render: r}
}

func (t *Base64Tool) Handler(c *gin.Context) {
	lang := c.GetString("lang")
	if lang == "" {
		lang = "en"
	}

	// 完整页面渲染
	
	// 1. SoftwareApplication Schema
	appSchema := map[string]any{
		"@type":               "SoftwareApplication",
		"name":                t.Render.Translate(lang, "tool_base64_title"),
		"applicationCategory": "DeveloperApplication",
		"operatingSystem":     "Web",
		"offers": map[string]string{
			"@type": "Offer",
			"price": "0",
		},
		"description": t.Render.Translate(lang, "tool_base64_desc"),
	}

	// 2. FAQPage Schema
	faqSchema := map[string]any{
		"@type":    "FAQPage",
		"mainEntity": []map[string]any{
			{
				"@type": "Question",
				"name":  t.Render.Translate(lang, "b64_seo_faq_1_q"),
				"acceptedAnswer": map[string]any{
					"@type": "Answer",
					"text":  t.Render.Translate(lang, "b64_seo_faq_1_a"),
				},
			},
			{
				"@type": "Question",
				"name":  t.Render.Translate(lang, "b64_seo_faq_2_q"),
				"acceptedAnswer": map[string]any{
					"@type": "Answer",
					"text":  t.Render.Translate(lang, "b64_seo_faq_2_a"),
				},
			},
		},
	}
	
	// 组合为 Graph Schema
	graphSchema := map[string]any{
		"@context": "https://schema.org",
		"@graph":   []any{appSchema, faqSchema},
	}

	t.Render.HTML(c, http.StatusOK, "base64.html", gin.H{
		"title":       "tool_base64_page_title",
		"description": "tool_base64_page_desc",
		"SchemaData":  graphSchema,
	})
}