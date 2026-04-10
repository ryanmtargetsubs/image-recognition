package config

import "os"

type Config struct {
	Port           string
	MaxUploadSize  int64
	TesseractLang  string
	UploadDir      string
	AllowedOrigins string

	OpenAIKey   string
	OpenAIModel string
	AIEnabled   bool
}

func Load() *Config {
	apiKey := getEnv("OPENAI_API_KEY", "")
	return &Config{
		Port:           getEnv("PORT", "3000"),
		MaxUploadSize:  50 * 1024 * 1024, // 50 MB
		TesseractLang:  getEnv("TESSERACT_LANG", "eng"),
		UploadDir:      getEnv("UPLOAD_DIR", "./uploads"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "*"),

		OpenAIKey:   apiKey,
		OpenAIModel: getEnv("OPENAI_MODEL", "gpt-4o-mini"),
		AIEnabled:   apiKey != "",
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
