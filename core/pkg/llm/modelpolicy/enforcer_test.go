package modelpolicy

import (
	"context"
	"testing"
)

func TestEnforcer_AllowedProviders(t *testing.T) {
	e := NewEnforcer()

	policy := &Policy{
		PolicyID: "test-providers",
		Version:  PolicyVersion,
		Name:     "Provider Policy",
		Enabled:  true,
		ModelConstraints: ModelConstraints{
			AllowedProviders: []string{"openai", "anthropic"},
		},
		Enforcement: Enforcement{
			Mode: EnforceModeEnforce,
		},
	}

	if err := e.LoadPolicy(policy); err != nil {
		t.Fatalf("LoadPolicy failed: %v", err)
	}

	// Allowed provider
	result := e.CheckRequest(context.Background(), "openai", "gpt-4", 1000, 500, 0.7, 0.05)
	if !result.Allowed {
		t.Error("OpenAI should be allowed")
	}

	// Denied provider
	result = e.CheckRequest(context.Background(), "unknown", "model-x", 1000, 500, 0.7, 0.05)
	if result.Allowed {
		t.Error("Unknown provider should be denied")
	}
	if len(result.Violations) == 0 {
		t.Error("Expected violation for unknown provider")
	}
}

func TestEnforcer_DeniedModels(t *testing.T) {
	e := NewEnforcer()

	policy := &Policy{
		PolicyID: "test-denied",
		Version:  PolicyVersion,
		Name:     "Denied Models",
		Enabled:  true,
		ModelConstraints: ModelConstraints{
			DeniedModels: []string{"gpt-3.5-turbo", "openai/davinci"},
		},
		Enforcement: Enforcement{
			Mode: EnforceModeEnforce,
		},
	}

	if err := e.LoadPolicy(policy); err != nil {
		t.Fatalf("LoadPolicy failed: %v", err)
	}

	// Denied model
	result := e.CheckRequest(context.Background(), "openai", "gpt-3.5-turbo", 1000, 500, 0.7, 0.01)
	if result.Allowed {
		t.Error("gpt-3.5-turbo should be denied")
	}

	// Allowed model
	result = e.CheckRequest(context.Background(), "openai", "gpt-4", 1000, 500, 0.7, 0.05)
	if !result.Allowed {
		t.Error("gpt-4 should be allowed")
	}
}

func TestEnforcer_TokenLimits(t *testing.T) {
	e := NewEnforcer()

	policy := &Policy{
		PolicyID: "test-tokens",
		Version:  PolicyVersion,
		Name:     "Token Limits",
		Enabled:  true,
		ModelConstraints: ModelConstraints{
			MaxContextTokens: 8000,
			MaxOutputTokens:  2000,
		},
		Enforcement: Enforcement{
			Mode: EnforceModeEnforce,
		},
	}

	if err := e.LoadPolicy(policy); err != nil {
		t.Fatalf("LoadPolicy failed: %v", err)
	}

	// Within limits
	result := e.CheckRequest(context.Background(), "openai", "gpt-4", 4000, 1000, 0.7, 0.05)
	if !result.Allowed {
		t.Error("Request within limits should be allowed")
	}

	// Exceeds context limit
	result = e.CheckRequest(context.Background(), "openai", "gpt-4", 10000, 1000, 0.7, 0.10)
	if result.Allowed {
		t.Error("Request exceeding context limit should be denied")
	}

	// Exceeds output limit
	result = e.CheckRequest(context.Background(), "openai", "gpt-4", 4000, 3000, 0.7, 0.10)
	if result.Allowed {
		t.Error("Request exceeding output limit should be denied")
	}
}

func TestEnforcer_BudgetConstraints(t *testing.T) {
	e := NewEnforcer()

	policy := &Policy{
		PolicyID:         "test-budget",
		Version:          PolicyVersion,
		Name:             "Budget Policy",
		Enabled:          true,
		ModelConstraints: ModelConstraints{},
		BudgetConstraints: &BudgetConstraints{
			DailyBudgetUSD:   10.0,
			PerRequestMaxUSD: 1.0,
			HardStopAtBudget: true,
		},
		Enforcement: Enforcement{
			Mode: EnforceModeEnforce,
		},
	}

	if err := e.LoadPolicy(policy); err != nil {
		t.Fatalf("LoadPolicy failed: %v", err)
	}

	// Within per-request max
	result := e.CheckRequest(context.Background(), "openai", "gpt-4", 1000, 500, 0.7, 0.50)
	if !result.Allowed {
		t.Error("Request within budget should be allowed")
	}

	// Exceeds per-request max
	result = e.CheckRequest(context.Background(), "openai", "gpt-4", 1000, 500, 0.7, 1.50)
	if result.Allowed {
		t.Error("Request exceeding per-request max should be denied")
	}
}

func TestEnforcer_TemperatureLimit(t *testing.T) {
	e := NewEnforcer()

	policy := &Policy{
		PolicyID:         "test-temp",
		Version:          PolicyVersion,
		Name:             "Temperature Policy",
		Enabled:          true,
		ModelConstraints: ModelConstraints{},
		Quality: &QualityRequirements{
			MaxTemperature: 0.5,
		},
		Enforcement: Enforcement{
			Mode: EnforceModeEnforce,
		},
	}

	if err := e.LoadPolicy(policy); err != nil {
		t.Fatalf("LoadPolicy failed: %v", err)
	}

	// Within temperature limit
	result := e.CheckRequest(context.Background(), "openai", "gpt-4", 1000, 500, 0.3, 0.05)
	if !result.Allowed {
		t.Error("Request within temperature limit should be allowed")
	}

	// Exceeds temperature limit
	result = e.CheckRequest(context.Background(), "openai", "gpt-4", 1000, 500, 0.8, 0.05)
	if result.Allowed {
		t.Error("Request exceeding temperature limit should be denied")
	}
}

