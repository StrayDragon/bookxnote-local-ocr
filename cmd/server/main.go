package main

import (
	"log"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/straydragon/bookxnote-local-ocr/internal/common/settings"
	"github.com/straydragon/bookxnote-local-ocr/internal/common/utils"
	"github.com/straydragon/bookxnote-local-ocr/internal/handlers"
	"github.com/straydragon/bookxnote-local-ocr/internal/middleware"
	"github.com/straydragon/bookxnote-local-ocr/internal/service"
	_ "github.com/straydragon/bookxnote-local-ocr/internal/swagger-doc"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title BookxNote Local OCR API
// @version 1.0
// @description This is a local OCR service for BookxNote that mimics Baidu OCR API
// @contact.name API Support
// @contact.url https://github.com/straydragon/bookxnote-local-ocr
// @license.name MIT
// @license.url https://github.com/straydragon/bookxnote-local-ocr/blob/main/LICENSE
// @host localhost:443
// @BasePath /
// @schemes https
func main() {
	// 权限检查
	if runtime.GOOS == "linux" {
		if err := utils.CheckCurrentProcessCaps([]string{"cap_net_bind_service"}); err != nil {
			log.Fatalf("%v", err)
		}
	}

	// 初始化服务
	svc, err := service.NewService()
	if err != nil {
		log.Fatalf("Failed to initialize service: %v", err)
	}

	// 设置gin
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// 注入服务实例到所有路由
	r.Use(middleware.InjectService(svc))

	// 注册路由: Hook 百度OCR相关API
	r.POST("/oauth/2.0/token", handlers.TokenHandler)
	r.POST("/rest/2.0/ocr/v1/accurate_basic", handlers.AccurateOCRHandler)

	// 注册内部服务路由: 用于实现其他附加功能
	rApp := r.Group("/_app")
	rAppConfig := rApp.Group("/config")
	{
		rAppConfig.GET("/Get", handlers.AppConfigGetHandler)
		rAppConfig.POST("/Set", handlers.AppConfigSetHandler)
	}

	// 添加Swagger文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.NoRoute(handlers.CatchAllHandler)
	// 启动服务器
	certPaths := settings.GetPathsFromCertDir("cert.pem", "key.pem")
	r.RunTLS(":443", certPaths[0], certPaths[1])
}
