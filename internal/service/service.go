package service

import (
	"context"
	"fmt"
	"log"

	"github.com/straydragon/bookxnote-local-ocr/internal/lib/customocr"
	"github.com/straydragon/bookxnote-local-ocr/internal/lib/ocr"
	"github.com/straydragon/bookxnote-local-ocr/internal/lib/umiocr"
)

// OCRClient 定义OCR客户端接口
type OCRClient interface {
	Recognize(base64Image string) (*ocr.OCRResult, error)
}

// Service 包含所有服务依赖
type Service struct {
	ocrClient        OCRClient
	configMgr        *ConfigManager
	postOCRProcessor *PostOCRProcessor
}

// NewService 创建新的服务实例
func NewService() (*Service, error) {
	configMgr, err := NewConfigManager()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	config, err := configMgr.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if config.OCR == nil {
		return nil, fmt.Errorf("OCR configuration is missing")
	}

	var ocrClient OCRClient
	switch config.OCR.Selected {
	case "umiocr":
		if config.OCR.UmiOCR == nil || config.OCR.UmiOCR.APIURL == "" {
			return nil, fmt.Errorf("UmiOCR configuration is missing or invalid")
		}
		ocrClient = umiocr.NewClient(config.OCR.UmiOCR.APIURL)
	case "custom":
		if config.OCR.Custom == nil || config.OCR.Custom.APIURL == "" {
			return nil, fmt.Errorf("Custom OCR configuration is missing or invalid")
		}
		ocrClient = customocr.NewClient(config.OCR.Custom.APIURL, config.OCR.Custom.APIKey)
	default:
		return nil, fmt.Errorf("unsupported OCR service: %s", config.OCR.Selected)
	}

	var postOCRProcessor *PostOCRProcessor
	if config.AfterOCR != nil && config.LLM != nil {
		postOCRProcessor, err = NewPostOCRProcessor()
		if err != nil {
			log.Printf("failed to create post-processor: %v", err)
		}
	}

	return &Service{
		ocrClient:        ocrClient,
		configMgr:        configMgr,
		postOCRProcessor: postOCRProcessor,
	}, nil
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
func (s *Service) ProcessOCRResult(ctx context.Context, ocrResult *ocr.OCRResult) (*ocr.OCRResult, error) {
	if s.postOCRProcessor == nil {
		return ocrResult, nil
	}
	return s.postOCRProcessor.ProcessOCRResult(ctx, ocrResult)
}

// Recognize 执行OCR识别
func (s *Service) Recognize(base64Image string) (*ocr.OCRResult, error) {
	return s.ocrClient.Recognize(base64Image)
}
