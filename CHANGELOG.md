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
 
## v0.2.9 — Audit & Snapshot Stabilization Release

- Deterministic Snapshot Hashing: `snapshotHash` included in export; `auditHash` when `--audit`.
- Hardened audit + snapshot pipeline across kernel, FS and memory adapters.
- Added export CLI command and verified audit verification tooling.

## v0.3.0 — Meaning Layer v1

- Added Meaning Layer v1: optional, version-scoped structured context attached to Units/Versions.
- Persistence: version-scoped sidecar files: `data/units/<unit-id>.<version-id>.meaning.json` and the `meaning_hash` stored on the version record.
- New audit event: `MEANING_SET` appended when a meaning is set. Payload includes `meaning_hash` and `meaning_path`.
- CLI: `digiemu meaning set` and `digiemu meaning show` thin adapters to the `SetMeaning` usecase.
- HTTP: `PUT /v1/units/{unit}/meaning` and `GET /v1/units/{unit}/meaning` endpoints added.
- Verify-audit: `verify-audit --strict-hash` detects tampering of meaning sidecars by comparing canonicalized sidecar content to recorded `meaning_hash`.


See RELEASE_NOTES_v0.2.9.md for full details.

