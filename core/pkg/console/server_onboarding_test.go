package console

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestServer() *Server {
	return &Server{
		cache:          make(map[string][]byte),
		pendingSignups: make(map[string]*pendingSignup),
		errorBudget:    100.0,
		systemStatus:   "HEALTHY",
		intents:        make(map[string]*operatorIntent),
		approvals:      make(map[string]*operatorApproval),
		operatorRuns:   make(map[string]*operatorRunState),
	}
}

func TestSignupAPI(t *testing.T) {
	srv := newTestServer()

	t.Run("successful signup", func(t *testing.T) {
		body := `{"email":"test@example.com","password":"securepass123"}`
		req := httptest.NewRequest(http.MethodPost, "/api/signup", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		srv.handleSignupAPI(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
		}

		var resp signupResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if resp.TenantID == "" {
			t.Fatal("tenantId should not be empty")
		}
		if resp.Message == "" {
			t.Fatal("message should not be empty")
		}
	})

	t.Run("missing email", func(t *testing.T) {
		body := `{"email":"","password":"securepass123"}`
		req := httptest.NewRequest(http.MethodPost, "/api/signup", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		srv.handleSignupAPI(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("password too short", func(t *testing.T) {
		body := `{"email":"short@test.com","password":"abc"}`
		req := httptest.NewRequest(http.MethodPost, "/api/signup", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		srv.handleSignupAPI(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("wrong method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/signup", nil)
		w := httptest.NewRecorder()

		srv.handleSignupAPI(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Fatalf("expected 405, got %d", w.Code)
		}
	})
}

func TestOnboardingVerifyAPI(t *testing.T) {
	srv := newTestServer()

	// First, create a signup
	srv.pendingSignups["verify@test.com"] = &pendingSignup{
		Email:    "verify@test.com",
		Password: "securepass",
		Code:     "123456",
		TenantID: "tenant-abc123",
		Verified: false,
	}

	t.Run("successful verify", func(t *testing.T) {
		body := `{"email":"verify@test.com","code":"123456"}`
		req := httptest.NewRequest(http.MethodPost, "/api/onboarding/verify", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		srv.handleOnboardingVerifyAPI(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}

		var resp onboardingVerifyResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if resp.TenantID != "tenant-abc123" {
			t.Fatalf("expected tenant-abc123, got %s", resp.TenantID)
		}
		if resp.APIKey == "" {
			t.Fatal("apiKey should not be empty")
		}
	})

	t.Run("wrong code", func(t *testing.T) {
		srv.pendingSignups["wrong@test.com"] = &pendingSignup{
			Email: "wrong@test.com", Code: "111111", TenantID: "t2", Verified: false,
		}
		body := `{"email":"wrong@test.com","code":"999999"}`
		req := httptest.NewRequest(http.MethodPost, "/api/onboarding/verify", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		srv.handleOnboardingVerifyAPI(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("unknown email", func(t *testing.T) {
		body := `{"email":"nobody@test.com","code":"123456"}`
		req := httptest.NewRequest(http.MethodPost, "/api/onboarding/verify", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		srv.handleOnboardingVerifyAPI(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})
}

func TestResendVerificationAPI(t *testing.T) {
	srv := newTestServer()

	srv.pendingSignups["resend@test.com"] = &pendingSignup{
		Email: "resend@test.com", Code: "111111", TenantID: "t3", Verified: false,
	}

	t.Run("successful resend", func(t *testing.T) {
		oldCode := srv.pendingSignups["resend@test.com"].Code

		body := `{"email":"resend@test.com"}`
		req := httptest.NewRequest(http.MethodPost, "/api/resend-verification", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		srv.handleResendVerificationAPI(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}

		newCode := srv.pendingSignups["resend@test.com"].Code
		if newCode == oldCode {
			t.Log("code may randomly match — acceptable in tests")
		}
	})

	t.Run("unknown email silently succeeds", func(t *testing.T) {
		body := `{"email":"ghost@test.com"}`
		req := httptest.NewRequest(http.MethodPost, "/api/resend-verification", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		srv.handleResendVerificationAPI(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 (silent), got %d", w.Code)
		}
	})
}

func TestOnboardingFullFlow(t *testing.T) {
	srv := newTestServer()

	// Step 1: Signup
	signupBody := `{"email":"flow@test.com","password":"mypassword123"}`
	req1 := httptest.NewRequest(http.MethodPost, "/api/signup", bytes.NewBufferString(signupBody))
	w1 := httptest.NewRecorder()
	srv.handleSignupAPI(w1, req1)

	if w1.Code != http.StatusCreated {
		t.Fatalf("signup: expected 201, got %d: %s", w1.Code, w1.Body.String())
	}

	var signupResp signupResponse
	_ = json.NewDecoder(w1.Body).Decode(&signupResp)

	// Get the verification code from internal state
	srv.onboardingMu.RLock()
	code := srv.pendingSignups["flow@test.com"].Code
	srv.onboardingMu.RUnlock()

	// Step 2: Verify with correct code
	verifyBody := `{"email":"flow@test.com","code":"` + code + `"}`
	req2 := httptest.NewRequest(http.MethodPost, "/api/onboarding/verify", bytes.NewBufferString(verifyBody))
	w2 := httptest.NewRecorder()
	srv.handleOnboardingVerifyAPI(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("verify: expected 200, got %d: %s", w2.Code, w2.Body.String())
	}

	var verifyResp onboardingVerifyResponse
	_ = json.NewDecoder(w2.Body).Decode(&verifyResp)

	if verifyResp.TenantID == "" {
		t.Fatal("tenantId should not be empty after verification")
	}
	if verifyResp.APIKey == "" {
		t.Fatal("apiKey should not be empty after verification")
	}
}
