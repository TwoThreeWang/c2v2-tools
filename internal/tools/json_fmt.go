package tools

import (
	"bytes"
	"c2v2/internal/pkg/render"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

// jsonToGoStruct converts JSON data to Go struct definition
func jsonToGoStruct(data interface{}, structName string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("type %s struct {\n", structName))
	generateGoStruct(&sb, data, 1, make(map[string]bool))
	sb.WriteString("}")
	return sb.String()
}

func generateGoStruct(sb *strings.Builder, data interface{}, indent int, visited map[string]bool) {
	indentStr := strings.Repeat("\t", indent)

	switch v := data.(type) {
	case map[string]interface{}:
		// Handle object
		for key, value := range v {
			fieldName := toPascalCase(key)
			fieldType := getGoType(value, visited)
			sb.WriteString(fmt.Sprintf("%s%s %s `json:\"%s\"`\n", indentStr, fieldName, fieldType, key))
		}
	case []interface{}:
		// Handle array - just process first element for type
		if len(v) > 0 {
			elemType := getGoType(v[0], visited)
			sb.WriteString(fmt.Sprintf("%sItems []%s `json:\"items\"`\n", indentStr, elemType))
		}
	default:
		// Handle primitive types
		typeStr := getGoType(data, visited)
		sb.WriteString(fmt.Sprintf("%sValue %s `json:\"value\"`\n", indentStr, typeStr))
	}
}

func getGoType(data interface{}, visited map[string]bool) string {
	switch v := data.(type) {
	case string:
		return "string"
	case float64:
		if v == float64(int64(v)) {
			return "int"
		}
		return "float64"
	case bool:
		return "bool"
	case map[string]interface{}:
		// For nested objects, create inline struct
		var sb strings.Builder
		sb.WriteString("struct {\n")
		generateGoStruct(&sb, v, 1, visited)
		sb.WriteString("\t}")
		return sb.String()
	case []interface{}:
		if len(v) > 0 {
			elemType := getGoType(v[0], visited)
			return "[]" + elemType
		}
		return "[]interface{}"
	case nil:
		return "interface{}"
	default:
		return "interface{}"
	}
}

// parseJSONError解析JSON解析错误，提取行号和列号信息
func parseJSONError(err error, input string) (string, int, int) {
	errorMsg := err.Error()

	// 尝试从错误消息中提取位置信息
	// 常见的JSON错误格式: "invalid character ',' after array element at line 3, column 15"
	re := regexp.MustCompile(`at line (\d+), column (\d+)`)
	matches := re.FindStringSubmatch(errorMsg)

	line := 0
	column := 0

	if len(matches) >= 3 {
		if lineNum, err := strconv.Atoi(matches[1]); err == nil {
			line = lineNum
		}
		if colNum, err := strconv.Atoi(matches[2]); err == nil {
			column = colNum
		}
	}

	// 如果没有找到行号，尝试计算
	if line == 0 && strings.Contains(errorMsg, "line") {
		// 简单的行号检测
		lines := strings.Split(input, "\n")
		for i, l := range lines {
			if strings.Contains(l, "{") || strings.Contains(l, "[") {
				line = i + 1
				break
			}
		}
	}

	return errorMsg, line, column
}

// getErrorContext获取错误行的上下文
func getErrorContext(input string, errorLine int) string {
	lines := strings.Split(input, "\n")
	if errorLine <= 0 || errorLine > len(lines) {
		return ""
	}

	// 显示错误行和前后几行
	start := max(0, errorLine-2)
	end := min(len(lines), errorLine+1)

	var context strings.Builder
	for i := start; i < end; i++ {
		lineNum := i + 1
		marker := "  "
		if lineNum == errorLine {
			marker = "> "
		}
		context.WriteString(fmt.Sprintf("%s%d: %s\n", marker, lineNum, lines[i]))
	}

	return context.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func toPascalCase(s string) string {
	if s == "" {
		return ""
	}

	// Split by common separators
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})

	var result strings.Builder
	for _, part := range parts {
		if len(part) > 0 {
			result.WriteString(strings.ToUpper(part[:1]))
			if len(part) > 1 {
				result.WriteString(strings.ToLower(part[1:]))
			}
		}
	}

	return result.String()
}

type JsonFmtTool struct {
	Render *render.Helper
}

