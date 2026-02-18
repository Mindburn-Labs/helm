# OSS-PROFILES Audit

## Scope
Per-region profiles (US, EU, RU, CN). Config loader. Runtime enforcement: networking, upstream allowlists, island mode, crypto policy, retention, ceremony.

## Reality
- **regional.yaml** (58 lines): Single file with all 4 regions. Ceremony, data_residency, compliance, encryption. [CODE]
- **config.go** (~945B): Minimal, no profile loading or enforcement. [CODE]
- **Quality: Partial** — Profiles documented, not enforced at runtime.

## Gaps
| # | Gap | Severity |
|---|-----|----------|
| 1 | Single file, not split per region | P1 |
| 2 | No config loader | P1 |
| 3 | No runtime enforcement | P1 |
| 4 | Missing networking/allowlist/island/retention fields | P1 |

## Recommendations
1. Split to `profile_{us,eu,ru,cn}.yaml` with full config.
2. Implement YAML loader in `config.go`.
3. Add enforcement functions.
