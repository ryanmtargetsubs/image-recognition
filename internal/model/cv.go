package model

// CVData represents structured data extracted from a CV/resume image.
type CVData struct {
	RawText      string    `json:"raw_text"`
	Name         string    `json:"name,omitempty"`
	Email        string    `json:"email,omitempty"`
	Phone        string    `json:"phone,omitempty"`
	Location     string    `json:"location,omitempty"`
	LinkedIn     string    `json:"linkedin,omitempty"`
	Website      string    `json:"website,omitempty"`
	Skills       []string  `json:"skills,omitempty"`
	Languages    []string  `json:"languages,omitempty"`
	Education    []Section `json:"education,omitempty"`
	Experience   []Section `json:"experience,omitempty"`
	Certificates []string  `json:"certificates,omitempty"`
	Summary      string    `json:"summary,omitempty"`
	AISummary    string    `json:"ai_summary,omitempty"`
	ProcessedBy  string    `json:"processed_by"`
}

// Section represents a generic CV section entry.
type Section struct {
	Title       string `json:"title,omitempty"`
	Subtitle    string `json:"subtitle,omitempty"`
	DateRange   string `json:"date_range,omitempty"`
	Description string `json:"description,omitempty"`
}

// UploadResponse is the API response after processing a CV image.
type UploadResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Data    *CVData `json:"data,omitempty"`
}

// ErrorResponse is a standard error payload.
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// HealthResponse is the health-check payload.
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}
