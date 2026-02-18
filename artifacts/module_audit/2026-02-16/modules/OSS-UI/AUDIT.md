# OSS-UI Audit

## Scope
Minimal control-room: P0 dashboard, approval inbox, ProofGraph timeline. Shows deny receipts, budget, chain integrity.

## Reality
**Nothing exists.** No UI directory, no React/Next.js/Vite code, no dashboard.
- **Quality: N/A** — Missing entirely.

## Gaps
| # | Gap | Severity |
|---|-----|----------|
| 1 | No control-room UI | P1 |
| 2 | No API for receipts/proofgraph/approvals | P1 |

## Recommendations
1. **[NEW]** `apps/control-room/` — Minimal Vite/React app.
2. 3 views: P0 dashboard, approvals inbox, ProofGraph timeline.
3. **[NEW]** REST API endpoints for UI data.
