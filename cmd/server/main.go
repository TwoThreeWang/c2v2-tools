package main

import (
	"c2v2/internal/app"
	"c2v2/internal/pkg/i18n"
	"log"
)

func main() {
	// 初始化 I18n
	i18nMgr := i18n.NewManager()
	if err := i18nMgr.LoadTranslations("locales"); err != nil {
		log.Fatalf("Failed to load translations: %v", err)
	}

	r := app.SetupRouter(i18nMgr)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}