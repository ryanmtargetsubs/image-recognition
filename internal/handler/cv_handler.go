package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ryanmtargetsubs/image-recognition/internal/model"
	"github.com/ryanmtargetsubs/image-recognition/internal/service"
)

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".bmp":  true,
	".tiff": true,
	".tif":  true,
	".webp": true,
}

// CVHandler holds HTTP handlers for CV processing endpoints.
type CVHandler struct {
	cvService *service.CVService
	uploadDir string
}

// NewCVHandler creates a new handler backed by the given service.
func NewCVHandler(cvService *service.CVService, uploadDir string) *CVHandler {
	return &CVHandler{cvService: cvService, uploadDir: uploadDir}
}

// UploadAndProcess handles POST /api/v1/cv/upload
func (h *CVHandler) UploadAndProcess(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error: "image field is required (multipart form)",
		})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error: fmt.Sprintf("unsupported file type %q; allowed: jpg, jpeg, png, bmp, tiff, webp", ext),
		})
	}

	if err := os.MkdirAll(h.uploadDir, 0o750); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error: "failed to prepare upload directory",
		})
	}

	// Save with a sanitized name under the upload dir.
	safeName := filepath.Base(file.Filename)
	dst := filepath.Join(h.uploadDir, safeName)
	if err := c.SaveFile(file, dst); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error: "failed to save uploaded file",
		})
	}
	defer os.Remove(dst)

	cv, err := h.cvService.ProcessImage(dst)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(model.ErrorResponse{
			Error: fmt.Sprintf("CV processing failed: %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.UploadResponse{
		Success: true,
		Message: "CV processed successfully",
		Data:    cv,
	})
}

// Health handles GET /api/v1/health
func (h *CVHandler) Health(c *fiber.Ctx) error {
	return c.JSON(model.HealthResponse{
		Status:  "ok",
		Service: "image-recognition",
	})
}
