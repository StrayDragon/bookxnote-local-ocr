package langchain

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type Client struct {
	llm *openai.LLM
}

type ClientOptions struct {
	APIKey     string
	APIBaseURL string
	ModelName  string
}

func NewClient(opts ClientOptions) (*Client, error) {
	llmOpts := []openai.Option{}

	if opts.ModelName != "" {
		llmOpts = append(llmOpts, openai.WithModel(opts.ModelName))
	}

	if opts.APIBaseURL != "" {
		llmOpts = append(llmOpts, openai.WithBaseURL(opts.APIBaseURL))
	}

	if opts.APIKey != "" {
		llmOpts = append(llmOpts, openai.WithToken(opts.APIKey))
	}

	llm, err := openai.New(llmOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI client: %w", err)
	}

	return &Client{
		llm: llm,
	}, nil
}

// GenerateText 使用LLM生成文本
func (c *Client) GenerateText(ctx context.Context, prompt string) (string, error) {
	resp, err := llms.GenerateFromSinglePrompt(ctx, c.llm, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate text: %w", err)
	}

	log.Printf("LLM response: %s", resp)
	return resp, nil
}

// TranslateText 翻译文本
func (c *Client) TranslateText(ctx context.Context, text, targetLang string) (string, error) {
	prompt := fmt.Sprintf("Translate the following text to %s:\n\n%s", targetLang, text)
	return c.GenerateText(ctx, prompt)
}

// GenerateNotes 根据OCR文本生成笔记
func (c *Client) GenerateNotes(ctx context.Context, text, promptTemplate string) (string, error) {
	prompt := fmt.Sprintf(promptTemplate, text)
	return c.GenerateText(ctx, prompt)
}

// FixTextLayout 修复OCR文本的布局
func (c *Client) FixTextLayout(ctx context.Context, text string) (string, error) {
	prompt := fmt.Sprintf("Fix the layout of the following OCR text, preserving paragraphs and removing unnecessary line breaks:\n\n%s", text)
	return c.GenerateText(ctx, prompt)
}