func TestEnforcer_AuditMode(t *testing.T) {
	e := NewEnforcer()

	policy := &Policy{
		PolicyID: "test-audit",
		Version:  PolicyVersion,
		Name:     "Audit Policy",
		Enabled:  true,
		ModelConstraints: ModelConstraints{
			DeniedModels: []string{"gpt-3.5-turbo"},
		},
		Enforcement: Enforcement{
			Mode: EnforceModeAudit,
		},
	}

	if err := e.LoadPolicy(policy); err != nil {
		t.Fatalf("LoadPolicy failed: %v", err)
	}

	// Should be allowed in audit mode even with violation
	result := e.CheckRequest(context.Background(), "openai", "gpt-3.5-turbo", 1000, 500, 0.7, 0.01)
	if !result.Allowed {
		t.Error("Request should be allowed in audit mode")
	}
	if !result.AuditOnly {
		t.Error("Should indicate audit only")
	}
	if len(result.Violations) == 0 {
		t.Error("Should still record violations")
	}
}

func TestEnforcer_FallbackChain(t *testing.T) {
	e := NewEnforcer()

	policy := &Policy{
		PolicyID: "test-fallback",
		Version:  PolicyVersion,
		Name:     "Fallback Policy",
		Enabled:  true,
		ModelConstraints: ModelConstraints{
			DeniedModels: []string{"gpt-4-turbo"},
		},
		Fallback: &FallbackConfig{
			Enabled:       true,
			FallbackChain: []string{"gpt-4", "gpt-3.5-turbo"},
		},
		Enforcement: Enforcement{
			Mode:       EnforceModeEnforce,
			FailAction: FailActionFallback,
		},
	}

	if err := e.LoadPolicy(policy); err != nil {
		t.Fatalf("LoadPolicy failed: %v", err)
	}

	result := e.CheckRequest(context.Background(), "openai", "gpt-4-turbo", 1000, 500, 0.7, 0.05)
	if result.Allowed {
		t.Error("Denied model should not be allowed")
	}
	if result.FallbackModel != "gpt-4" {
		t.Errorf("Expected fallback to gpt-4, got %s", result.FallbackModel)
	}
}

func TestEnforcer_UsageTracking(t *testing.T) {
	e := NewEnforcer()

	e.RecordUsage(1.50, 5000)
	e.RecordUsage(0.75, 2500)

	stats := e.GetUsageStats()
	if dailySpend, ok := stats["daily_spend_usd"].(float64); !ok || dailySpend != 2.25 {
		t.Errorf("Expected daily spend 2.25, got %v", stats["daily_spend_usd"])
	}
}

func TestEnforcer_ConcurrentTracking(t *testing.T) {
	e := NewEnforcer()

	e.AcquireConcurrent()
	e.AcquireConcurrent()

	stats := e.GetUsageStats()
	if concurrent, ok := stats["concurrent"].(int64); !ok || concurrent != 2 {
		t.Errorf("Expected concurrent 2, got %v", stats["concurrent"])
	}

	e.ReleaseConcurrent()
	stats = e.GetUsageStats()
	if concurrent, ok := stats["concurrent"].(int64); !ok || concurrent != 1 {
		t.Errorf("Expected concurrent 1, got %v", stats["concurrent"])
	}
}

func TestEnforcer_ViolationHandler(t *testing.T) {
	e := NewEnforcer()

	var handlerCalled bool
	e.AddViolationHandler(func(ctx context.Context, result CheckResult) {
		handlerCalled = true
	})

	policy := &Policy{
		PolicyID: "test-handler",
		Version:  PolicyVersion,
		Name:     "Handler Policy",
		Enabled:  true,
		ModelConstraints: ModelConstraints{
			DeniedModels: []string{"denied-model"},
		},
		Enforcement: Enforcement{
			Mode: EnforceModeEnforce,
		},
	}

	if err := e.LoadPolicy(policy); err != nil {
		t.Fatalf("LoadPolicy failed: %v", err)
	}

	e.CheckRequest(context.Background(), "provider", "denied-model", 1000, 500, 0.7, 0.05)

	if !handlerCalled {
		t.Error("Violation handler should be called")
	}
}

func TestEnforcer_DisabledPolicy(t *testing.T) {
	e := NewEnforcer()

	policy := &Policy{
		PolicyID: "disabled",
		Version:  PolicyVersion,
		Name:     "Disabled",
		Enabled:  false,
		ModelConstraints: ModelConstraints{
			DeniedModels: []string{"all-models"},
		},
		Enforcement: Enforcement{
			Mode: EnforceModeEnforce,
		},
	}

	if err := e.LoadPolicy(policy); err != nil {
		t.Fatalf("LoadPolicy failed: %v", err)
	}

	result := e.CheckRequest(context.Background(), "any", "all-models", 1000, 500, 0.7, 0.05)
	if !result.Allowed {
		t.Error("Disabled policy should not block requests")
	}
}
