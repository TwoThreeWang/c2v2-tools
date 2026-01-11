package tools

import (
	"c2v2/internal/pkg/render"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PasswordTool handles Password Generator page requests
type PasswordTool struct {
	renderHelper *render.Helper
}

// NewPasswordTool creates a new Password Generator tool
func NewPasswordTool(renderHelper *render.Helper) *PasswordTool {
	return &PasswordTool{
		renderHelper: renderHelper,
	}
}

// Handler handles HTTP requests for the Password Generator tool
func (t *PasswordTool) Handler(c *gin.Context) {
	lang := c.GetString("lang")
	if lang == "" {
		lang = "en"
	}

	// Create SEO schema data
	appSchema := map[string]any{
		"@type":               "SoftwareApplication",
		"name":                t.renderHelper.Translate(lang, "tool_password_title"),
		"applicationCategory": "SecurityApplication",
		"operatingSystem":     "Web",
		"offers": map[string]string{
			"@type": "Offer",
			"price": "0",
		},
		"description": t.renderHelper.Translate(lang, "tool_password_desc"),
	}

	faqSchema := map[string]any{
		"@type": "FAQPage",
		"mainEntity": []map[string]any{
			{
				"@type": "Question",
				"name":  t.renderHelper.Translate(lang, "pwd_seo_faq_1_q"),
				"acceptedAnswer": map[string]any{
					"@type": "Answer",
					"text":  t.renderHelper.Translate(lang, "pwd_seo_faq_1_a"),
				},
			},
		},
	}

	graphSchema := map[string]any{
		"@context": "https://schema.org",
		"@graph":   []any{appSchema, faqSchema},
	}

	t.renderHelper.HTML(c, http.StatusOK, "password_generator.html", gin.H{
		"title":       "tool_password_page_title",
		"description": "tool_password_page_desc",
		"keywords":    "tool_password_keywords",
		"SchemaData":  graphSchema,
	})
}
