# OSS-PROXY Gaps

| # | Gap | Severity | Impact |
|---|-----|----------|--------|
| 1 | Handler stub does not forward to upstream | P0 | Server-mode proxy non-functional |
| 2 | Tool calls not routed through kernel pipeline | P0 | No governance on tool execution |
| 3 | No agentic re-invocation loop | P0 | Multi-turn agents break |
| 4 | Receipts use JSONL not ProofGraph | P0 | No DAG chain integrity |
| 5 | Reason codes are ad-hoc | P0 | Not auditor-grade |
| 6 | No budget enforcement | P0 | Budget deny demo impossible |
| 7 | Missing tool_calls/tools in request/response types | P1 | Wire format incomplete |
| 8 | No SSE streaming | P2 | Modern SDK clients broken |
| 9 | Two divergent proxy implementations | P1 | Maintenance burden |
