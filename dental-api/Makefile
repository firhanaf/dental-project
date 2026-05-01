.PHONY: run build migrate test lint

run:
	go run ./cmd/server

build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/dental-api ./cmd/server

# Jalankan semua migration secara berurutan
migrate:
	@echo "Running migrations..."
	@for f in migrations/*.sql; do \
		echo "  → $$f"; \
		psql "$$DATABASE_URL" -f $$f; \
	done
	@echo "Migration done."

test:
	go test ./... -v

lint:
	golangci-lint run ./...

# Docker helpers
up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f go-api

rebuild:
	docker compose build go-api && docker compose up -d go-api
