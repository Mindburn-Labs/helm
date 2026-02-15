package arc

import (
	"context"

	"golang.org/x/time/rate"
)

// BaseConnector provides common functionality for source connectors.
type BaseConnector struct {
	id         string
	trustClass TrustClass
	version    string
	limiter    *rate.Limiter
}

// NewBaseConnector creates a new BaseConnector with rate limiting.
func NewBaseConnector(id string, trustClass TrustClass, version string, r rate.Limit, b int) *BaseConnector {
	return &BaseConnector{
		id:         id,
		trustClass: trustClass,
		version:    version,
		limiter:    rate.NewLimiter(r, b),
	}
}

func (c *BaseConnector) ID() string {
	return c.id
}

func (c *BaseConnector) TrustClass() TrustClass {
	return c.trustClass
}

// Version returns the version of this connector.
func (c *BaseConnector) Version() string {
	return c.version
}

// Wait blocks until the rate limiter allows an event.
func (c *BaseConnector) Wait(ctx context.Context) error {
	return c.limiter.Wait(ctx)
}
