// Package modelpolicy implements model gateway policy enforcement
// including budget ceilings, rate limits, and quality requirements.
package modelpolicy

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// PolicyVersion is the current schema version.
const PolicyVersion = "1.0.0"

// EnforcementMode defines how policy violations are handled.
type EnforcementMode string

const (
	EnforceModeEnforce  EnforcementMode = "enforce"
	EnforceModeAudit    EnforcementMode = "audit"
	EnforceModeDisabled EnforcementMode = "disabled"
)

// FailAction defines what happens on policy violation.
type FailAction string

const (
	FailActionBlock    FailAction = "block"
	FailActionFallback FailAction = "fallback"
	FailActionDegrade  FailAction = "degrade"
	FailActionAlert    FailAction = "alert"
)

// Policy defines a model gateway policy.
type Policy struct {
	PolicyID          string               `json:"policy_id"`
	Version           string               `json:"version"`
	Name              string               `json:"name"`
	Description       string               `json:"description,omitempty"`
	Enabled           bool                 `json:"enabled"`
	ModelConstraints  ModelConstraints     `json:"model_constraints"`
	BudgetConstraints *BudgetConstraints   `json:"budget_constraints,omitempty"`
	RateLimits        *RateLimits          `json:"rate_limits,omitempty"`
	Quality           *QualityRequirements `json:"quality_requirements,omitempty"`
	Fallback          *FallbackConfig      `json:"fallback_config,omitempty"`
	Enforcement       Enforcement          `json:"enforcement"`
}

// ModelConstraints defines allowed models.
type ModelConstraints struct {
	AllowedProviders        []string `json:"allowed_providers,omitempty"`
	AllowedModels           []string `json:"allowed_models,omitempty"`
	DeniedModels            []string `json:"denied_models,omitempty"`
	RequireStructuredOutput bool     `json:"require_structured_output"`
	RequireToolUse          bool     `json:"require_tool_use"`
	PreferLocal             bool     `json:"prefer_local"`
	MaxContextTokens        int      `json:"max_context_tokens,omitempty"`
	MaxOutputTokens         int      `json:"max_output_tokens,omitempty"`
}

// BudgetConstraints defines spending limits.
type BudgetConstraints struct {
	DailyBudgetUSD        float64 `json:"daily_budget_usd,omitempty"`
	MonthlyBudgetUSD      float64 `json:"monthly_budget_usd,omitempty"`
	PerRequestMaxUSD      float64 `json:"per_request_max_usd,omitempty"`
	BudgetPeriodReset     string  `json:"budget_period_reset,omitempty"`
	AlertThresholdPercent float64 `json:"alert_threshold_percent,omitempty"`
	HardStopAtBudget      bool    `json:"hard_stop_at_budget"`
}

// RateLimits defines request rate limits.
type RateLimits struct {
	RequestsPerMinute  int `json:"requests_per_minute,omitempty"`
	RequestsPerHour    int `json:"requests_per_hour,omitempty"`
	TokensPerMinute    int `json:"tokens_per_minute,omitempty"`
	ConcurrentRequests int `json:"concurrent_requests,omitempty"`
	BurstAllowance     int `json:"burst_allowance,omitempty"`
}

// QualityRequirements defines model quality requirements.
type QualityRequirements struct {
	MinModelVersion      string  `json:"min_model_version,omitempty"`
	RequireDeterministic bool    `json:"require_deterministic_mode"`
	RequireContentFilter bool    `json:"require_content_filter"`
	RequireCitation      bool    `json:"require_citation"`
	MaxTemperature       float64 `json:"max_temperature,omitempty"`
}

// FallbackConfig defines fallback behavior.
type FallbackConfig struct {
	Enabled             bool     `json:"enabled"`
	FallbackChain       []string `json:"fallback_chain,omitempty"`
	FallbackOn          []string `json:"fallback_on,omitempty"`
	MaxFallbackAttempts int      `json:"max_fallback_attempts"`
}

