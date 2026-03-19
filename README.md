# wolfram-cag

A Go CLI for the Wolfram Conversational API Gateway (CAG) endpoints.

This tool provides command-line access to:

- `WolframAlphaContext` (POST)
- `WolframAlphaResult` (GET)
- `WolframLanguageCompute` (POST)
- `WolframLanguageHints` (POST)

Base API URL (default):

- `https://services.wolfram.com/api/cag/v1`

---

## Features

- Cobra-based CLI with subcommands
- API key via flag or environment variable
- Typed error handling
- `text` and `json` output modes
- Batch mode with worker pool for `result` and `compute`
- Unit, integration, fuzz, and race-tested code paths

---

## Requirements

- Go 1.22+
- A Wolfram API key

---

## Build

```bash
go build ./cmd/wolfram-cag
```

Run directly:

```bash
go run ./cmd/wolfram-cag --help
```

---

## Authentication

You can provide the API key in either way:

1. `--api-key` flag (highest priority)
2. `WOLFRAM_APP_ID` environment variable

Examples:

```bash
export WOLFRAM_APP_ID="your-api-key"
go run ./cmd/wolfram-cag context --context "what is quantum tunneling?"
```

or:

```bash
go run ./cmd/wolfram-cag --api-key "your-api-key" context --context "what is quantum tunneling?"
```

If neither is set, the CLI exits with an auth error.

---

## Global Flags

```text
--api-key string       Wolfram API key (overrides WOLFRAM_APP_ID)
--base-url string      API base URL (default https://services.wolfram.com/api/cag/v1)
--output string        text|json (default text)
--timeout-secs int     HTTP timeout in seconds (default 30)
--workers int          Worker count for batch operations (default 4)
--verbose              Print resolved runtime config
```

---

## Commands

### 1) `context`
Call WolframAlphaContext API.

```bash
go run ./cmd/wolfram-cag context \
  --context "I am helping a student learn calculus"
```

### 2) `result`
Call WolframAlphaResult API.

Single input:

```bash
go run ./cmd/wolfram-cag result --input "integrate x^2"
```

With optional query parameters:

```bash
go run ./cmd/wolfram-cag result \
  --input "weather in Boston" \
  --units metric \
  --location Boston \
  --format plaintext
```

Supported optional flags:

- `--assumption`
- `--format`
- `--units`
- `--location`
- `--latlong`
- `--timeout`
- `--maxwidth`

Batch mode:

```bash
go run ./cmd/wolfram-cag result --input-file ./queries.txt --workers 8
```

Where `queries.txt` is newline-delimited:

```text
integrate x^2
population of japan
distance from earth to moon
```

### 3) `compute`
Call WolframLanguageCompute API.

Single input:

```bash
go run ./cmd/wolfram-cag compute --code "Integrate[x^2, x]"
```

With optional compute parameters:

```bash
go run ./cmd/wolfram-cag compute \
  --code "Table[n^2, {n, 1, 10}]" \
  --time-constraint 5 \
  --line 1 \
  --max-chars 1000
```

Batch mode:

```bash
go run ./cmd/wolfram-cag compute --code-file ./wl_code.txt --workers 4
```

Where `wl_code.txt` is newline-delimited Wolfram Language expressions.

### 4) `hints`
Call WolframLanguageHints API.

```bash
go run ./cmd/wolfram-cag hints \
  --context "How do I compute eigenvalues in Wolfram Language?"
```

---

## Output Modes

- `--output text` (default): prints the primary `result` field when available.
- `--output json`: pretty-prints full JSON response.

Example:

```bash
go run ./cmd/wolfram-cag --output json result --input "integrate x^2"
```

---

## Error Behavior

The client maps failures into typed categories (network/timeout/http status/encoding/decoding/invalid args).

Common HTTP statuses from CAG APIs include:

- `400` bad input
- `401` missing/invalid API key
- `501` input not interpretable

In batch mode (`result` / `compute` with file input), the tool processes inputs concurrently and preserves output ordering; it reports per-item failures and returns a non-zero exit if any item fails.

---

## API Mapping

| CLI Command | Endpoint | Method |
|---|---|---|
| `context` | `/WolframAlphaContext` | `POST` |
| `result` | `/WolframAlphaResult` | `GET` |
| `compute` | `/WolframLanguageCompute` | `POST` |
| `hints` | `/WolframLanguageHints` | `POST` |

Requests include:

- Header: `Authorization: <api-key>`
- Header: `Content-Type: application/json` for POST requests

---

## Development

Run tests:

```bash
go test ./...
```

Run race detector:

```bash
go test -race ./...
```

Run fuzz tests (short example):

```bash
go test ./api -run=^$ -fuzz=Fuzz -fuzztime=2s
go test ./client -run=^$ -fuzz=Fuzz -fuzztime=2s
```

Format code:

```bash
gofmt -w ./...
```

---

## Notes

- This project is currently CLI-focused and intentionally lightweight.
- You can override `--base-url` for local mocks or testing servers.
