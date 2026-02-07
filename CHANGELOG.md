## v0.1.1

- CLI/Core: support --desc/--description end-to-end (domain/ports/usecase/fs/http)
- Docs: improved quickstart + example responses
- Release: create/patch release by tag (avoid untagged drafts), script supports -Tag
- CI: add go vet (optional race via RUN_RACE env)
# Changelog

## v0.1.0
- Kernel: Units & Versions with validation
- Ports & DTO contracts
- In-memory + FS JSON repository
- HTTP API with consistent error format
- CLI (cmd/digiemu) + HTTP API cmd (cmd/api)
- CI workflow (gofmt check + go test)

