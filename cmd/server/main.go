package main

import (
	"c2v2/internal/app"
	"c2v2/internal/pkg/i18n"
	"log"
)

func main() {
	// 加载配置
	cfg := app.LoadConfig()

	// 初始化 I18n
	i18nMgr := i18n.NewManager()
	if err := i18nMgr.LoadTranslations("locales"); err != nil {
		log.Fatalf("加载翻译文件失败: %v", err)
	}

	r := app.SetupRouter(i18nMgr, cfg)
	log.Printf("服务器启动在 :%s (域名: %s)", cfg.Port, cfg.Domain)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
