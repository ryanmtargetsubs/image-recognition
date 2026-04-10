# Build stage
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache gcc musl-dev tesseract-ocr-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -o /app/server ./cmd/server

# Runtime stage
FROM alpine:3.20

RUN apk add --no-cache tesseract-ocr tesseract-ocr-data-eng

COPY --from=builder /app/server /usr/local/bin/server

RUN adduser -D appuser
USER appuser

EXPOSE 3000

ENTRYPOINT ["server"]
