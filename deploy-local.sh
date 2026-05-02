#!/usr/bin/env bash
# Jalankan dari laptop sebelum deploy ke EC2.
# Usage: ./deploy-local.sh <ec2-ip> [path-to-pem]
#
# Contoh:
#   ./deploy-local.sh 54.123.45.67 ~/.ssh/dental-key.pem
set -e

# Terima format "54.x.x.x" atau "ubuntu@54.x.x.x"
EC2_IP="${1:?Masukkan EC2 IP. Contoh: ./deploy-local.sh 54.123.45.67}"
PEM="${2:-~/.ssh/id_rsa}"
# Strip "ubuntu@" jika sudah ada, lalu tambah kembali
EC2_IP="${EC2_IP#ubuntu@}"
REMOTE="ubuntu@$EC2_IP"
REMOTE_DIR="~/dental-project"

# ── Build frontend ────────────────────────────────────────────
cd dental-web

if [ ! -d node_modules ]; then
  echo "▶ Install dependencies (pertama kali)..."
  node_modules/.bin/npm install || npm.cmd install || npm install
fi

echo "▶ Build React frontend..."
node node_modules/typescript/bin/tsc -b
node node_modules/vite/bin/vite.js build
cd ..

echo "  Build selesai → dental-web/dist/"

# ── Kirim dist ke EC2 ─────────────────────────────────────────
echo "▶ Upload dist ke EC2 ($EC2_IP)..."
ssh -i "$PEM" "$REMOTE" "mkdir -p $REMOTE_DIR/dental-web/dist"
scp -i "$PEM" -r dental-web/dist/. "$REMOTE:$REMOTE_DIR/dental-web/dist/"
echo "  Upload selesai."

# ── Pull kode terbaru di EC2 ──────────────────────────────────
echo "▶ Pull latest code di EC2..."
ssh -i "$PEM" "$REMOTE" "cd $REMOTE_DIR && git pull --ff-only"

# ── Jalankan deploy.sh di EC2 ─────────────────────────────────
echo "▶ Jalankan deploy.sh di EC2..."
ssh -i "$PEM" "$REMOTE" "cd $REMOTE_DIR && bash deploy.sh"
