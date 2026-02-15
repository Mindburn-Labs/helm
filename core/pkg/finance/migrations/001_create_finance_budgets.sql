CREATE TABLE IF NOT EXISTS finance_budgets (
    id TEXT PRIMARY KEY,
    resource_type TEXT NOT NULL,
    budget_limit BIGINT NOT NULL,
    consumed BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_finance_budgets_resource_type ON finance_budgets (resource_type);
