#!/usr/bin/env bash
# Migration runner dengan tracking — hanya apply file SQL yang belum pernah dijalankan
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_DIR"

# Load .env
set -a
# shellcheck disable=SC1091
source .env
set +a

DB_USER="${DB_USER:-dental}"
DB_NAME="${DB_NAME:-dentaldb}"

# Helper: jalankan SQL via stdin ke postgres container
psql_exec() {
  docker compose exec -T postgres psql -U "$DB_USER" -d "$DB_NAME" "$@"
}

# Buat tabel tracker jika belum ada
psql_exec -c "
  CREATE TABLE IF NOT EXISTS schema_migrations (
    filename   VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
  );
" >/dev/null

echo "=== Checking migrations ==="

applied=0
skipped=0

for f in migrations/*.sql; do
  [ -f "$f" ] || continue
  filename="$(basename "$f")"

  # Cek apakah sudah pernah diapply
  count=$(psql_exec -t -A -c \
    "SELECT COUNT(*) FROM schema_migrations WHERE filename = '$filename';" \
    2>/dev/null | tr -d '[:space:]')

  if [ "$count" = "0" ]; then
    echo "  → Applying : $filename"

    # Jalankan SQL file via stdin
    if psql_exec < "$f"; then
      psql_exec -c \
        "INSERT INTO schema_migrations (filename) VALUES ('$filename');" \
        >/dev/null
      echo "  ✓ Applied  : $filename"
      applied=$((applied + 1))
    else
      echo "  ✗ FAILED   : $filename — deployment dibatalkan"
      exit 1
    fi
  else
    echo "  · Skip     : $filename (already applied)"
    skipped=$((skipped + 1))
  fi
done

echo "=== Migrations: $applied applied, $skipped skipped ==="
