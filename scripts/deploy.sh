#!/usr/bin/env bash
# Deploy script — dijalankan di VPS oleh GitHub Actions
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_DIR"

log() { echo "[$(date '+%H:%M:%S')] $*"; }

# Guard: pastikan .env ada
if [ ! -f .env ]; then
  echo "ERROR: .env tidak ditemukan di $PROJECT_DIR"
  echo "Jalankan: cp .env.example .env && nano .env"
  exit 1
fi

log "=== [1/5] Git pull ==="
git pull origin main

log "=== [2/5] Build frontend ==="
docker run --rm \
  -v "$PROJECT_DIR/dental-web":/app \
  -w /app \
  node:20-alpine \
  sh -c "npm ci --prefer-offline && npm run build"

log "=== [3/5] Start postgres & tunggu ready ==="
docker compose up -d postgres

# Baca DB_USER dan DB_NAME dari .env
DB_USER=$(grep -E '^DB_USER=' .env | cut -d= -f2)
DB_NAME=$(grep -E '^DB_NAME=' .env | cut -d= -f2)

for i in $(seq 1 30); do
  if docker compose exec -T postgres pg_isready -U "$DB_USER" -d "$DB_NAME" >/dev/null 2>&1; then
    log "Postgres ready."
    break
  fi
  [ "$i" -eq 30 ] && { log "ERROR: Postgres tidak kunjung ready."; exit 1; }
  sleep 2
done

log "=== [4/5] Run migrations ==="
bash "$PROJECT_DIR/scripts/migrate.sh"

log "=== [5/5] Rebuild & restart semua service ==="
docker compose build go-api
docker compose up -d

# Beri waktu nginx naik, lalu reload config
sleep 3
docker compose exec -T nginx nginx -s reload 2>/dev/null || true

log "=== Deploy selesai: $(date) ==="
