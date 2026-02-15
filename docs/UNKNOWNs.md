# HELM OSS v0.1 — Unknowns and Validation Plans
Date: 2026-02-15

## 1. ACID Kill Test (Phase 6)

**Status**: Not yet validated (requires live Postgres with concurrent connections)

**Validation Plan**:
1. Set up Postgres with `docker-compose up -d`
2. Run 100 parallel budget-constrained calls via `go test -race`
3. Assert exactly 1 success under budget=1
4. Assert restart simulation does not allow re-execution

**Risk**: Medium — the code structure is correct (receipt-based idempotency + FOR UPDATE), but no concurrency stress test exists yet.

---

## 2. Deterministic EvidencePack Tarball Bytes (Phase 9)

**Status**: Not yet validated with actual export

**Validation Plan**:
1. Run `helm export --evidence <dir> --out /tmp/pack1 --audit`
2. Run same command again to `/tmp/pack2`
3. `sha256sum pack1/* pack2/*` must match
4. Verify sorted paths, epoch mtime, stable uid/gid

**Risk**: Low — `export_pack.go` sorts entries and uses fixed mtime, but needs end-to-end test.

---

## 3. Offline Replay Verify (Phase 9)

**Status**: CLI wired, but not tested with network disabled

**Validation Plan**:
1. Export an EvidencePack from a real session
2. Disconnect network (or use `unshare --net`)
3. Run `helm replay --evidence <pack> --verify`
4. Must succeed without any network calls

**Risk**: Low — replay operates on local tape files only.

---

## 4. SDK Clients (Phase 12)

**Status**: Not yet implemented

**Validation Plan**:
1. Create `sdk/typescript/src/client.ts` with typed fetch-based client
2. Create `sdk/python/helm_client/client.py` with httpx-based client
3. Run examples against `docker compose` stack
4. Verify receipts are returned

**Risk**: Medium — no code exists yet.

---

## 5. CI Gate Expansion (Phase 14)

**Status**: 7 gates exist in `.github/workflows/helm_core_gates.yml`

**Validation Plan**:
1. Add ACID kill test gate
2. Add time-travel replay invariant gate
3. Add EvidencePack determinism gate
4. Verify all gates block PRs

**Risk**: Low — existing gate infrastructure works.

---

## 6. SBOM and Provenance (Phase 15)

**Status**: Not yet implemented

**Validation Plan**:
1. Add `cyclonedx-gomod` or `spdx-sbom-generator` to release pipeline
2. Generate SBOM as release artifact
3. Create `release_verify.sh` that validates checksums + SBOM
4. Add SLSA-style provenance attestation

**Risk**: Low — tooling is well-established.

---

## 7. Instruction/Fuel Metering in WASI (Phase 7)

**Status**: wazero does not natively expose instruction counting

**Validation Plan**:
1. Document limitation: wazero enforces time+memory but not instruction metering
2. ERR_COMPUTE_GAS_EXHAUSTED is implemented in budget package but not wired to actual WASM instruction counting
3. Mitigation: context.WithTimeout provides reliable termination

**Risk**: Low — documented limitation, not a blocker. Time limits are enforced.
