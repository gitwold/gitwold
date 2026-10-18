# Changelog

## 0.1.0 - dev

Self-hosted Git hosting: personal folders and organizations share the `/{owner}/...` namespace.

- Accounts: register, login, logout, sessions
- Profile: email, avatar upload and delete
- Access tokens: create and revoke in the API and UI (`gitwold_` prefix)
- SSH keys for Git over SSH
- Nested folders in a personal namespace and inside an organization
- Repositories: create, update, move, delete, public and private, default branch
- Transfer a personal repository into an organization
- Organizations: create, update, delete, members and roles, repository ACL
- Git over HTTP and SSH (clone, fetch, push)
- Browse: tree, blob, commits, diff, readme, languages, blame, raw download
- Issues: open, close, comments, labels
- Pull requests: open, close, comments, diff, merge
- Search repositories by name
- Activity feed for push, issues, pull requests, membership, and transfers
- YAML config with environment overrides (`GITWOLD_*`)
- UI in English and Russian
- Rate limiting, backup script, Docker image and Compose
