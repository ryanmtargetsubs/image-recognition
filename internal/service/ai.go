package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ryanmtargetsubs/image-recognition/internal/model"
)

const openAIURL = "https://api.openai.com/v1/chat/completions"

const systemPrompt = `You are a CV/resume data extraction expert. Given raw OCR text from a CV image, extract structured information and return ONLY valid JSON matching this exact schema:

{
  "name": "Full name of the candidate",
  "email": "Email address",
  "phone": "Phone number",
  "location": "City, Country or address",
  "linkedin": "LinkedIn URL if present",
  "website": "Personal website or portfolio URL if present",
  "skills": ["skill1", "skill2"],
  "languages": ["language1", "language2"],
  "education": [
    {
      "title": "Degree or certification name",
      "subtitle": "Institution name",
      "date_range": "Start - End",
      "description": "Additional details"
    }
  ],
  "experience": [
    {
      "title": "Job title",
      "subtitle": "Company name",
      "date_range": "Start - End",
      "description": "Key responsibilities and achievements"
    }
  ],
  "certificates": ["cert1", "cert2"],
  "summary": "Brief professional summary extracted from the CV",
  "ai_summary": "Your 2-3 sentence assessment of this candidate's profile, strengths, and seniority level"
}

Rules:
- Return ONLY the JSON object, no markdown fences, no explanation.
- Use empty string "" for missing text fields, empty array [] for missing list fields.
- Preserve the original language of the CV content.
- For ai_summary, provide your own professional assessment of the candidate.
- Normalize phone numbers and emails when possible.
- Merge duplicate sections intelligently (e.g. "Work History" and "Experience" are the same).`

// AIService handles LLM-based CV text analysis.
type AIService struct {
	apiKey string
	model  string
	client *http.Client
}

// NewAIService creates an AI service configured with an OpenAI API key.
func NewAIService(apiKey, model string) *AIService {
	return &AIService{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// aiCVResult mirrors the JSON we ask the LLM to produce.
type aiCVResult struct {
	Name         string         `json:"name"`
	Email        string         `json:"email"`
	Phone        string         `json:"phone"`
	Location     string         `json:"location"`
	LinkedIn     string         `json:"linkedin"`
	Website      string         `json:"website"`
	Skills       []string       `json:"skills"`
	Languages    []string       `json:"languages"`
	Education    []model.Section `json:"education"`
	Experience   []model.Section `json:"experience"`
	Certificates []string       `json:"certificates"`
	Summary      string         `json:"summary"`
	AISummary    string         `json:"ai_summary"`
}

// AnalyzeCV sends the raw OCR text to the OpenAI API and returns structured CV data.
func (s *AIService) AnalyzeCV(ctx context.Context, rawText string) (*model.CVData, error) {
	reqBody := chatRequest{
		Model:       s.model,
		Temperature: 0.1,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: fmt.Sprintf("Extract structured data from this CV text:\n\n%s", rawText)},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, openAIURL, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("openai request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var chatResp chatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if chatResp.Error != nil {
		return nil, fmt.Errorf("openai error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("openai returned no choices")
	}

	content := chatResp.Choices[0].Message.Content

	var result aiCVResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("parse AI output: %w (raw: %.500s)", err, content)
	}

	cv := &model.CVData{
		RawText:      rawText,
		Name:         result.Name,
		Email:        result.Email,
		Phone:        result.Phone,
		Location:     result.Location,
		LinkedIn:     result.LinkedIn,
		Website:      result.Website,
		Skills:       result.Skills,
		Languages:    result.Languages,
		Education:    result.Education,
		Experience:   result.Experience,
		Certificates: result.Certificates,
		Summary:      result.Summary,
		AISummary:    result.AISummary,
		ProcessedBy:  "ai",
	}

	return cv, nil
}
