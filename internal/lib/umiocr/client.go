package umiocr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Client struct {
	BaseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
	}
}

type APIRecognizeReq struct {
	Base64  string                 `json:"base64"`
	Options map[string]interface{} `json:"options"`
}

type APIRecognizeRespData struct {
	Text  string   `json:"text"`
	Score float64  `json:"score"`
	Box   [][2]int `json:"box"`
	End   string   `json:"end"`
}

type APIRecognizeResp struct {
	Code int                    `json:"code"`
	Data []APIRecognizeRespData `json:"data"`
}

func (c *Client) Recognize(base64Image string) (*APIRecognizeResp, error) {
	req := APIRecognizeReq{
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

	var result APIRecognizeResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	log.Printf("ocr result (raw: umi-ocr): %+v\n", result)

	return &result, nil
}
