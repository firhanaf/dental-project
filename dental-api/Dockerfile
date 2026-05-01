# ── Stage 1: builder ─────────────────────────────────────────
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Download deps dulu (cache layer)
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o dental-api ./cmd/server

# ── Stage 2: runtime ─────────────────────────────────────────
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

# Non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Direktori upload (override oleh volume di compose)
RUN mkdir -p /uploads && chown appuser:appgroup /uploads

WORKDIR /app
COPY --from=builder /app/dental-api .
RUN chown appuser:appgroup dental-api

USER appuser
EXPOSE 8080
CMD ["./dental-api"]
