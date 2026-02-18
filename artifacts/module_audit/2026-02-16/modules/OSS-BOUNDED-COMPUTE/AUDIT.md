# OSS-BOUNDED-COMPUTE Audit

## Scope
Budget enforcement (daily/monthly). WASI sandbox (gas/time/memory). Fail-closed.

## Reality
- **enforcer.go** (150 lines): SimpleEnforcer with Check, Memory + Postgres storage, fail-closed. [CODE]
- **sandbox.go**: WASI sandbox with gas metering, time/memory limits. [CODE]
- CI "sandbox" job tests WASI isolation. [CODE]
- **Quality: OK** — Both functional, not wired into proxy.

## Gaps
| # | Gap | Severity |
|---|-----|----------|
| 1 | Budget not wired into proxy | P0 |
| 2 | No max_iterations in proxy | P0 |
| 3 | No wallclock timeout in proxy | P1 |

## Recommendations
1. Wire `budget.SimpleEnforcer.Check` into proxy with `--daily-limit`, `--monthly-limit`.
2. Add `--max-iterations` (default 10) and `--max-wallclock` (default 120s).
