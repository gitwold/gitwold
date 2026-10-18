# Gitwold

Lightweight self-hosted Git hosting.

Built with go-git, SQLite, GORM.

Most Git hostings group projects with organizations and teams.
Gitwold skips that model: each user has a personal namespace, and repositories live in nested folders - like a filesystem.

```
user/
├── work/
│   ├── backend/api.git
│   └── frontend/web.git
└── playground/hello.git
```

## Run

```bash
go run ./cmd/gitwold

cd web && yarn install && yarn dev
```

HTTP API listens on `:8080`. SSH git transport listens on `:2222` (`GITWOLD_SSH_ADDR`).

### Docker

```bash
docker compose up --build
```

Data persists in volumes `gitwold-data` and `gitwold-backups`.

### Backup

```bash
./scripts/backup.sh
# or: GITWOLD_DATA_DIR=./data GITWOLD_BACKUP_DIR=./backups ./scripts/backup.sh
```

Writes a timestamped `tar.gz` with SQLite, bare repos, avatars, and SSH host keys.

API reference: [API.md](API.md)

Locales: [langs/](langs/) (`en` / `ru`, `api.json` + `web.json`)
