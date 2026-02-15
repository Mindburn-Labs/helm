package console

import (
	"net/http"
)

// OIDC Skeleton Implementation

func (s *Server) handleAuthLoginAPI(w http.ResponseWriter, r *http.Request) {
	// RED TEAM: Removed mock login flow.
	// Production must require real OIDC provider.
	http.Error(w, "OIDC Provider Not Configured (Mock Removed)", http.StatusServiceUnavailable)
}

func (s *Server) handleAuthCallbackAPI(w http.ResponseWriter, r *http.Request) {
	// RED TEAM: Removed mock callback.
	http.Error(w, "OIDC Callback Not Implemented (Mock Removed)", http.StatusServiceUnavailable)
}

// Ensure these are registered in server.go