// Enforcement defines enforcement configuration.
type Enforcement struct {
	Mode       EnforcementMode `json:"mode"`
	FailAction FailAction      `json:"fail_action"`
}

// Violation represents a policy violation.
type Violation struct {
	PolicyID       string    `json:"policy_id"`
	ConstraintType string    `json:"constraint_type"`
	Description    string    `json:"description"`
	Severity       string    `json:"severity"`
	Timestamp      time.Time `json:"timestamp"`
}

// CheckResult is the outcome of a policy check.
type CheckResult struct {
	Allowed       bool        `json:"allowed"`
	Violations    []Violation `json:"violations,omitempty"`
	AuditOnly     bool        `json:"audit_only"`
	FallbackModel string      `json:"fallback_model,omitempty"`
}

// UsageTracker tracks budget and rate limit usage.
type UsageTracker struct {
	mu              sync.Mutex
	dailySpendUSD   float64
	monthlySpendUSD float64
	dailyReset      time.Time
	monthlyReset    time.Time

	minuteRequests int64
	hourRequests   int64
	minuteTokens   int64
	minuteReset    time.Time
	hourReset      time.Time

	concurrent int64
}

// Enforcer enforces model policies.
type Enforcer struct {
	mu       sync.RWMutex
	policies map[string]*Policy
	tracker  *UsageTracker
	handlers []ViolationHandler
}

// ViolationHandler is called when violations occur.
type ViolationHandler func(ctx context.Context, result CheckResult)

// NewEnforcer creates a new policy enforcer.
func NewEnforcer() *Enforcer {
	now := time.Now().UTC()
	return &Enforcer{
		policies: make(map[string]*Policy),
		tracker: &UsageTracker{
			dailyReset:   now.Truncate(24 * time.Hour).Add(24 * time.Hour),
			monthlyReset: time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC),
			minuteReset:  now.Add(time.Minute),
			hourReset:    now.Add(time.Hour),
		},
	}
}

