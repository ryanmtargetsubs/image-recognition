package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/ryanmtargetsubs/image-recognition/internal/config"
	"github.com/ryanmtargetsubs/image-recognition/internal/handler"
	"github.com/ryanmtargetsubs/image-recognition/internal/middleware"
	"github.com/ryanmtargetsubs/image-recognition/internal/router"
	"github.com/ryanmtargetsubs/image-recognition/internal/service"
)

func main() {
	cfg := config.Load()

	ocrService := service.NewOCRService(cfg.TesseractLang)
	parser := service.NewCVParser()

	var aiService *service.AIService
	if cfg.AIEnabled {
		aiService = service.NewAIService(cfg.OpenAIKey, cfg.OpenAIModel)
		log.Printf("AI processing enabled (model: %s)", cfg.OpenAIModel)
	} else {
		log.Println("AI processing disabled (set OPENAI_API_KEY to enable)")
	}

	cvService := service.NewCVService(ocrService, parser, aiService)
	cvHandler := handler.NewCVHandler(cvService, cfg.UploadDir)

	app := fiber.New(fiber.Config{
		BodyLimit: int(cfg.MaxUploadSize),
	})

	middleware.Setup(app, cfg.AllowedOrigins)
	router.Setup(app, cvHandler)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting image-recognition server on %s", addr)
	log.Fatal(app.Listen(addr))
}
