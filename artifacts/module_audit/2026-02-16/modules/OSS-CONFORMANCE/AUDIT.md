# OSS-CONFORMANCE Audit

## Scope
`helm conform` runner with profiles, gates, EvidencePack output. Deterministic JSON. Level L1/L2. CI integration.

## Reality
- **engine.go** (274 lines): Engine with RegisterGate, Run, EvidencePack output, 6 profiles. [CODE]
- **reason_codes.go** (84 lines): 23 stable reason codes with AllReasonCodes(). [CODE]
- **conform.go** CLI (139 lines): --profile, --json, --gates. [CODE]
- 24 gates (G0-G12, Gx_envelope, Gx_tenant). [CODE]
- **Quality: OK** — Engine functional, output not byte-deterministic.

## Gaps
| # | Gap | Severity |
|---|-----|----------|
| 1 | No `--level L1\|L2` alias | P1 |
| 2 | Output not byte-deterministic (`json.MarshalIndent`) | P1 |
| 3 | No git commit / env fingerprint in output | P1 |
| 4 | No `helm conform` CI job | P1 |

## Recommendations
1. Add `--level L1|L2` mapping to gate subsets.
2. Use JCS for output canonicalization.
3. Add git commit + env fingerprint to report.
4. Add CI job.
