package service

import (
	"context"
	"log"

	"github.com/straydragon/bookxnote-local-ocr/internal/lib/umiocr"
)

// Service 包含所有服务依赖
type Service struct {
	UmiOCRClient     *umiocr.Client
	configMgr        *ConfigManager
	postOCRProcessor *PostOCRProcessor
}

// NewService 创建新的服务实例
func NewService(ocrClient *umiocr.Client) *Service {
	configMgr, err := NewConfigManager()
	if err != nil {
		log.Printf("Failed to load config: %v, using default settings", err)
	}

	config, err := configMgr.GetConfig()
	if err != nil {
		log.Printf("解析配置失败 > %v", err)
	}

	var postOCRProcessor *PostOCRProcessor
	if config != nil && config.AfterOCR != nil && config.LLM != nil {
		postOCRProcessor, err = NewPostOCRProcessor()
		if err != nil {
			log.Printf("创建后处理处理器失败 > %v", err)
		}
	}
	return &Service{
		UmiOCRClient:     ocrClient,
		configMgr:        configMgr,
		postOCRProcessor: postOCRProcessor,
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

// ProcessOCRResult 处理OCR结果
func (s *Service) ProcessOCRResult(ctx context.Context, ocrResult *umiocr.APIRecognizeResp) (*umiocr.APIRecognizeResp, error) {
	if s.postOCRProcessor == nil {
		return ocrResult, nil
	}
	return s.postOCRProcessor.ProcessOCRResult(ctx, ocrResult)
}
