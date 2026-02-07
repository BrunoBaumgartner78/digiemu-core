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

## Running API

Run the minimal HTTP API (uses FS adapter):

```bash
DIGIEMU_DATA_DIR=./data DIGIEMU_ADDR=:8080 go run ./cmd/api
```

Create a unit:

```bash
curl -X POST -H "Content-Type: application/json" -d '{"name":"my-unit","description":"desc"}' http://localhost:8080/v1/units
```

Create a version:

```bash
curl -X POST -H "Content-Type: application/json" -d '{"content":"hello"}' http://localhost:8080/v1/units/my-unit/versions
```

## HTTP API (cmd/api)

Quickstart — start the server and use curl (copy/paste):

1) Start server:

```bash
go run ./cmd/api --addr :8080 --data ./data
```

2) Create a unit (returns `key` in response):

```bash
curl -s -X POST http://localhost:8080/v1/units \
	-H "Content-Type: application/json" \
	-d '{"name":"Demo Unit","description":"Demo"}'
```

3) Create a version for that unit (use returned `{key}`):

```bash
curl -s -X POST http://localhost:8080/v1/units/demo-unit/versions \
	-H "Content-Type: application/json" \
	-d '{"content":"v1"}'
```

Note: the create-unit response contains the `key` you should use in the versions endpoint.

## Local Quickstart (CLI + FS)

This repository includes a tiny CLI that uses the FS adapter for local demos.

Build and run the CLI examples below — the default data directory is `./data`.

Create a unit (auto-generates a key from the title when `--key` is omitted):

```bash
go run ./cmd/digiemu unit create --title "Demo Unit" --desc "Demo" --data ./data
```

Create a version for an existing unit (use the unit key):

```bash
go run ./cmd/digiemu version create --unit demo-unit --content "v1" --data ./data
```

Start the HTTP API (same FS storage used by the CLI):

```bash
go run ./cmd/digiemu serve --addr :8080 --data ./data
```

Flags supported by the CLI:

- `--data`: data directory (default `./data`)
- `--addr`: server address (for `serve`)
- `--title`, `--desc`: unit creation
- `--key`: optional explicit unit key
- `--unit`, `--content`: version creation
