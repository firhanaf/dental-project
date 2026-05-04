#!/usr/bin/env bash
# Jalankan SEKALI di VPS yang sudah ada DB-nya.
# Membuat tabel schema_migrations dan menandai semua file SQL yang sudah ada
# sebagai "already applied" — sehingga migrate.sh tidak re-run dari awal.
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_DIR"

set -a
# shellcheck disable=SC1091
source .env
set +a

DB_USER="${DB_USER:-dental}"
DB_NAME="${DB_NAME:-dentaldb}"

psql_exec() {
  docker compose exec -T postgres psql -U "$DB_USER" -d "$DB_NAME" "$@"
}

echo "=== Membuat tabel schema_migrations ==="
psql_exec -c "
  CREATE TABLE IF NOT EXISTS schema_migrations (
    filename   VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
  );
"

echo "=== Menandai semua migration yang sudah ada sebagai applied ==="
for f in migrations/*.sql; do
  [ -f "$f" ] || continue
  filename="$(basename "$f")"

  psql_exec -c \
    "INSERT INTO schema_migrations (filename) VALUES ('$filename') ON CONFLICT DO NOTHING;" \
    >/dev/null

  echo "  ✓ Marked: $filename"
done

echo ""
echo "=== Bootstrap selesai ==="
echo "Sekarang migrate.sh akan skip semua file di atas."
echo "File SQL BARU yang kamu tambahkan ke depannya akan dijalankan otomatis."
