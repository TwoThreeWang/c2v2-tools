package tools

import (
	"html/template"
)

// Tool 定义工具列表项的数据结构
type Tool struct {
	ID       string
	NameKey  string        // 对应 en.json/zh.json 中的 key
	DescKey  string        // 对应 key
	URL      string        // 相对 URL 路径
	IconHTML template.HTML // SVG 图标，自动视为安全 HTML
}

// Category 定义工具分类
type Category struct {
	ID      string
	NameKey string
	DescKey string
	Tools   []Tool
}

// 定义所有工具
var (
	ToolBase64 = Tool{
		ID:       "base64",
		NameKey:  "tool_base64_title",
		DescKey:  "tool_base64_desc",
		URL:      "/base64",
		IconHTML: template.HTML(`<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4"></path></svg>`),
	}

	ToolJSON = Tool{
		ID:       "json-fmt",
		NameKey:  "tool_json_title",
		DescKey:  "tool_json_desc",
		URL:      "/json-fmt",
		IconHTML: template.HTML(`<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16m-7 6h7"></path></svg>`),
	}

	ToolHTML = Tool{
		ID:       "html-fmt",
		NameKey:  "tool_html_title",
		DescKey:  "tool_html_desc",
		URL:      "/html-fmt",
		IconHTML: template.HTML(`<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"></path></svg>`),
	}

	ToolCSS = Tool{
		ID:       "css-fmt",
		NameKey:  "tool_css_title",
		DescKey:  "tool_css_desc",
		URL:      "/css-fmt",
		IconHTML: template.HTML(`<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01"></path></svg>`),
	}
)

// Categories 返回所有工具分类（用于首页渲染）
func Categories() []Category {
	return []Category{
		{
			ID:      "encoders",
			NameKey: "cat_encoders_title",
			DescKey: "cat_encoders_desc",
			Tools:   []Tool{ToolBase64},
		},
		{
			ID:      "formatters",
			NameKey: "cat_formatters_title",
			DescKey: "cat_formatters_desc",
			Tools:   []Tool{ToolJSON, ToolHTML, ToolCSS},
		},
	}
}

// AllTools 返回所有工具的扁平列表（用于搜索）
func AllTools() []Tool {
	return []Tool{ToolBase64, ToolJSON, ToolHTML, ToolCSS}
}

// AllRoutes 返回所有需要包含在 Sitemap 中的路由
// 包括工具页面和静态页面
func AllRoutes() []string {
	return []string{
		"",         // 首页
		"base64",   // Base64 工具
		"json-fmt", // JSON 格式化
		"html-fmt", // HTML 格式化
		"css-fmt",  // CSS 格式化
		"about",    // 关于页面
		"privacy",  // 隐私政策
		"terms",    // 服务条款
		"contact",  // 联系我们
	}
}
