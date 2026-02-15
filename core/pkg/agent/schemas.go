package agent

// Tool definitions exposed to the LLM.

type KernelToolName string

const (
	ToolCreateObligation  KernelToolName = "create_obligation"
	ToolSearchObligations KernelToolName = "search_obligations"
	ToolProposePlan       KernelToolName = "propose_plan"
	ToolRequestDecision   KernelToolName = "request_decision"
	ToolMCPToolSearch     KernelToolName = "mcp_tool_search"
	ToolCallMCPTool       KernelToolName = "call_mcp_tool"
	ToolSubmitModule      KernelToolName = "submit_module_bundle"
	ToolRequestActivation KernelToolName = "request_module_activation"
)

// GoalSpec represents the LLM's understanding of intent.
type GoalSpec struct {
	Intent          string   `json:"intent"`
	SuccessCriteria []string `json:"success_criteria"`
}

// CreateObligationParams
type CreateObligationParams struct {
	Intent         string `json:"intent"`
	IdempotencyKey string `json:"idempotency_key"`
}

// ProposePlanParams
type ProposePlanParams struct {
	ObligationID string         `json:"obligation_id"`
	Steps        []PlanStepSpec `json:"steps"`
}

type PlanStepSpec struct {
	ToolName  string         `json:"tool_name"`
	Inputs    map[string]any `json:"inputs"`
	Reasoning string         `json:"reasoning"`
}

// CallMCPToolParams
type CallMCPToolParams struct {
	ToolName   string         `json:"tool_name"`
	DecisionID string         `json:"decision_id"` // REQUIRED: The gating token
	Params     map[string]any `json:"params"`
}
