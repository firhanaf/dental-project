#!/usr/bin/env bash
# Jalankan di EC2 setelah deploy-local.sh selesai dari laptop.
# Tidak butuh Node.js di server.
set -e

# ── Cek .env ─────────────────────────────────────────────────
if [ ! -f .env ]; then
  echo "ERROR: .env tidak ditemukan."
  echo "       Jalankan: cp .env.example .env  lalu edit nilainya."
  exit 1
fi

# ── Cek React build sudah ada ────────────────────────────────
if [ ! -d dental-web/dist ] || [ -z "$(ls -A dental-web/dist)" ]; then
  echo "ERROR: dental-web/dist kosong."
  echo "       Jalankan deploy-local.sh dari laptop terlebih dahulu."
  exit 1
fi

echo "▶ Build Go API image..."
docker compose build go-api

echo "▶ Start semua container..."
docker compose up -d

echo "▶ Tunggu database siap..."
docker compose exec -T postgres sh -c \
  'until pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB" 2>/dev/null; do sleep 1; done'
echo "  Database ready."

echo "▶ Cek migration..."
DB_USER=$(grep ^DB_USER .env | cut -d= -f2 | tr -d ' ')
DB_NAME=$(grep ^DB_NAME .env | cut -d= -f2 | tr -d ' ')

TABLE_EXISTS=$(docker compose exec -T postgres psql \
  -U "$DB_USER" -d "$DB_NAME" \
  -tAc "SELECT to_regclass('public.users')")

if [ "$TABLE_EXISTS" = "" ]; then
  echo "▶ Jalankan migration..."
  for f in migrations/*.sql; do
    echo "  → $f"
    docker compose exec -T postgres psql -U "$DB_USER" -d "$DB_NAME" < "$f"
  done
  echo "  Migration selesai."
else
  echo "  Migration sudah ada, skip."
fi

echo ""
echo "✓ Deploy selesai!"
PUBLIC_IP=$(curl -sf --max-time 3 http://checkip.amazonaws.com 2>/dev/null \
  || hostname -I | awk '{print $1}')
echo "  URL  : http://$PUBLIC_IP"
echo "  Log  : docker compose logs -f"
echo "  Stop : docker compose down"
