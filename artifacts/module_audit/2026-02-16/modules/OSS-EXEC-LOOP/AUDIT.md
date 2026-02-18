# OSS-EXEC-LOOP Audit

## Scope
Forward model request to upstream, intercept tool_calls, route each through KernelBridge -> Guardian -> Executor, re-invoke LLM with tool outputs until completion. Strict limits: max tool iterations, max wallclock, bounded compute budgets.

## Reality
- `proxy_cmd.go` (388 lines): Single-pass upstream forward + tool_call interception. No loop. [CODE]
- `core/pkg/kernel/`: 85 files with cybernetics loop, effect boundary, interop. [CODE]
- No `KernelBridge` abstraction materialized as named component.
- Missing: agentic re-invocation loop, max_iterations, wallclock timeout.

## Conformance
| Check | Status |
|-------|--------|
| Upstream forwarding | ✅ proxy CLI |
| Tool call interception | ✅ proxy CLI |
| KernelBridge abstraction | ❌ missing |
| Guardian/Executor routing | ❌ not wired |
| Re-invocation loop | ❌ missing |
| Max iterations limit | ❌ missing |
| Wallclock budget | ❌ missing |

## Gaps
| # | Gap | Severity |
|---|-----|----------|
| 1 | No agentic loop (tool_call -> execute -> re-invoke) | P0 |
| 2 | No KernelBridge routing to Guardian -> Executor | P0 |
| 3 | No max_iterations enforcement | P0 |
| 4 | No wallclock timeout in proxy | P1 |

## Recommendations
1. Create `KernelBridge` adapter: accepts tool_call, calls Guardian, calls Executor, returns result.
2. Implement agentic loop in proxy with `--max-iterations` (default 10) and `--max-wallclock` (default 120s).
3. **[NEW]** `core/pkg/bridge/kernel_bridge.go`
4. **[MODIFY]** `core/cmd/helm/proxy_cmd.go` — add loop + bridge
