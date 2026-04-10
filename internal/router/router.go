package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ryanmtargetsubs/image-recognition/internal/handler"
)

// Setup registers all routes on the Fiber app.
func Setup(app *fiber.App, cvHandler *handler.CVHandler) {
	api := app.Group("/api/v1")

	api.Get("/health", cvHandler.Health)
	api.Post("/cv/upload", cvHandler.UploadAndProcess)
}
