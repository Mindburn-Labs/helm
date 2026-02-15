package finance

import (
	"testing"
)

func TestMoney_Add(t *testing.T) {
	m1 := NewMoney(100, "USD")
	m2 := NewMoney(50, "USD")

	sum, err := m1.Add(m2)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if sum.AmountMinor != 150 {
		t.Errorf("Expected 150, got %d", sum.AmountMinor)
	}
}

func TestMoney_Add_Mismatch(t *testing.T) {
	m1 := NewMoney(100, "USD")
	m2 := NewMoney(50, "EUR")

	_, err := m1.Add(m2)
	if err == nil {
		t.Error("Expected currency mismatch error")
	}
}

func TestBudget_Enforcement(t *testing.T) {
	tracker := NewInMemoryTracker()
	b := Budget{
		ID:           "test-budget",
		ResourceType: "USD",
		Limit:        1000, // $10.00
		Consumed:     0,
	}
	tracker.budgets[b.ID] = &b

	cost1 := Cost{Money: NewMoney(500, "USD")} // $5.00
	if err := tracker.Consume(b.ID, cost1); err != nil {
		t.Errorf("First consume failed: %v", err)
	}

	cost2 := Cost{Money: NewMoney(600, "USD")} // $6.00 (Total $11.00 > Limit)
	if err := tracker.Consume(b.ID, cost2); err == nil {
		t.Error("Expected budget exceeded error")
	}

	cost3 := Cost{Money: NewMoney(100, "USD")} // $1.00 (Total $6.00 < Limit)
	if err := tracker.Consume(b.ID, cost3); err != nil {
		t.Errorf("Third consume failed: %v", err)
	}

	if tracker.budgets[b.ID].Consumed != 600 {
		t.Errorf("Expected consumed 600, got %d", tracker.budgets[b.ID].Consumed)
	}
}
