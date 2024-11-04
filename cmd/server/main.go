package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/straydragon/bookxnote-local-ocr/internal/common/settings"
	"github.com/straydragon/bookxnote-local-ocr/internal/handlers"
	"github.com/straydragon/bookxnote-local-ocr/internal/lib/umiocr"
	"github.com/straydragon/bookxnote-local-ocr/internal/middleware"
	"github.com/straydragon/bookxnote-local-ocr/internal/service"
)

func main() {
	// 加载配置
	config, err := service.LoadConfig()
	if err != nil {
		log.Printf("Failed to load config: %v, using default settings", err)
		config = &service.DefaultConfig
	}

	// 初始化服务
	ocrClient := umiocr.NewClient(config.OCR.UmiOCR.APIURL)
	svc := service.NewService(ocrClient)

	// 设置gin
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// 注入服务实例到所有路由
	r.Use(middleware.InjectService(svc))

	// 注册路由
	r.POST("/oauth/2.0/token", handlers.TokenHandler)
	r.POST("/rest/2.0/ocr/v1/accurate_basic", handlers.AccurateOCRHandler)
	r.NoRoute(handlers.CatchAllHandler)

	// 启动服务器
	r.RunTLS(":443", settings.GetPathFromCertDir("cert.pem"), settings.GetPathFromCertDir("key.pem"))
}
