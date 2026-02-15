package finance

// Cost represents the consumption associated with an action.
// It can be financial (Money), computational (Tokens), or abstract (Risk).
type Cost struct {
	Money    Money `json:"money,omitempty"`
	Tokens   int64 `json:"tokens,omitempty"`   // LLM tokens
	Compute  int64 `json:"compute,omitempty"`  // CPU ms
	Requests int64 `json:"requests,omitempty"` // API calls
}

// Add sums two Costs.
func (c Cost) Add(other Cost) (Cost, error) {
	m, err := c.Money.Add(other.Money)
	if err != nil && !c.Money.IsZero() && !other.Money.IsZero() {
		return Cost{}, err
	}
	// Handle zero money case aggressively for now
	if c.Money.IsZero() {
		m = other.Money
	}

	return Cost{
		Money:    m,
		Tokens:   c.Tokens + other.Tokens,
		Compute:  c.Compute + other.Compute,
		Requests: c.Requests + other.Requests,
	}, nil
}

// Estimator interface for predicting costs.
type Estimator interface {
	Estimate(input interface{}) (Cost, error)
}
