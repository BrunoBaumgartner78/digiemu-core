# digiemu-core (Kernel)

[![CI](https://github.com/BrunoBaumgartner78/digiemu-core/actions/workflows/ci.yml/badge.svg)](https://github.com/BrunoBaumgartner78/digiemu-core/actions/workflows/ci.yml)

## Run tests
```bash
go test ./...
```

## Adapters

FS adapter (development JSON store):

- Location: `internal/kernel/adapters/fs`
- Usage: instantiate with a base path. Example:

	repo := fs.NewUnitRepo("./data")

This adapter stores each unit (and its versions) as a JSON file under the provided base path. It's intended for prototyping and local development.