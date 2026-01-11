package tools

import (
	"c2v2/internal/pkg/render"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HeicTool struct {
	Render *render.Helper
}

func NewHeicTool(r *render.Helper) *HeicTool {
	return &HeicTool{Render: r}
}

func (t *HeicTool) Handler(c *gin.Context) {
	lang := c.GetString("lang")
	if lang == "" {
		lang = "en"
	}

	// 1. SoftwareApplication Schema
	appSchema := map[string]any{
		"@type":               "SoftwareApplication",
		"name":                t.Render.Translate(lang, "tool_heic_title"),
		"applicationCategory": "MultimediaApplication",
		"operatingSystem":     "Web",
		"offers": map[string]string{
			"@type": "Offer",
			"price": "0",
		},
		"description": t.Render.Translate(lang, "tool_heic_desc"),
	}

	// 2. FAQPage Schema
	faqSchema := map[string]any{
		"@type": "FAQPage",
		"mainEntity": []map[string]any{
			{
				"@type": "Question",
				"name":  t.Render.Translate(lang, "heic_seo_faq_1_q"),
				"acceptedAnswer": map[string]any{
					"@type": "Answer",
					"text":  t.Render.Translate(lang, "heic_seo_faq_1_a"),
				},
			},
			{
				"@type": "Question",
				"name":  t.Render.Translate(lang, "heic_seo_faq_2_q"),
				"acceptedAnswer": map[string]any{
					"@type": "Answer",
					"text":  t.Render.Translate(lang, "heic_seo_faq_2_a"),
				},
			},
		},
	}

	// 3. HowTo Schema (Conversion Steps)
	howToSchema := map[string]any{
		"@type": "HowTo",
		"name":  t.Render.Translate(lang, "heic_seo_h2_how"),
		"step": []map[string]any{
			{
				"@type": "HowToStep",
				"text":  t.Render.Translate(lang, "heic_seo_step_1"),
				"name":  t.Render.Translate(lang, "heic_seo_step_1_title"),
			},
			{
				"@type": "HowToStep",
				"text":  t.Render.Translate(lang, "heic_seo_step_2"),
				"name":  t.Render.Translate(lang, "heic_seo_step_2_title"),
			},
			{
				"@type": "HowToStep",
				"text":  t.Render.Translate(lang, "heic_seo_step_3"),
				"name":  t.Render.Translate(lang, "heic_seo_step_3_title"),
			},
		},
	}

	// Combined Schema
	graphSchema := map[string]any{
		"@context": "https://schema.org",
		"@graph":   []any{appSchema, faqSchema, howToSchema},
	}

	t.Render.HTML(c, http.StatusOK, "heic.html", gin.H{
		"title":       "tool_heic_page_title",
		"description": "tool_heic_page_desc",
		"SchemaData":  graphSchema,
		"ICON_HOWTO":  "M13 10V3L4 14h7v7l9-11h-7z",
		"ICON_FAQ":    "M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z",
	})
}
