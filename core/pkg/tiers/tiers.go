// Package tiers defines product tier definitions for HELM.
// Tiers map to limits, features, and pricing.
package tiers

// TierID identifies a product tier.
type TierID string

const (
	TierFree       TierID = "free"
	TierPro        TierID = "pro"
	TierEnterprise TierID = "enterprise"
)

// Limits defines resource limits for a tier.
type Limits struct {
	DailyExecutions   int64 // -1 = unlimited
	MonthlyTokens     int64 // -1 = unlimited
	StorageGB         int64 // -1 = unlimited
	ConcurrentAgents  int   // -1 = unlimited
	RetentionDays     int   // How long to keep receipts/evidence
	MaxToolsPerIntent int   // Max tools per single intent
}

// Tier represents a product tier with limits, features, and pricing.
type Tier struct {
	ID            TierID
	Name          string
	Description   string
	Limits        Limits
	Features      []string
	PricePerMonth int64 // cents, -1 = custom pricing
}

// All available tiers
var (
	Free = Tier{
		ID:          TierFree,
		Name:        "Free",
		Description: "For individuals and small projects",
		Limits: Limits{
			DailyExecutions:   100,
			MonthlyTokens:     100_000,
			StorageGB:         1,
			ConcurrentAgents:  1,
			RetentionDays:     30,
			MaxToolsPerIntent: 5,
		},
		Features:      []string{"basic_governance", "basic_receipts"},
		PricePerMonth: 0,
	}

	Pro = Tier{
		ID:          TierPro,
		Name:        "Pro",
		Description: "For teams and production workloads",
		Limits: Limits{
			DailyExecutions:   10_000,
			MonthlyTokens:     10_000_000,
			StorageGB:         100,
			ConcurrentAgents:  10,
			RetentionDays:     365,
			MaxToolsPerIntent: 20,
		},
		Features: []string{
			"basic_governance",
			"basic_receipts",
			"advanced_receipts",
			"api_access",
			"priority_support",
			"custom_policies",
		},
		PricePerMonth: 9900, // $99
	}

	Enterprise = Tier{
		ID:          TierEnterprise,
		Name:        "Enterprise",
		Description: "For large organizations with compliance needs",
		Limits: Limits{
			DailyExecutions:   -1, // unlimited
			MonthlyTokens:     -1,
			StorageGB:         -1,
			ConcurrentAgents:  -1,
			RetentionDays:     -1, // unlimited
			MaxToolsPerIntent: -1,
		},
		Features: []string{
			"all",
			"hsm",
			"sso",
			"sla",
			"dedicated_support",
			"custom_integrations",
			"audit_exports",
			"compliance_reports",
		},
		PricePerMonth: -1, // custom
	}

	// AllTiers contains all available tiers
	AllTiers = map[TierID]Tier{
		TierFree:       Free,
		TierPro:        Pro,
		TierEnterprise: Enterprise,
	}
)

// Get returns a tier by ID, or nil if not found.
func Get(id TierID) *Tier {
	tier, ok := AllTiers[id]
	if !ok {
		return nil
	}
	return &tier
}

// HasFeature checks if a tier has a specific feature.
func (t *Tier) HasFeature(feature string) bool {
	for _, f := range t.Features {
		if f == feature || f == "all" {
			return true
		}
	}
	return false
}

// IsUnlimited checks if a limit is unlimited (-1).
func IsUnlimited(limit int64) bool {
	return limit < 0
}
