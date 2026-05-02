#!/bin/sh
# Fix ownership of mounted volume (runs as root before switching to appuser)
chown -R appuser:appgroup /uploads
exec su-exec appuser ./dental-api
