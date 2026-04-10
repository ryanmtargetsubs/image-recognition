package service

import (
	"fmt"

	"github.com/otiai10/gosseract/v2"
)

// OCRService handles text extraction from images using Tesseract.
type OCRService struct {
	lang string
}

// NewOCRService creates an OCR service with the specified language.
func NewOCRService(lang string) *OCRService {
	return &OCRService{lang: lang}
}

// ExtractText reads text content from the image at the given file path.
func (s *OCRService) ExtractText(imagePath string) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	if err := client.SetLanguage(s.lang); err != nil {
		return "", fmt.Errorf("set language %q: %w", s.lang, err)
	}

	if err := client.SetImage(imagePath); err != nil {
		return "", fmt.Errorf("set image %q: %w", imagePath, err)
	}

	text, err := client.Text()
	if err != nil {
		return "", fmt.Errorf("extract text: %w", err)
	}

	return text, nil
}
