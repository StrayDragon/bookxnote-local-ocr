package service

import (
	"log"

	"github.com/straydragon/bookxnote-local-ocr/internal/lib/umiocr"
)

// Service 包含所有服务依赖
type Service struct {
	UmiOCRClient *umiocr.Client
	configMgr    *ConfigManager
}

// NewService 创建新的服务实例
func NewService(ocrClient *umiocr.Client) *Service {
	configMgr, err := NewConfigManager()
	if err != nil {
		log.Printf("Failed to load config: %v, using default settings", err)
	}

	return &Service{
		UmiOCRClient: ocrClient,
		configMgr:    configMgr,
	}
}

// GetConfigValue 获取配置值
func (s *Service) GetConfigValue(key string) interface{} {
	return s.configMgr.Get(key)
}

// SetConfigValue 设置配置值
func (s *Service) SetConfigValue(key string, value interface{}) error {
	return s.configMgr.Set(key, value)
}
