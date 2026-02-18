# OSS-PROXY Recommendations

1. **Converge on proxy_cmd.go path** - Delete or rewrite `openai_proxy.go` to delegate to the CLI proxy logic. The CLI proxy already has upstream forwarding and tool interception.

2. **Wire kernel pipeline** - Import and initialize Guardian, Executor, ProofGraph in `runProxyCmd`. For each tool_call: `guardian.EvaluateDecision` -> `guardian.IssueExecutionIntent` -> `executor.Execute` -> `proofgraph.Append`.

3. **Implement agentic loop** - Max 10 iterations. Forward -> intercept tool_calls -> govern each -> collect results -> re-invoke LLM. Terminate on `finish_reason: stop` or max iterations.

4. **Replace JSONL with ProofGraph** - `receiptStore` -> `proofgraph.Graph` + filesystem persistence. Each receipt becomes a DAG node.

5. **Use stable reason codes** - Map all deny outcomes to `conform.Reason*` constants. Add proxy-specific codes. Return standard error JSON.

6. **Add tool_calls/tools to wire format** - Extend `OpenAIChatRequest` and `OpenAIChatResponse` with `tools`, `tool_calls`, `tool_choice`, `function_call` fields per OpenAI API spec.

7. **Wire budget enforcement** - Initialize `budget.SimpleEnforcer` in proxy. Check before each tool execution. Add `--daily-limit`, `--monthly-limit` CLI flags.
