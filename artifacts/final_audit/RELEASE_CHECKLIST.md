# Release Readiness Checklist — HELM OSS v0.1

## Build & Test
- [x] `go build ./...` — clean
- [x] `go test ./pkg/... -count=1` — all packages pass
- [x] `go vet ./...` — clean
- [ ] `go test -race ./pkg/...` — not run (requires extended CI time)

## Crucible Use Cases
- [ ] UC-001: PEP Allow
- [ ] UC-002: PEP Fail-Closed
- [ ] UC-003: Approval Ceremony
- [ ] UC-004: WASM Transform
- [ ] UC-005: WASM Exhaustion
- [ ] UC-006: Idempotency
- [ ] UC-007: Export CLI Build
- [ ] UC-008: Replay CLI Build
- [ ] UC-009: Output Drift
- [ ] UC-010: Trust Rotation
- [ ] UC-011: Island Mode
- [ ] UC-012: Conformance Gates

> **Note:** Use cases require Docker + running services. Run `make crucible` to validate.

## Docker Path
- [ ] `docker build -t helm:latest .` — root Dockerfile
- [ ] `docker compose up -d` — dev stack
- [ ] Health check responds

## Docs Verified
- [x] README.md — commands work, ports match
- [x] DEMO.md — commands produce expected outputs
- [x] QUICKSTART.md — 8-step proof loop verified
- [x] Internal links — all valid after rename

## Security
- [x] No leaked secrets (API keys, tokens, private keys)
- [x] Docker compose passwords labeled DEV ONLY
- [x] `.env` is gitignored
- [x] No debug `fmt.Print*` in production code

## TCB Isolation
- [x] 8 TCB packages: zero `net/http` imports
- [x] 8 TCB packages: zero `os/exec` imports
- [x] 8 TCB packages: zero vendor SDK imports
- [x] No TODO/FIXME/HACK in production code

## Cutline Truth
- [x] `docs/OSS_CUTLINE.md` matches shipped packages
- [x] No overclaims

## Supply Chain
- [ ] `sbom.json` — generate with `make sbom`
- [x] LICENSE present (BUSL-1.1)
- [x] SECURITY.md present (clear reporting path)

## Repo Cleanliness
- [x] No stale binaries (41MB `core/helm` deleted)
- [x] No vendor directories
- [x] `.gitignore` covers bin/, .env, vendor/, node_modules/
