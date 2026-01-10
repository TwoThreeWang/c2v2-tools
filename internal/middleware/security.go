package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders 返回一个添加安全响应头的中间件
// 这些头可以防止常见的 Web 安全漏洞
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 防止点击劫持攻击
		c.Header("X-Frame-Options", "SAMEORIGIN")

		// 防止 MIME 类型嗅探
		c.Header("X-Content-Type-Options", "nosniff")

		// 启用 XSS 过滤器
		c.Header("X-XSS-Protection", "1; mode=block")

		// 控制 Referrer 信息发送
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// 限制功能和 API 的使用（权限策略）
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		c.Next()
	}
}