// LoadPolicy adds or updates a policy.
func (e *Enforcer) LoadPolicy(policy *Policy) error {
	if policy.Version != PolicyVersion {
		return fmt.Errorf("unsupported policy version: %s", policy.Version)
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	e.policies[policy.PolicyID] = policy
	return nil
}

// LoadPolicyJSON removed - was dead code

// AddViolationHandler registers a handler for violations.
func (e *Enforcer) AddViolationHandler(h ViolationHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.handlers = append(e.handlers, h)
}

// CheckRequest validates a model request against policies.
func (e *Enforcer) CheckRequest(ctx context.Context, providerID, modelID string, inputTokens, outputTokens int, temperature float64, estimatedCostUSD float64) CheckResult {
	e.mu.RLock()
	defer e.mu.RUnlock()

	result := CheckResult{Allowed: true}

	for _, policy := range e.policies {
		if !policy.Enabled || policy.Enforcement.Mode == EnforceModeDisabled {
			continue
		}

		violations := e.checkPolicy(policy, providerID, modelID, inputTokens, outputTokens, temperature, estimatedCostUSD)
		result.Violations = append(result.Violations, violations...)

		if len(violations) > 0 {
			if policy.Enforcement.Mode == EnforceModeEnforce {
				result.Allowed = false

				// Check for fallback
				if policy.Fallback != nil && policy.Fallback.Enabled && len(policy.Fallback.FallbackChain) > 0 {
					result.FallbackModel = policy.Fallback.FallbackChain[0]
				}
			} else {
				result.AuditOnly = true
			}
		}
	}

	if len(result.Violations) > 0 {
		e.notifyHandlers(ctx, result)
	}

	return result
}

// checkPolicy checks a request against a single policy.
//
//nolint:gocognit // complexity acceptable
//nolint:gocyclo // complexity acceptable
func (e *Enforcer) checkPolicy(policy *Policy, providerID, modelID string, inputTokens, outputTokens int, temperature float64, estimatedCostUSD float64) []Violation {
	var violations []Violation

	// Model constraints
	mc := policy.ModelConstraints

	// Check allowed providers
	if len(mc.AllowedProviders) > 0 {
		allowed := false
		for _, p := range mc.AllowedProviders {
			if p == providerID {
				allowed = true
				break
			}
		}
		if !allowed {
			violations = append(violations, Violation{
				PolicyID:       policy.PolicyID,
				ConstraintType: "model.allowed_providers",
				Description:    fmt.Sprintf("Provider %s not in allowed list", providerID),
				Severity:       "high",
				Timestamp:      time.Now().UTC(),
			})
		}
	}

	// Check denied models
	for _, denied := range mc.DeniedModels {
		if denied == modelID || denied == providerID+"/"+modelID {
			violations = append(violations, Violation{
				PolicyID:       policy.PolicyID,
				ConstraintType: "model.denied_models",
				Description:    fmt.Sprintf("Model %s is denied", modelID),
				Severity:       "critical",
				Timestamp:      time.Now().UTC(),
			})
		}
	}

	// Check token limits
	if mc.MaxContextTokens > 0 && inputTokens > mc.MaxContextTokens {
		violations = append(violations, Violation{
			PolicyID:       policy.PolicyID,
			ConstraintType: "model.max_context_tokens",
			Description:    fmt.Sprintf("Input tokens %d exceeds max %d", inputTokens, mc.MaxContextTokens),
			Severity:       "medium",
			Timestamp:      time.Now().UTC(),
		})
	}

	if mc.MaxOutputTokens > 0 && outputTokens > mc.MaxOutputTokens {
		violations = append(violations, Violation{
			PolicyID:       policy.PolicyID,
			ConstraintType: "model.max_output_tokens",
			Description:    fmt.Sprintf("Output tokens %d exceeds max %d", outputTokens, mc.MaxOutputTokens),
			Severity:       "medium",
			Timestamp:      time.Now().UTC(),
		})
	}

	// Quality requirements
	if policy.Quality != nil {
		qr := policy.Quality
		if qr.MaxTemperature > 0 && temperature > qr.MaxTemperature {
			violations = append(violations, Violation{
				PolicyID:       policy.PolicyID,
				ConstraintType: "quality.max_temperature",
				Description:    fmt.Sprintf("Temperature %.2f exceeds max %.2f", temperature, qr.MaxTemperature),
				Severity:       "low",
				Timestamp:      time.Now().UTC(),
			})
		}
	}

	// Budget constraints
	if policy.BudgetConstraints != nil {
		bc := policy.BudgetConstraints
		e.tracker.mu.Lock()

		// Check per-request max
		if bc.PerRequestMaxUSD > 0 && estimatedCostUSD > bc.PerRequestMaxUSD {
			violations = append(violations, Violation{
				PolicyID:       policy.PolicyID,
				ConstraintType: "budget.per_request_max",
				Description:    fmt.Sprintf("Request cost $%.4f exceeds max $%.4f", estimatedCostUSD, bc.PerRequestMaxUSD),
				Severity:       "high",
				Timestamp:      time.Now().UTC(),
			})
		}

		// Check daily budget
		if bc.DailyBudgetUSD > 0 && e.tracker.dailySpendUSD+estimatedCostUSD > bc.DailyBudgetUSD {
			violations = append(violations, Violation{
				PolicyID:       policy.PolicyID,
				ConstraintType: "budget.daily_budget",
				Description:    fmt.Sprintf("Daily budget $%.2f would be exceeded", bc.DailyBudgetUSD),
				Severity:       "critical",
				Timestamp:      time.Now().UTC(),
			})
		}

		// Check monthly budget
		if bc.MonthlyBudgetUSD > 0 && e.tracker.monthlySpendUSD+estimatedCostUSD > bc.MonthlyBudgetUSD {
			violations = append(violations, Violation{
				PolicyID:       policy.PolicyID,
				ConstraintType: "budget.monthly_budget",
				Description:    fmt.Sprintf("Monthly budget $%.2f would be exceeded", bc.MonthlyBudgetUSD),
				Severity:       "critical",
				Timestamp:      time.Now().UTC(),
			})
		}

		e.tracker.mu.Unlock()
	}

	// Rate limits
	if policy.RateLimits != nil {
		rl := policy.RateLimits

		// Check concurrent requests
		if rl.ConcurrentRequests > 0 {
			current := atomic.LoadInt64(&e.tracker.concurrent)
			if int(current) >= rl.ConcurrentRequests {
				violations = append(violations, Violation{
					PolicyID:       policy.PolicyID,
					ConstraintType: "rate.concurrent_requests",
					Description:    fmt.Sprintf("Concurrent limit %d reached", rl.ConcurrentRequests),
					Severity:       "high",
					Timestamp:      time.Now().UTC(),
				})
			}
		}

		// Check RPM
		if rl.RequestsPerMinute > 0 {
			e.tracker.mu.Lock()
			e.maybeResetCounters()
			if int(e.tracker.minuteRequests) >= rl.RequestsPerMinute {
				violations = append(violations, Violation{
					PolicyID:       policy.PolicyID,
					ConstraintType: "rate.requests_per_minute",
					Description:    fmt.Sprintf("RPM limit %d reached", rl.RequestsPerMinute),
					Severity:       "high",
					Timestamp:      time.Now().UTC(),
				})
			}
			e.tracker.mu.Unlock()
		}
	}

	return violations
}

// RecordUsage records usage for tracking.
func (e *Enforcer) RecordUsage(costUSD float64, tokens int) {
	e.tracker.mu.Lock()
	defer e.tracker.mu.Unlock()

	e.maybeResetCounters()

	e.tracker.dailySpendUSD += costUSD
	e.tracker.monthlySpendUSD += costUSD
	e.tracker.minuteRequests++
	e.tracker.hourRequests++
	e.tracker.minuteTokens += int64(tokens)
}

// AcquireConcurrent acquires a concurrent request slot.
func (e *Enforcer) AcquireConcurrent() {
	atomic.AddInt64(&e.tracker.concurrent, 1)
}

// ReleaseConcurrent releases a concurrent request slot.
func (e *Enforcer) ReleaseConcurrent() {
	atomic.AddInt64(&e.tracker.concurrent, -1)
}

// maybeResetCounters resets time-based counters if needed.
func (e *Enforcer) maybeResetCounters() {
	now := time.Now().UTC()

	if now.After(e.tracker.minuteReset) {
		e.tracker.minuteRequests = 0
		e.tracker.minuteTokens = 0
		e.tracker.minuteReset = now.Add(time.Minute)
	}

	if now.After(e.tracker.hourReset) {
		e.tracker.hourRequests = 0
		e.tracker.hourReset = now.Add(time.Hour)
	}

	if now.After(e.tracker.dailyReset) {
		e.tracker.dailySpendUSD = 0
		e.tracker.dailyReset = now.Truncate(24 * time.Hour).Add(24 * time.Hour)
	}

	if now.After(e.tracker.monthlyReset) {
		e.tracker.monthlySpendUSD = 0
		e.tracker.monthlyReset = time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	}
}

// notifyHandlers notifies violation handlers.
func (e *Enforcer) notifyHandlers(ctx context.Context, result CheckResult) {
	for _, h := range e.handlers {
		h(ctx, result)
	}
}

// GetUsageStats returns current usage statistics.
func (e *Enforcer) GetUsageStats() map[string]interface{} {
	e.tracker.mu.Lock()
	defer e.tracker.mu.Unlock()

	e.maybeResetCounters()

	return map[string]interface{}{
		"daily_spend_usd":   e.tracker.dailySpendUSD,
		"monthly_spend_usd": e.tracker.monthlySpendUSD,
		"minute_requests":   e.tracker.minuteRequests,
		"hour_requests":     e.tracker.hourRequests,
		"concurrent":        atomic.LoadInt64(&e.tracker.concurrent),
	}
}
