package service

import (
	"context"
	"log"

	"github.com/ryanmtargetsubs/image-recognition/internal/model"
)

// CVService orchestrates OCR extraction and CV parsing.
type CVService struct {
	ocr    *OCRService
	parser *CVParser
	ai     *AIService
}

// NewCVService wires up the CV processing pipeline.
// Pass ai as nil to disable AI processing and use regex-only parsing.
func NewCVService(ocr *OCRService, parser *CVParser, ai *AIService) *CVService {
	return &CVService{ocr: ocr, parser: parser, ai: ai}
}

// ProcessImage runs OCR on the image file, then uses AI (with regex fallback) to parse the CV.
func (s *CVService) ProcessImage(imagePath string) (*model.CVData, error) {
	rawText, err := s.ocr.ExtractText(imagePath)
	if err != nil {
		return nil, err
	}

	// Try AI-powered extraction first.
	if s.ai != nil {
		cv, aiErr := s.ai.AnalyzeCV(context.Background(), rawText)
		if aiErr == nil {
			return cv, nil
		}
		log.Printf("AI processing failed, falling back to regex parser: %v", aiErr)
	}

	// Fallback: regex-based parser.
	cv := s.parser.Parse(rawText)
	cv.ProcessedBy = "regex"
	return cv, nil
}
