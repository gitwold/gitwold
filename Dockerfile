FROM golang:1.27.1-bookworm AS api

WORKDIR /src

ENV GOTOOLCHAIN=auto

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o /out/gitwold ./cmd/gitwold

FROM node:22-bookworm AS web

WORKDIR /web

COPY web/package.json web/yarn.lock ./

RUN yarn install --frozen-lockfile

COPY web/ ./

RUN yarn build

FROM debian:bookworm-slim

RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates git sqlite3 \
  && rm -rf /var/lib/apt/lists/* \

WORKDIR /app

COPY --from=api /out/gitwold /usr/local/bin/gitwold
COPY --from=web /web/dist /app/web/dist
COPY scripts/backup.sh /usr/local/bin/gitwold-backup

RUN chmod +x /usr/local/bin/gitwold-backup && mkdir -p /app/data /app/backups \

ENV GITWOLD_SSH_ADDR=:2222

EXPOSE 8080 2222

VOLUME ["/app/data", "/app/backups"]

CMD ["gitwold"]
