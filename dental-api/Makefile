.PHONY: run build migrate test lint up down logs rebuild dev-up dev-down dev-logs dev-migrate

# ── Development (native) ─────────────────────────────────────
run:
	go run ./cmd/server

build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/dental-api ./cmd/server

test:
	go test ./... -v

lint:
	golangci-lint run ./...

# ── Migration ────────────────────────────────────────────────
# Set DATABASE_URL sebelum menjalankan, contoh:
#   export DATABASE_URL="postgresql://dental:pass@localhost:5432/dentaldb?sslmode=disable"
#   make migrate
migrate:
	@echo "Running migrations..."
	@for f in ../migrations/*.sql; do \
		echo "  → $$f"; \
		psql "$$DATABASE_URL" -f $$f; \
	done
	@echo "Migration done."

# ── Local Dev (postgres only, Go jalan native) ───────────────
dev-up:
	docker compose -f ../docker-compose.dev.yml up -d
	@echo ""
	@echo "✓ Postgres siap di localhost:5432"
	@echo "✓ Adminer  siap di http://localhost:8081"
	@echo ""
	@echo "Langkah selanjutnya:"
	@echo "  1. cp .env.example .env  (jika belum ada)"
	@echo "  2. Pastikan DB_HOST=localhost di .env"
	@echo "  3. make dev-migrate  (jalankan migration)"
	@echo "  4. make run"

dev-down:
	docker compose -f ../docker-compose.dev.yml down

dev-logs:
	docker compose -f ../docker-compose.dev.yml logs -f postgres

dev-migrate:
	@echo "Running migrations (local dev)..."
	@export DB_HOST=$$(grep '^DB_HOST' .env | cut -d= -f2); \
	 export DB_PORT=$$(grep '^DB_PORT' .env | cut -d= -f2); \
	 export DB_NAME=$$(grep '^DB_NAME' .env | cut -d= -f2); \
	 export DB_USER=$$(grep '^DB_USER' .env | cut -d= -f2); \
	 export DB_PASSWORD=$$(grep '^DB_PASSWORD' .env | cut -d= -f2); \
	 for f in ../migrations/*.sql; do \
		echo "  → $$f"; \
		PGPASSWORD=$$DB_PASSWORD psql -h $$DB_HOST -p $$DB_PORT -U $$DB_USER -d $$DB_NAME -f $$f; \
	 done
	@echo "Migration done."

# ── Production Docker (full stack) ───────────────────────────
up:
	docker compose -f ../docker-compose.yml up -d

down:
	docker compose -f ../docker-compose.yml down

logs:
	docker compose -f ../docker-compose.yml logs -f go-api

rebuild:
	docker compose -f ../docker-compose.yml build go-api
	docker compose -f ../docker-compose.yml up -d go-api
