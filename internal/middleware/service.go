package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/straydragon/bookxnote-local-ocr/internal/service"
)

// InjectService 注入服务实例到上下文
func InjectService(svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("service", svc)
		c.Next()
	}
}
