package app

import (
	"os"
	"strings"
)

// Config 存储应用程序配置
type Config struct {
	// Domain 是应用程序的完整域名，例如 "https://www.example.com"
	Domain string
	// Port 是服务器监听的端口
	Port string
	// SupportedLangs 是支持的语言列表
	SupportedLangs []string
	// DefaultLang 是默认语言
	DefaultLang string
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Domain:         "http://localhost:5006",
		Port:           "5006",
		SupportedLangs: []string{"en", "zh"},
		DefaultLang:    "en",
	}
}

// LoadConfig 从环境变量加载配置，未设置的项使用默认值
func LoadConfig() *Config {
	cfg := DefaultConfig()

	// 从环境变量读取 DOMAIN
	if domain := os.Getenv("DOMAIN"); domain != "" {
		cfg.Domain = strings.TrimRight(domain, "/")
	}

	// 从环境变量读取 PORT
	if port := os.Getenv("PORT"); port != "" {
		cfg.Port = port
	}

	// 从环境变量读取支持的语言（逗号分隔）
	if langs := os.Getenv("SUPPORTED_LANGS"); langs != "" {
		cfg.SupportedLangs = strings.Split(langs, ",")
	}

	// 从环境变量读取默认语言
	if defaultLang := os.Getenv("DEFAULT_LANG"); defaultLang != "" {
		cfg.DefaultLang = defaultLang
	}

	return cfg
}

// IsSupportedLang 检查语言是否在支持列表中
func (c *Config) IsSupportedLang(lang string) bool {
	for _, l := range c.SupportedLangs {
		if l == lang {
			return true
		}
	}
	return false
}
