# OSS-CI-GATES Audit

## Scope
CI enforcing build, test, conformance, SBOM, provenance, linting. Fail-on-regression.

## Reality
- **helm_core_gates.yml** (252 lines): 15 jobs — build, unit tests, fuzz, PEP, proofgraph, sandbox, ceremony, evidence-pack, TCB isolation, race detection, use-cases, lint (TODO check), SBOM, provenance, doccheck. Go 1.24. [CODE]
- **release.yml**: Release workflow. **sdk_gates.yml**: SDK CI. [CODE]
- **Quality: OK** — Comprehensive, missing helm conform job.

## Gaps
| # | Gap | Severity |
|---|-----|----------|
| 1 | No `helm conform` CI job | P1 |
| 2 | No conformance determinism check | P1 |

## Recommendations
1. Add `helm conform --level L1 --json` CI job.
2. Add determinism check: run twice, diff output.
