package main

import (
	"github.com/gin-gonic/gin"
	"github.com/straydragon/bookxnote-local-ocr/internal/handlers"
	"github.com/straydragon/bookxnote-local-ocr/internal/settings"
)

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// r.Use(middleware.RequestDebuggingLogger())
	r.POST("/oauth/2.0/token", handlers.TokenHandler)
	r.POST("/rest/2.0/ocr/v1/accurate_basic", handlers.AccurateOCRHandler)
	r.NoRoute(handlers.CatchAllHandler)

	r.RunTLS(":443", settings.GetPathFromCertDir("cert.pem"), settings.GetPathFromCertDir("key.pem"))
}
