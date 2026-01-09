package app

// ToolRoutes 定义了所有工具的路径后缀
// 不包含语言前缀，例如 "base64" (对应 /base64 和 /zh/base64)
// 空字符串 "" 代表首页
var ToolRoutes = []string{
	"",         // 首页
	"base64",   // Base64 工具
	"json-fmt", // JSON 格式化 (预留)
}
