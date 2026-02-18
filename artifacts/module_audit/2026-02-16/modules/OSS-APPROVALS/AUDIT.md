# OSS-APPROVALS Audit

## Scope
RFC-005 approval ceremonies: timelock, hold, challenge, domain separation. HITL gates.

## Reality
- **ceremony.go**: CeremonyEnforcer with timelock, hold, challenge, domain separation. [CODE]
- G8_HITL conformance gate. Regional profiles define ceremony params. Tests exist. [CODE]
- **Quality: OK** — Mechanism exists, no UI inbox.

## Gaps
| # | Gap | Severity |
|---|-----|----------|
| 1 | No approval inbox UI | P1 |
| 2 | No REST API for approvals | P1 |

## Recommendations
1. REST API: `GET /api/approvals/pending`, `POST /api/approvals/:id/approve`.
2. Surface in control-room UI inbox.
