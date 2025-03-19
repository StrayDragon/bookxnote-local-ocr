package umiocr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/straydragon/bookxnote-local-ocr/internal/lib/ocr"
)

type Client struct {
	BaseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
	}
}

type apiRecognizeReq struct {
	Base64  string                 `json:"base64"`
	Options map[string]interface{} `json:"options"`
}

type apiRecognizeRespData struct {
	Text  string   `json:"text"`
	Score float64  `json:"score"`
	Box   [][2]int `json:"box"`
	End   string   `json:"end"`
}

type apiRecognizeResp struct {
	Code int                    `json:"code"`
	Data []apiRecognizeRespData `json:"data"`
}

func (c *Client) Recognize(base64Image string) (*ocr.OCRResult, error) {
	req := apiRecognizeReq{
		Base64: base64Image,
		Options: map[string]interface{}{
			"data.format": "dict",
		},
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	resp, err := http.Post(c.BaseURL+"/api/ocr",
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, fmt.Errorf("call OCR service failed: %w", err)
	}
	defer resp.Body.Close()

	var result apiRecognizeResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	log.Printf("ocr result (raw: umi-ocr): %+v\n", result)

	// Convert to common OCR result type
	commonResult := &ocr.OCRResult{
		Code: result.Code,
		Data: make([]ocr.OCRTextBox, len(result.Data)),
	}

	for i, item := range result.Data {
		commonResult.Data[i] = ocr.OCRTextBox{
			Text:       item.Text,
			Confidence: item.Score,
			Box:        item.Box,
		}
	}

	return commonResult, nil
}
