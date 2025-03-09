package service

import (
	"context"
	"fmt"
	"log"

	"github.com/straydragon/bookxnote-local-ocr/internal/lib/langchain"
	"github.com/straydragon/bookxnote-local-ocr/internal/lib/umiocr"
)

// OCR后处理管理
type PostOCRProcessor struct {
	langchainClient *langchain.Client
	config          *Config
}

// reloadConfig 重新加载配置
func (p *PostOCRProcessor) reloadConfig() {
	cfg, err := GetUserConfig()
	if err != nil {
		log.Printf("failed to get user config: %v", err)
	}
	p.config = cfg
}

func NewPostOCRProcessor() (*PostOCRProcessor, error) {
	config, err := GetUserConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get user config: %w", err)
	}
	if config.LLM == nil || len(config.LLM.Models) == 0 {
		return nil, fmt.Errorf("no LLM models configured")
	}

	var defaultModel *ConfigLLMModel
	for _, model := range config.LLM.Models {
		if model.Ident == config.AfterOCR.Translate.UseIdent {
			defaultModel = model
			break
		}
	}

	if defaultModel == nil {
		// Use the first model as default
		defaultModel = config.LLM.Models[0]
	}

	langchainClient, err := langchain.NewClient(langchain.ClientOptions{
		APIKey:     defaultModel.ApiKey,
		APIBaseURL: defaultModel.ApiBaseUrl,
		ModelName:  defaultModel.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create LangChain client: %w", err)
	}

	return &PostOCRProcessor{
		langchainClient: langchainClient,
		config:          config,
	}, nil
}

// ProcessOCRResult 处理OCR结果
func (p *PostOCRProcessor) ProcessOCRResult(ctx context.Context, ocrResult *umiocr.APIRecognizeResp) (*umiocr.APIRecognizeResp, error) {
	if ocrResult == nil || len(ocrResult.Data) == 0 {
		return ocrResult, nil
	}

	var fullText string
	for _, item := range ocrResult.Data {
		fullText += item.Text + "\n"
	}

	processedText, err := p.processText(ctx, fullText)
	if err != nil {
		return nil, fmt.Errorf("failed to process OCR text: %w", err)
	}

	result := &umiocr.APIRecognizeResp{
		Code: ocrResult.Code,
		Data: []umiocr.APIRecognizeRespData{
			{
				Text:  processedText,
				Score: 1.0,
				Box:   [][2]int{{0, 0}, {0, 0}, {0, 0}, {0, 0}},
				End:   "",
			},
		},
	}

	return result, nil
}

// processText 使用langchain处理文本
func (p *PostOCRProcessor) processText(ctx context.Context, text string) (string, error) {
	processedText := text
	p.reloadConfig()

	if p.config.AfterOCR.AutoFixContent != nil && p.config.AfterOCR.AutoFixContent.Enabled {
		log.Println("Applying auto fix content")
		var err error
		processedText, err = p.langchainClient.FixTextLayout(ctx, processedText)
		if err != nil {
			return "", fmt.Errorf("failed to fix text layout: %w", err)
		}
	}

	if p.config.AfterOCR.Translate != nil && p.config.AfterOCR.Translate.Enabled {
		log.Printf("Translating text to %s", p.config.AfterOCR.Translate.TargetLang)
		var err error
		processedText, err = p.langchainClient.TranslateText(ctx, processedText, p.config.AfterOCR.Translate.TargetLang)
		if err != nil {
			return "", fmt.Errorf("failed to translate text: %w", err)
		}
	}

	if p.config.AfterOCR.GenerateByLLM != nil && p.config.AfterOCR.GenerateByLLM.Enabled {
		log.Println("Generating notes")

		var promptTemplate string
		for _, prompt := range p.config.LLM.Prompts {
			if prompt.Ident == p.config.AfterOCR.GenerateByLLM.PromptIdent {
				promptTemplate = prompt.Prompt
				break
			}
		}

		if promptTemplate == "" {
			return "", fmt.Errorf("prompt template not found: %s", p.config.AfterOCR.GenerateByLLM.PromptIdent)
		}

		var err error
		processedText, err = p.langchainClient.GenerateNotes(ctx, processedText, promptTemplate)
		if err != nil {
			return "", fmt.Errorf("failed to generate notes: %w", err)
		}
	}

	return processedText, nil
}
