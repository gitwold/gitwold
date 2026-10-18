#!/usr/bin/env bash

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
DATA_DIR="${GITWOLD_DATA_DIR:-$ROOT/data}"
OUT_DIR="${GITWOLD_BACKUP_DIR:-$ROOT/backups}"
STAMP="$(date -u +%Y%m%dT%H%M%SZ)"
WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

mkdir -p "$OUT_DIR"

if [[ ! -d "$DATA_DIR" ]]; then
  echo "data dir not found: $DATA_DIR" >&2
  exit 1
fi

mkdir -p "$WORK/data"

if [[ -f "$DATA_DIR/db.sqlite" ]]; then
  if command -v sqlite3 >/dev/null 2>&1; then
    sqlite3 "$DATA_DIR/db.sqlite" ".backup '$WORK/data/db.sqlite'"
  else
    cp "$DATA_DIR/db.sqlite" "$WORK/data/db.sqlite"
    [[ -f "$DATA_DIR/db.sqlite-wal" ]] && cp "$DATA_DIR/db.sqlite-wal" "$WORK/data/db.sqlite-wal"
    [[ -f "$DATA_DIR/db.sqlite-shm" ]] && cp "$DATA_DIR/db.sqlite-shm" "$WORK/data/db.sqlite-shm"
  fi
fi

if [[ -d "$DATA_DIR/repos" ]]; then
  mkdir -p "$WORK/data/repos"
  cp -a "$DATA_DIR/repos/." "$WORK/data/repos/"
fi

if [[ -d "$DATA_DIR/avatars" ]]; then
  mkdir -p "$WORK/data/avatars"
  cp -a "$DATA_DIR/avatars/." "$WORK/data/avatars/"
fi

if [[ -d "$DATA_DIR/ssh" ]]; then
  mkdir -p "$WORK/data/ssh"
  cp -a "$DATA_DIR/ssh/." "$WORK/data/ssh/"
fi

ARCHIVE="$OUT_DIR/gitwold-backup-$STAMP.tar.gz"
tar -C "$WORK" -czf "$ARCHIVE" data
echo "$ARCHIVE"
