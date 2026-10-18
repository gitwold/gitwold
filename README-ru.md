# Gitwold

Лёгкий self-hosted Git-хостинг.

Стек: go-git, SQLite, GORM.

Большинство Git-хостингов группируют проекты через организации и команды.
Gitwold этого не делает: у каждого пользователя своё пространство, а репозитории лежат во вложенных папках - как в файловой системе.

```
user/
├── work/
│   ├── backend/api.git
│   └── frontend/web.git
└── playground/hello.git
```

## Запуск

```bash
go run ./cmd/gitwold

cd web && yarn install && yarn dev
```

HTTP API слушает `:8080`. SSH для git - `:2222` (`GITWOLD_SSH_ADDR`).

### Docker

```bash
docker compose up --build
```

Данные в volumes `gitwold-data` и `gitwold-backups`.

### Бэкап

```bash
./scripts/backup.sh
```

Создаёт `tar.gz` с SQLite, bare-репозиториями, аватарами и SSH host keys.

Справка по API: [API.md](API.md)

Локализация: [langs/](langs/) (`en` / `ru`, `api.json` + `web.json`)
