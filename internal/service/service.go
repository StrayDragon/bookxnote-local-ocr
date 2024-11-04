package service

import (
	"github.com/straydragon/bookxnote-local-ocr/internal/lib/umiocr"
)

// Service 包含所有服务依赖
type Service struct {
	UmiOCRClient *umiocr.Client
}

// NewService 创建新的服务实例
func NewService(ocrClient *umiocr.Client) *Service {
	return &Service{
		UmiOCRClient: ocrClient,
	}
}
