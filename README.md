# image-recognition

Go backend service (Fiber v2) that reads CV/resume images, runs OCR via Tesseract, and uses **AI (OpenAI)** to intelligently extract structured data — handling any CV format.

When `OPENAI_API_KEY` is set, the service sends OCR text to an LLM for smart extraction. If AI is unavailable or fails, it falls back to regex-based parsing automatically.

## Prerequisites

- **Go 1.22+**
- **Tesseract OCR** installed on the host

```bash
# Ubuntu / Debian
sudo apt-get install -y tesseract-ocr libtesseract-dev

# macOS
brew install tesseract
```

## Getting Started

```bash
# Install Go dependencies
go mod tidy

# Run the server
go run ./cmd/server
```

The server starts on port **3000** by default (override with `PORT` env var).

## API Endpoints

### Health Check

```
GET /api/v1/health
```

Response:
```json
{ "status": "ok", "service": "image-recognition" }
```

### Upload & Process CV

```
POST /api/v1/cv/upload
Content-Type: multipart/form-data
```

| Field   | Type | Description                          |
|---------|------|--------------------------------------|
| `image` | file | CV/resume image (jpg, png, bmp, tiff, webp) |

Response:
```json
{
  "success": true,
  "message": "CV processed successfully",
  "data": {
    "raw_text": "...",
    "name": "Jane Doe",
    "email": "jane@example.com",
    "phone": "+1 555-123-4567",
    "location": "San Francisco, CA",
    "linkedin": "https://linkedin.com/in/janedoe",
    "website": "https://janedoe.dev",
    "skills": ["Go", "Python", "Docker", "Kubernetes"],
    "languages": ["English", "Spanish"],
    "education": [{ "title": "BSc Computer Science", "subtitle": "MIT", "date_range": "2016-2020" }],
    "experience": [{ "title": "Software Engineer", "subtitle": "Acme Corp", "date_range": "2020-Present", "description": "Led backend team..." }],
    "certificates": ["AWS Solutions Architect"],
    "summary": "Experienced backend developer...",
    "ai_summary": "Senior backend engineer with 4+ years of experience, strong in cloud-native Go development.",
    "processed_by": "ai"
  }
}
```

### cURL Example

```bash
curl -X POST http://localhost:3000/api/v1/cv/upload \
  -F "image=@/path/to/cv-scan.png"
```

## Environment Variables

| Variable          | Default        | Description                     |
|-------------------|----------------|---------------------------------|
| `PORT`            | `3000`         | Server listen port              |
| `OPENAI_API_KEY`  | *(empty)*      | OpenAI API key (enables AI mode)|
| `OPENAI_MODEL`    | `gpt-4o-mini`  | OpenAI model to use             |
| `TESSERACT_LANG`  | `eng`          | Tesseract language pack         |
| `UPLOAD_DIR`      | `./uploads`    | Temp directory for uploads      |
| `ALLOWED_ORIGINS` | `*`            | CORS allowed origins            |

## Docker

```bash
docker build -t image-recognition .
docker run -p 3000:3000 -e OPENAI_API_KEY=sk-... image-recognition
```

## Project Structure

```
cmd/server/main.go              Entry point
internal/
  config/config.go              Environment-based config
  handler/cv_handler.go         Fiber HTTP handlers
  middleware/middleware.go       CORS, logger, recover
  model/cv.go                   Request/response models
  router/router.go              Route registration
  service/
    ai.go                       OpenAI LLM CV analysis
    ocr.go                      Tesseract OCR wrapper
    parser.go                   Regex fallback parser
    cv_service.go               Orchestrator (AI → regex fallback)
Dockerfile                      Multi-stage build
```