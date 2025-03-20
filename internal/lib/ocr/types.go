package ocr

// OCRResult represents a common OCR result structure used across different OCR implementations
type OCRResult struct {
	Code int          `json:"code"`
	Data []OCRTextBox `json:"data"`
}

// OCRTextBox represents a single text detection result
type OCRTextBox struct {
	Text       string   `json:"text"`
	Confidence float64  `json:"confidence"`
	Box        [][2]int `json:"box,omitempty"`
}

// OCRRequest represents a common OCR request structure
type OCRRequest struct {
	ImageBase64 string                 `json:"image"`
	Options     map[string]interface{} `json:"options,omitempty"`
}
