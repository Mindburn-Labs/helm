# OSS-SCHEMAS-MANIFEST Audit

## Scope
Pinned tool schemas in manifest. JCS canonicalization (RFC 8785). Schema validation for tool args and outputs.

## Reality
- **manifest.go**: ToolSchema definitions, hash computation. [CODE]
- **canon.go**: JCS canonicalization per RFC 8785. [CODE]
- **validator.go**: ValidateToolArgs against pinned schemas. [CODE]
- Executor has OutputSchemaRegistry interface. Proxy uses JCS but NOT pinned schemas.
- **Quality: OK** — Exists, not wired into proxy.

## Gaps
| # | Gap | Severity |
|---|-----|----------|
| 1 | Proxy uses JCS hash only, not ValidateToolArgs | P0 |
| 2 | No tool output schema validation in proxy | P1 |

## Recommendations
1. Wire `manifest.ValidateToolArgs` into proxy tool_call handling.
2. Populate OutputSchemaRegistry from manifest in proxy context.
