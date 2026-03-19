
# Wolfram CAG CLI Implementation Ledger

## Goal
Implement a Go CLI that calls Wolfram CAG APIs with Cobra, typed errors, sync HTTP, and built-in concurrency hooks.

## Technical anchors
- Base URL: `https://services.wolfram.com/api/cag/v1`
- Endpoints:
  - `POST /WolframAlphaContext`
  - `GET /WolframAlphaResult`
  - `POST /WolframLanguageCompute`
  - `POST /WolframLanguageHints`
- Auth header: `Authorization: <API key>`
- Env var fallback: `WOLFRAM_APP_ID`

## Phases

### Phase 1 — CLI skeleton + config wiring ✅ DONE
- [x] Initialize module + Cobra command structure
- [x] Add root/global flags
- [x] Add subcommands: `context`, `result`, `compute`, `hints`
- [x] Config resolution precedence (`--api-key` > env)
- [x] Output mode validation (`text|json`)

### Phase 2 — Typed errors + shared HTTP client ✅ DONE
- [x] Create `errs` package with static error types
- [x] Map request errors to timeout/network categories
- [x] Create shared sync `client.Client` wrapper
- [x] Add JSON decode helper
- [x] Wire resolved config/client in root pre-run

### Phase 3 — Endpoint handlers (current)
- [ ] Implement API calls for each subcommand
- [ ] `context`: POST `{context}`
- [ ] `result`: GET `input` (+ optional params later)
- [ ] `compute`: POST `{code}`
- [ ] `hints`: POST `{context}`
- [ ] Print responses in text/json modes

### Phase 4 — Output model + UX hardening
- [ ] Define response structs for known fields
- [ ] Consistent `stderr` formatting and exit behavior
- [ ] Add support for richer optional request params

### Phase 5 — Concurrency from day 1 (batch mode)
- [ ] Add `--input-file` processing for `result` and `compute`
- [ ] Worker pool with goroutines/channels
- [ ] Preserve deterministic output ordering

### Phase 6 — Tests + quality gates
- [ ] Unit tests for config and validation
- [ ] Integration tests with `httptest.Server`
- [ ] Fuzz tests for decode/parsing surfaces
- [ ] `go test -race ./...` clean

## Current file layout
```text
cmd/wolfram-cag/main.go
cli/
config/
client/
errs/
```

## Validation commands
```bash
go build ./...
go test ./...
go test -race ./...
```

## Notes
- Keep implementation simple and explicit.
- Avoid adding dependencies unless justified.
- Preserve typed error handling without `anyhow`-style wrappers.