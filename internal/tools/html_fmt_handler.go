package tools

import (
	"c2v2/internal/pkg/render"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HTMLFmtTool handles HTML formatting operations
type HTMLFmtTool struct {
	renderHelper *render.Helper
}

// NewHTMLFmtTool creates a new HTML formatting tool
func NewHTMLFmtTool(renderHelper *render.Helper) *HTMLFmtTool {
	return &HTMLFmtTool{
		renderHelper: renderHelper,
	}
}

// Handler handles HTTP requests for the HTML formatting tool
func (t *HTMLFmtTool) Handler(c *gin.Context) {
	action := c.PostForm("action")
	input := c.PostForm("input")
	lang := c.GetString("lang")
	if lang == "" {
		lang = "en"
	}

	// Handle GET request (show the form)
	if c.Request.Method == "GET" {
		// Create SEO schema data
		appSchema := map[string]any{
			"@type":               "SoftwareApplication",
			"name":                t.renderHelper.Translate(lang, "tool_html_title"),
			"applicationCategory": "DeveloperApplication",
			"operatingSystem":     "Web",
			"offers": map[string]string{
				"@type": "Offer",
				"price": "0",
			},
			"description": t.renderHelper.Translate(lang, "tool_html_desc"),
		}

		faqSchema := map[string]any{
			"@type": "FAQPage",
			"mainEntity": []map[string]any{
				{
					"@type": "Question",
					"name":  t.renderHelper.Translate(lang, "html_seo_faq_1_q"),
					"acceptedAnswer": map[string]any{
						"@type": "Answer",
						"text":  t.renderHelper.Translate(lang, "html_seo_faq_1_a"),
					},
				},
				{
					"@type": "Question",
					"name":  t.renderHelper.Translate(lang, "html_seo_faq_2_q"),
					"acceptedAnswer": map[string]any{
						"@type": "Answer",
						"text":  t.renderHelper.Translate(lang, "html_seo_faq_2_a"),
					},
				},
			},
		}

		graphSchema := map[string]any{
			"@context": "https://schema.org",
			"@graph":   []any{appSchema, faqSchema},
		}

		t.renderHelper.HTML(c, http.StatusOK, "html_fmt.html", gin.H{
			"title":       "tool_html_page_title",
			"description": "tool_html_page_desc",
			"keywords":    "tool_html_keywords",
			"SchemaData":  graphSchema,
		})
		return
	}

	// Handle POST request (process the form)
	if action == "validate" {
		result := FormatHTML(input)
		if result.Error != "" {
			c.HTML(http.StatusOK, "html_fmt_result.html", gin.H{
				"error":  result.Error,
				"line":   result.Line,
				"column": result.Column,
			})
		} else {
			c.HTML(http.StatusOK, "html_fmt_result.html", gin.H{
				"error": "",
			})
		}
		return
	}

	c.Status(http.StatusBadRequest)
}
