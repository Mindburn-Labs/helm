# OSS-EXEC-LOOP Scope
Forward model request to upstream, intercept tool_calls, route each through KernelBridge -> Guardian -> Executor, re-invoke LLM with tool outputs until completion. Strict limits: max tool iterations, max wallclock, bounded compute budgets.
