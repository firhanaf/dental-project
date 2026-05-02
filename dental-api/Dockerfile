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

# su-exec: lightweight tool untuk switch user di entrypoint
RUN apk add --no-cache ca-certificates tzdata su-exec

# Non-root user (UID/GID default Alpine: 1000/1000 via -S)
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Direktori upload — ownership di-fix ulang oleh entrypoint saat runtime
RUN mkdir -p /uploads

WORKDIR /app
COPY --from=builder /app/dental-api .
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

# Tidak pakai USER appuser di sini — entrypoint.sh yang switch ke appuser
# setelah chown /uploads
EXPOSE 8080
ENTRYPOINT ["./entrypoint.sh"]