func NewJsonFmtTool(r *render.Helper) *JsonFmtTool {
	return &JsonFmtTool{Render: r}
}

func (t *JsonFmtTool) Handler(c *gin.Context) {
	lang := c.GetString("lang")
	if lang == "" {
		lang = "en"
	}

	// HTMX 请求处理
	if c.GetHeader("HX-Request") == "true" {
		input := c.PostForm("input")
		action := c.PostForm("action")
		var result string
		var isError bool

		input = strings.TrimSpace(input)

		if input == "" {
			result = t.Render.Translate(lang, "json_error_empty")
			isError = true
		} else {
			// 尝试解析 JSON
			var jsonObj interface{}
			err := json.Unmarshal([]byte(input), &jsonObj)
			if err != nil {
				// 解析错误，获取详细的错误信息
				errorMsg, errorLine, errorColumn := parseJSONError(err, input)
				context := getErrorContext(input, errorLine)

				var errorDetails strings.Builder
				errorDetails.WriteString(t.Render.Translate(lang, "json_error_invalid"))
				errorDetails.WriteString(errorMsg)

				if errorLine > 0 {
					errorDetails.WriteString(fmt.Sprintf("\n\n%s %d, %s %d",
						t.Render.Translate(lang, "json_error_line"), errorLine,
						t.Render.Translate(lang, "json_error_column"), errorColumn))
					if context != "" {
						errorDetails.WriteString(fmt.Sprintf("\n\n%s:\n", t.Render.Translate(lang, "json_error_context")))
						errorDetails.WriteString(context)
					}
				}

				result = errorDetails.String()
				isError = true
			} else {
				switch action {
				case "minify":
					// 压缩
					compactedBuffer := new(bytes.Buffer)
					if err := json.Compact(compactedBuffer, []byte(input)); err != nil {
						result = "Error minifying JSON"
						isError = true
					} else {
						result = compactedBuffer.String()
					}
				case "validate":
					// 验证 - 返回成功消息
					result = "✓ Valid JSON"
					isError = false
				case "to_go":
					// JSON 转 Go Struct
					result = jsonToGoStruct(jsonObj, "AutoGenerated")
				case "to_yaml":
					// JSON 转 YAML
					yamlData, err := yaml.Marshal(jsonObj)
					if err != nil {
						result = "Error converting to YAML: " + err.Error()
						isError = true
					} else {
						result = string(yamlData)
					}
				default:
					// 格式化 (默认)
					formatted, _ := json.MarshalIndent(jsonObj, "", "  ")
					result = string(formatted)
				}
			}
		}

		// 计算统计信息
		charCount := utf8.RuneCountInString(result)
		byteCount := len(result)

		t.Render.HTML(c, http.StatusOK, "json_fmt_result.html", gin.H{
			"result":    result,
			"isError":   isError,
			"charCount": charCount,
			"byteCount": byteCount,
		})
		return
	}

	// 完整页面渲染 Schema
	appSchema := map[string]any{
		"@type":               "SoftwareApplication",
		"name":                t.Render.Translate(lang, "tool_json_title"),
		"applicationCategory": "DeveloperApplication",
		"operatingSystem":     "Web",
		"offers": map[string]string{
			"@type": "Offer",
			"price": "0",
		},
		"description": t.Render.Translate(lang, "tool_json_desc"),
	}

	faqSchema := map[string]any{
		"@type": "FAQPage",
		"mainEntity": []map[string]any{
			{
				"@type": "Question",
				"name":  t.Render.Translate(lang, "json_seo_faq_1_q"),
				"acceptedAnswer": map[string]any{
					"@type": "Answer",
					"text":  t.Render.Translate(lang, "json_seo_faq_1_a"),
				},
			},
			{
				"@type": "Question",
				"name":  t.Render.Translate(lang, "json_seo_faq_2_q"),
				"acceptedAnswer": map[string]any{
					"@type": "Answer",
					"text":  t.Render.Translate(lang, "json_seo_faq_2_a"),
				},
			},
		},
	}

	graphSchema := map[string]any{
		"@context": "https://schema.org",
		"@graph":   []any{appSchema, faqSchema},
	}

	t.Render.HTML(c, http.StatusOK, "json_fmt.html", gin.H{
		"title":       "tool_json_page_title",
		"description": "tool_json_page_desc",
		"SchemaData":  graphSchema,
	})
}
