package service

import (
	"regexp"
	"strings"

	"github.com/ryanmtargetsubs/image-recognition/internal/model"
)

var (
	emailRe = regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
	phoneRe = regexp.MustCompile(`[\+]?[(]?[0-9]{1,4}[)]?[-\s./0-9]{7,15}`)
)

// sectionHeaders maps normalised header keywords to CV section types.
var sectionHeaders = map[string]string{
	"education":            "education",
	"academic":             "education",
	"qualification":        "education",
	"experience":           "experience",
	"employment":           "experience",
	"work history":         "experience",
	"professional history": "experience",
	"skill":               "skills",
	"competenc":           "skills",
	"technical":           "skills",
	"summary":             "summary",
	"objective":           "summary",
	"profile":             "summary",
	"about me":            "summary",
}

// CVParser extracts structured data from raw CV text.
type CVParser struct{}

// NewCVParser creates a new CVParser.
func NewCVParser() *CVParser {
	return &CVParser{}
}

// Parse takes raw OCR text and returns structured CV data.
func (p *CVParser) Parse(rawText string) *model.CVData {
	cv := &model.CVData{RawText: rawText}

	cv.Email = extractFirst(emailRe, rawText)
	cv.Phone = extractFirst(phoneRe, rawText)
	cv.Name = guessName(rawText)

	sections := splitSections(rawText)

	for sectionType, body := range sections {
		switch sectionType {
		case "skills":
			cv.Skills = parseSkills(body)
		case "education":
			cv.Education = parseSectionEntries(body)
		case "experience":
			cv.Experience = parseSectionEntries(body)
		case "summary":
			cv.Summary = strings.TrimSpace(body)
		}
	}

	return cv
}

func extractFirst(re *regexp.Regexp, text string) string {
	match := re.FindString(text)
	return strings.TrimSpace(match)
}

// guessName takes the first non-empty line that isn't an email/phone as the name.
func guessName(text string) string {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if emailRe.MatchString(line) || phoneRe.MatchString(line) {
			continue
		}
		lower := strings.ToLower(line)
		isHeader := false
		for kw := range sectionHeaders {
			if strings.Contains(lower, kw) {
				isHeader = true
				break
			}
		}
		if isHeader {
			continue
		}
		if len(line) > 60 {
			continue
		}
		return line
	}
	return ""
}

// splitSections splits text by recognised section headers.
func splitSections(text string) map[string]string {
	lines := strings.Split(text, "\n")
	result := make(map[string]string)
	currentSection := ""
	var buf strings.Builder

	for _, line := range lines {
		lower := strings.ToLower(strings.TrimSpace(line))
		matched := ""
		for kw, secType := range sectionHeaders {
			if strings.Contains(lower, kw) {
				matched = secType
				break
			}
		}
		if matched != "" {
			if currentSection != "" {
				result[currentSection] = buf.String()
			}
			currentSection = matched
			buf.Reset()
			continue
		}
		if currentSection != "" {
			buf.WriteString(line)
			buf.WriteString("\n")
		}
	}
	if currentSection != "" {
		result[currentSection] = buf.String()
	}
	return result
}

func parseSkills(text string) []string {
	var skills []string
	// Try comma-separated, bullet, or newline-separated.
	delimiters := regexp.MustCompile(`[,;|\n•●▪\-]`)
	parts := delimiters.Split(text, -1)
	for _, part := range parts {
		s := strings.TrimSpace(part)
		if s != "" && len(s) < 60 {
			skills = append(skills, s)
		}
	}
	return skills
}

func parseSectionEntries(text string) []model.Section {
	var entries []model.Section
	// Split on double newlines or lines that look like titles (short, no period).
	blocks := strings.Split(text, "\n\n")
	for _, block := range blocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}
		lines := strings.SplitN(block, "\n", 3)
		entry := model.Section{}
		if len(lines) >= 1 {
			entry.Title = strings.TrimSpace(lines[0])
		}
		if len(lines) >= 2 {
			entry.Subtitle = strings.TrimSpace(lines[1])
		}
		if len(lines) >= 3 {
			entry.Description = strings.TrimSpace(lines[2])
		}
		// Try to pull out a date range from title or subtitle.
		dateRe := regexp.MustCompile(`\d{4}\s*[-–]\s*(\d{4}|[Pp]resent|[Cc]urrent)`)
		for _, l := range lines {
			if m := dateRe.FindString(l); m != "" {
				entry.DateRange = m
				break
			}
		}
		entries = append(entries, entry)
	}
	return entries
}
