# 08 — Supply Chain & Releases

**Score: 4/5** · Gate ≥3 · **✅ PASS**

---

## CI Pipeline (Verified: `helm_core_gates.yml`, 336 lines, 16 jobs)

### Job Matrix

| Job | Lines | Purpose | Status |
|-----|-------|---------|--------|
| `build` | L16-25 | `go build ./...` | ✅ |
| `test` | L27-38 | `go test ./pkg/... -cover` | ✅ |
| `sandbox` | L40-71 | WASI sandbox tests | ✅ |
| `ceremony` | L73-88 | Ceremony/escalation tests | ✅ |
| `race` | L90-121 | Race detector (`-race -count=1`) | ✅ |
| `use-cases` | L123-136 | `scripts/usecases/run_all.sh` | ✅ |
| `lint` | L138-172 | golangci-lint + govulncheck + TODO check | ✅ |
| `sbom` | L174-197 | CycloneDX SBOM generation + validation | ✅ |
| `provenance` | L199-217 | SHA-256 build provenance | ✅ |
| `doc-hash` | L219-225 | Canonical doc hash verification | ✅ |
| `sdk-build` | L227-256 | TS (vitest) + Python (pytest) + Rust (cargo) | ✅ |
| `doccheck` | L258-263 | `tools/doccheck/main.go` | ✅ |
| `conformance-gate` | L265-300 | Conformance L1/L2 gate | ✅ |
| `tcb-check` | L284-301 | TCB boundary verification | ✅ |
| `evidence-determinism` | L303-328 | EvidencePack determinism | ⚠️ No-op (see below) |
| `examples-smoke` | L330-335 | `scripts/ci/examples_smoke.sh` | ✅ |

### Lint Configuration

```yaml
# L146-154
- name: golangci-lint
  uses: golangci/golangci-lint-action@v5
  with:
    version: latest
    args: --timeout=5m
- name: govulncheck
  run: |
    go install golang.org/x/vuln/cmd/govulncheck@latest
    govulncheck ./...
```

golangci-lint@v5 and govulncheck ARE configured. Zero-TODO policy enforced (L155-172).

### Evidence-Determinism Test Is a No-Op

```bash
# L316-317 — swallows the error with || echo
./bin/helm export --evidence ./data/evidence --out /tmp/pack1.tar.gz 2>/dev/null \
    || echo "Export requires running kernel (skip in CI dry-run)"
```

The test can't run without a live kernel, so it always falls through to the `echo` branch. The comparison logic (L318-328) never executes. **This is a false gate.**

---

## Release Pipeline (Verified: `release.yml`)

### Architecture
```
Tag Push (v*)
    │
    ├── Cross-compile (linux/darwin × amd64/arm64)
    │     └── CGO_ENABLED=0, -ldflags="-s -w"
    ├── SHA-256 checksums → SHA256SUMS.txt
    ├── CycloneDX SBOM → sbom.json
    ├── Docker multi-arch → GHCR
    │     ├── gcr.io/distroless/static-debian12:nonroot
    │     └── OCI labels (source, revision, version)
    ├── GitHub Release
    ├── SLSA Provenance (actions/attest-build-provenance@v2)
    └── SDK Publishing
          ├── PyPI (Trusted Publishing / OIDC)
          ├── npm (--provenance flag)
          ├── crates.io
          └── Maven Central
```

### Supply Chain Scorecard

| Criterion | Status | Evidence |
|-----------|--------|----------|
| golangci-lint | ✅ | `helm_core_gates.yml` L146-150 |
| govulncheck | ✅ | `helm_core_gates.yml` L151-154 |
| Race detector | ✅ | `helm_core_gates.yml` L121: `-race` |
| SLSA provenance | ✅ | `actions/attest-build-provenance@v2` |
| CycloneDX SBOM | ✅ | Generated and published |
| SPDX SBOM | ❌ | Only CycloneDX |
| Distroless runtime | ✅ | `gcr.io/distroless/static-debian12:nonroot` |
| Multi-arch images | ✅ | linux/amd64 + linux/arm64 |
| NPM provenance | ✅ | `--provenance` flag |
| PyPI trusted publishing | ✅ | OIDC `id-token: write` |
| CGO_ENABLED=0 | ✅ | Static binary |
| Cosign signatures | ❌ | SHA-256 only, no identity verification |
| OpenSSF Scorecard | ❌ | Not running |
| Container scanning | ❌ | No Trivy/Grype in CI |
| Pinned base image digest | ❌ | Uses tag, not digest |
| Reproducible builds | ⚠️ | `-ldflags="-s -w"` strips symbols |

---

## SLSA Level

| Requirement | Status | Gap |
|-------------|--------|-----|
| Build L1: Scripted build | ✅ | — |
| Build L2: Authenticated provenance | ✅ | `actions/attest-build-provenance@v2` |
| Build L3: Hardened platform | ⚠️ | GitHub-hosted runners (not isolated) |
| Source L1: Version controlled | ✅ | Git |
| Dependencies L1: Listed | ✅ | go.mod + SBOM |

**Assessment:** SLSA Build L2.

---

## Dockerfile Security

```dockerfile
FROM golang:1.24-alpine AS builder     # ✅ Multi-stage
RUN CGO_ENABLED=0 go build ...         # ✅ Static binary
FROM gcr.io/distroless/static-debian12:nonroot  # ✅ Distroless + nonroot
```

| Check | Status |
|-------|--------|
| Multi-stage build | ✅ |
| Distroless runtime | ✅ |
| Non-root user | ✅ |
| No secrets in image | ✅ |
| Image scanning | ❌ Not in CI |
| Pinned base digest | ❌ Uses tag |

---

## Score: 4/5

**Justification:**
- ✅ SLSA Build L2 provenance
- ✅ golangci-lint + govulncheck + race detector in CI
- ✅ CycloneDX SBOM + Trusted Publishing (OIDC)
- ✅ Distroless + nonroot + static binary
- ✅ 16-job CI pipeline with conformance, TCB check, ceremony tests
- ❌ No Cosign signatures
- ❌ No OpenSSF Scorecard
- ❌ Evidence-determinism test is a no-op
- ❌ No container image scanning

### To reach 5/5:
1. Add Cosign signing for binaries and Docker images
2. Enable OpenSSF Scorecard
3. Fix evidence-determinism test (needs mock kernel or integration env)
4. Add Trivy/Grype container scanning
5. Pin base images by digest
