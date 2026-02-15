package console

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRecordCost(t *testing.T) {
	mm := NewMetricsManager()

	mm.RecordCost(CostRecord{
		TenantID:  "t1",
		Model:     "gpt-4",
		Operation: "eval",
		Tokens:    1500,
		CostCents: 4.50,
		Timestamp: time.Now(),
	})

	mm.RecordCost(CostRecord{
		TenantID:  "t1",
		Model:     "gpt-4",
		Operation: "tool_call",
		Tokens:    500,
		CostCents: 1.50,
	})

	summary := mm.GetCostSummary("t1")
	require.Equal(t, "t1", summary.TenantID)
	require.InDelta(t, 6.0, summary.TotalCents, 0.01)
	require.Equal(t, int64(2000), summary.TotalTokens)
	require.Equal(t, 2, summary.TotalOps)
}

func TestCostSummaryByModel(t *testing.T) {
	mm := NewMetricsManager()

	mm.RecordCost(CostRecord{TenantID: "t1", Model: "gpt-4", Operation: "eval", Tokens: 1000, CostCents: 3.0})
	mm.RecordCost(CostRecord{TenantID: "t1", Model: "claude-3", Operation: "eval", Tokens: 800, CostCents: 2.4})
	mm.RecordCost(CostRecord{TenantID: "t1", Model: "gpt-4", Operation: "eval", Tokens: 500, CostCents: 1.5})

	summary := mm.GetCostSummary("t1")
	require.Len(t, summary.ByModel, 2)

	// Check aggregation
	modelCosts := make(map[string]float64)
	for _, b := range summary.ByModel {
		modelCosts[b.Key] = b.CostCents
	}
	require.InDelta(t, 4.5, modelCosts["gpt-4"], 0.01)
	require.InDelta(t, 2.4, modelCosts["claude-3"], 0.01)
}

func TestCostSummaryByOperation(t *testing.T) {
	mm := NewMetricsManager()

	mm.RecordCost(CostRecord{TenantID: "t1", Model: "gpt-4", Operation: "eval", Tokens: 1000, CostCents: 3.0})
	mm.RecordCost(CostRecord{TenantID: "t1", Model: "gpt-4", Operation: "tool_call", Tokens: 200, CostCents: 0.6})
	mm.RecordCost(CostRecord{TenantID: "t1", Model: "gpt-4", Operation: "eval", Tokens: 500, CostCents: 1.5})

	summary := mm.GetCostSummary("t1")
	require.Len(t, summary.ByOperation, 2)

	opCosts := make(map[string]float64)
	for _, b := range summary.ByOperation {
		opCosts[b.Key] = b.CostCents
	}
	require.InDelta(t, 4.5, opCosts["eval"], 0.01)
	require.InDelta(t, 0.6, opCosts["tool_call"], 0.01)
}

func TestCostSummaryTenantFilter(t *testing.T) {
	mm := NewMetricsManager()

	mm.RecordCost(CostRecord{TenantID: "t1", Model: "gpt-4", CostCents: 3.0})
	mm.RecordCost(CostRecord{TenantID: "t2", Model: "gpt-4", CostCents: 5.0})

	// Filtered to t1
	s1 := mm.GetCostSummary("t1")
	require.InDelta(t, 3.0, s1.TotalCents, 0.01)

	// Filtered to t2
	s2 := mm.GetCostSummary("t2")
	require.InDelta(t, 5.0, s2.TotalCents, 0.01)

	// Empty tenant = all
	sAll := mm.GetCostSummary("")
	require.InDelta(t, 8.0, sAll.TotalCents, 0.01)
}

func TestCheckBudgetAlerts(t *testing.T) {
	mm := NewMetricsManager()

	mm.RecordCost(CostRecord{TenantID: "t1", CostCents: 85.0})
	mm.RecordCost(CostRecord{TenantID: "t2", CostCents: 95.0})
	mm.RecordCost(CostRecord{TenantID: "t3", CostCents: 105.0})
	mm.RecordCost(CostRecord{TenantID: "t4", CostCents: 50.0})

	budgets := map[string]float64{
		"t1": 100.0,
		"t2": 100.0,
		"t3": 100.0,
		"t4": 100.0,
	}

	alerts := mm.CheckBudgetAlerts(budgets)

	alertMap := make(map[string]string)
	for _, a := range alerts {
		alertMap[a.TenantID] = a.AlertType
	}

	require.Equal(t, "warning", alertMap["t1"])  // 85%
	require.Equal(t, "critical", alertMap["t2"]) // 95%
	require.Equal(t, "exceeded", alertMap["t3"]) // 105%
	require.Empty(t, alertMap["t4"])             // 50% — no alert
}

func TestCheckBudgetAlerts_NoBudget(t *testing.T) {
	mm := NewMetricsManager()
	mm.RecordCost(CostRecord{TenantID: "t1", CostCents: 85.0})

	// No budgets defined → no alerts
	alerts := mm.CheckBudgetAlerts(map[string]float64{})
	require.Empty(t, alerts)
}

func TestCostRecordRingBuffer(t *testing.T) {
	mm := NewMetricsManager()

	// Insert just over the ring buffer cap
	for i := 0; i < 10005; i++ {
		mm.RecordCost(CostRecord{TenantID: "t1", CostCents: 0.01})
	}

	mm.mu.RLock()
	require.LessOrEqual(t, len(mm.costRecords), 10000)
	mm.mu.RUnlock()
}
