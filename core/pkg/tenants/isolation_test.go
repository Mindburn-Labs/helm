package tenants

import (
	"testing"
)

func TestIsolatedAccess(t *testing.T) {
	c := NewIsolationChecker()
	c.RegisterResource("t1", "db-1")
	c.RegisterResource("t1", "db-2")

	receipt := c.CheckAccess("t1", []string{"db-1", "db-2"})
	if !receipt.Isolated {
		t.Fatalf("expected isolated, got violations: %v", receipt.Violations)
	}
	if receipt.ChecksPassed != 2 {
		t.Fatalf("expected 2 passed, got %d", receipt.ChecksPassed)
	}
}

func TestCrossTenantViolation(t *testing.T) {
	c := NewIsolationChecker()
	c.RegisterResource("t1", "db-1")
	c.RegisterResource("t2", "db-2")

	receipt := c.CheckAccess("t1", []string{"db-1", "db-2"})
	if receipt.Isolated {
		t.Fatal("expected cross-tenant violation")
	}
	if receipt.ChecksFailed != 1 {
		t.Fatalf("expected 1 failure, got %d", receipt.ChecksFailed)
	}
	if len(receipt.Violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(receipt.Violations))
	}
}

func TestUnregisteredResourceAllowed(t *testing.T) {
	c := NewIsolationChecker()
	c.RegisterResource("t1", "db-1")

	// "new-db" not registered to any tenant, should pass
	receipt := c.CheckAccess("t1", []string{"db-1", "new-db"})
	if !receipt.Isolated {
		t.Fatal("unregistered resource should not cause violation")
	}
}

func TestVerifyIsolationClean(t *testing.T) {
	c := NewIsolationChecker()
	c.RegisterResource("t1", "db-1")
	c.RegisterResource("t2", "db-2")

	ok, _ := c.VerifyIsolation()
	if !ok {
		t.Fatal("expected clean isolation")
	}
}

func TestVerifyIsolationConflict(t *testing.T) {
	c := NewIsolationChecker()
	c.RegisterResource("t1", "shared-db")
	c.RegisterResource("t2", "shared-db")

	ok, violations := c.VerifyIsolation()
	if ok {
		t.Fatal("expected conflict for shared resource")
	}
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestIsolationReceiptHash(t *testing.T) {
	c := NewIsolationChecker()
	c.RegisterResource("t1", "db-1")

	receipt := c.CheckAccess("t1", []string{"db-1"})
	if receipt.ContentHash == "" {
		t.Fatal("expected content hash")
	}
}

func TestMultipleTenants(t *testing.T) {
	c := NewIsolationChecker()
	c.RegisterResource("t1", "a")
	c.RegisterResource("t2", "b")
	c.RegisterResource("t3", "c")

	r1 := c.CheckAccess("t1", []string{"a"})
	r2 := c.CheckAccess("t2", []string{"b"})
	r3 := c.CheckAccess("t3", []string{"c"})

	if !r1.Isolated || !r2.Isolated || !r3.Isolated {
		t.Fatal("all tenants accessing own resources should be isolated")
	}

	cross := c.CheckAccess("t1", []string{"b"}) // t1 accessing t2's resource
	if cross.Isolated {
		t.Fatal("expected violation")
	}
}
