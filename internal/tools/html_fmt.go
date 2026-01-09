package tools

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// HTMLFormatResult represents the result of HTML formatting
type HTMLFormatResult struct {
	Formatted string
	Error     string
	Line      int
	Column    int
}

// FormatHTML formats and prettifies HTML code
func FormatHTML(input string) HTMLFormatResult {
	if strings.TrimSpace(input) == "" {
		return HTMLFormatResult{
			Formatted: "",
			Error:     "",
		}
	}

	// Parse the HTML
	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		// Extract line and column from error
		line, column := extractErrorPosition(err.Error())
		return HTMLFormatResult{
			Formatted: "",
			Error:     fmt.Sprintf("HTML解析错误: %s", err.Error()),
			Line:      line,
			Column:    column,
		}
	}

	// Format the HTML
	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		line, column := extractErrorPosition(err.Error())
		return HTMLFormatResult{
			Formatted: "",
			Error:     fmt.Sprintf("HTML渲染错误: %s", err.Error()),
			Line:      line,
			Column:    column,
		}
	}

	// Pretty print with indentation - No longer needed as we do it in frontend
	// Keeping a minimal version for validation support if needed
	return HTMLFormatResult{
		Formatted: "",
		Error:     "",
	}
}

// extractErrorPosition extracts line and column from error message
func extractErrorPosition(errorMsg string) (int, int) {
	// Look for patterns like "line 1, column 2" or "1:2"
	patterns := []string{
		`line (\d+), column (\d+)`,
		`(\d+):(\d+)`,
		`position (\d+):(\d+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(errorMsg)
		if len(matches) >= 3 {
			line := 0
			column := 0
			fmt.Sscanf(matches[1], "%d", &line)
			fmt.Sscanf(matches[2], "%d", &column)
			return line, column
		}
	}

	return 0, 0
}
