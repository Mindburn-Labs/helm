# OSS-DOCS-ADOPTION Audit

## Scope
Developer adoption docs: 1-line integration, use case docs with assertions, reference docs, threat model.

## Reality
- **README.md**: 1-line integration wedge (OPENAI_BASE_URL). [CODE]
- **INTEGRATE_IN_5_MIN.md**: Quick-start guide. [CODE]
- **DEMO.md**: Demo instructions. [CODE]
- **UC docs** (UC-001..UC-012): ~120B pointer files each. No assertions. [CODE]
- **THREAT_MODEL.md**: Exists, no proxy section. [CODE]
- **run_all.sh** (89 lines): UC orchestrator. [CODE]
- **Quality: Partial** — Skeleton exists, UC docs are stubs.

## Gaps
| # | Gap | Severity |
|---|-----|----------|
| 1 | UC docs are stubs (~120B each) | P2 |
| 2 | No REASON_CODES.md reference | P2 |
| 3 | No CONFORMANCE_LEVELS.md | P2 |
| 4 | No EVIDENCEPACK_VERIFY.md tutorial | P2 |
| 5 | No rogue agent demo script | P1 |
| 6 | Proxy not in threat model | P2 |

## Recommendations
1. Expand UC docs with input/output/assertions.
2. Generate REASON_CODES.md from `conform.AllReasonCodes()`.
3. Create CONFORMANCE_LEVELS.md and EVIDENCEPACK_VERIFY.md.
4. **[NEW]** `scripts/demo/rogue_agent.sh`.
5. Add proxy section to THREAT_MODEL.md.
