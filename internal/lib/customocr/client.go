package customocr

import (
	"context"
	"fmt"
	"log"

	"github.com/straydragon/bookxnote-local-ocr/internal/client/openapi"
	"github.com/straydragon/bookxnote-local-ocr/internal/lib/ocr"
)

type Client struct {
	apiClient *openapi.APIClient
}

func NewClient(baseURL string, apiKey string) *Client {
	config := openapi.NewConfiguration()
	config.Servers = []openapi.ServerConfiguration{
		{
			URL: baseURL,
		},
	}
	config.AddDefaultHeader("Authorization", "Bearer "+apiKey)

	return &Client{
		apiClient: openapi.NewAPIClient(config),
	}
}

func (c *Client) Recognize(base64Image string) (*ocr.OCRResult, error) {
	ctx := context.Background()

	// Make the API call directly with the base64 string
	resp, _, err := c.apiClient.DefaultAPI.PostOcrByBxnLocalOcr(ctx).Base64Image(base64Image).Execute()
	if err != nil {
		return nil, fmt.Errorf("call OCR service failed: %w", err)
	}

	// Convert response to common format
	result := &ocr.OCRResult{
		Code: 0,
		Data: []ocr.OCRTextBox{
			{
				Text:       resp.Data.GetText(),
				Confidence: float64(resp.Data.GetConfidence()),
				Box:        [][2]int{{0, 0}, {0, 0}, {0, 0}, {0, 0}},
			},
		},
	}

	log.Printf("ocr result (raw: custom-ocr): %+v\n", result)

	return result, nil
}
