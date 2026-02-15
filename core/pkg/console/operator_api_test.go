package console

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newTestOperatorServer creates a Server with in-memory stores for operator tests.
func newTestOperatorServer() *Server {
	store := &inMemoryReceiptStore{
		receipts: nil,
	}
	return &Server{
		cache:          make(map[string][]byte),
		pendingSignups: make(map[string]*pendingSignup),
		errorBudget:    100.0,
		systemStatus:   "HEALTHY",
		receiptStore:   store,
		intents:        make(map[string]*operatorIntent),
		approvals:      make(map[string]*operatorApproval),
		operatorRuns:   make(map[string]*operatorRunState),
	}
}

func TestOperatorCreateIntent(t *testing.T) {
	srv := newTestOperatorServer()

	t.Run("creates intent and returns 201", func(t *testing.T) {
		body := `{"type":"deploy","description":"Deploy new API version","params":{"version":"2.0"}}`
		req := httptest.NewRequest(http.MethodPost, "/api/intents", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		srv.handleCreateIntent(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
		}

		var intent operatorIntent
		if err := json.NewDecoder(w.Body).Decode(&intent); err != nil {
			t.Fatalf("decode: %v", err)
		}

		if intent.IntentID == "" {
			t.Fatal("intent_id should not be empty")
		}
		if intent.Type != "deploy" {
			t.Fatalf("expected type=deploy, got %s", intent.Type)
		}
		if intent.Status != IntentDraft {
			t.Fatalf("expected status=draft, got %s", intent.Status)
		}
		if intent.Description != "Deploy new API version" {
			t.Fatalf("expected description='Deploy new API version', got %s", intent.Description)
		}
	})

	t.Run("rejects empty type", func(t *testing.T) {
		body := `{"type":"","description":"test"}`
		req := httptest.NewRequest(http.MethodPost, "/api/intents", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		srv.handleCreateIntent(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("rejects wrong method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/intents", nil)
		w := httptest.NewRecorder()

		srv.handleCreateIntent(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Fatalf("expected 405, got %d", w.Code)
		}
	})
}

func TestOperatorListIntents(t *testing.T) {
	srv := newTestOperatorServer()

	// Seed some intents
	srv.intents["int-1"] = &operatorIntent{
		IntentID: "int-1", Type: "deploy", Description: "Deploy v1",
		Status: IntentDraft, CreatedAt: "2026-02-08T00:00:00Z", UpdatedAt: "2026-02-08T00:00:00Z",
	}
	srv.intents["int-2"] = &operatorIntent{
		IntentID: "int-2", Type: "scale", Description: "Scale up",
		Status: IntentSubmitted, CreatedAt: "2026-02-08T01:00:00Z", UpdatedAt: "2026-02-08T01:00:00Z",
	}

	t.Run("lists all intents", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/intents", nil)
		w := httptest.NewRecorder()

		srv.handleListIntents(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}

		var resp intentsListResponse
		_ = json.NewDecoder(w.Body).Decode(&resp)
		if resp.Total != 2 {
			t.Fatalf("expected total=2, got %d", resp.Total)
		}
	})

	t.Run("filters by status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/intents?status=submitted", nil)
		w := httptest.NewRecorder()

		srv.handleListIntents(w, req)

		var resp intentsListResponse
		_ = json.NewDecoder(w.Body).Decode(&resp)
		if resp.Total != 1 {
			t.Fatalf("expected 1 submitted intent, got %d", resp.Total)
		}
		if resp.Intents[0].IntentID != "int-2" {
			t.Fatalf("expected int-2, got %s", resp.Intents[0].IntentID)
		}
	})
}

func TestOperatorGetIntent(t *testing.T) {
	srv := newTestOperatorServer()
	srv.intents["int-1"] = &operatorIntent{
		IntentID: "int-1", Type: "deploy", Description: "Deploy v1",
		Status: IntentDraft, CreatedAt: "2026-02-08T00:00:00Z", UpdatedAt: "2026-02-08T00:00:00Z",
	}

	t.Run("returns intent by ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/intents/int-1", nil)
		w := httptest.NewRecorder()

		srv.handleGetIntent(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}

		var intent operatorIntent
		_ = json.NewDecoder(w.Body).Decode(&intent)
		if intent.IntentID != "int-1" {
			t.Fatalf("expected int-1, got %s", intent.IntentID)
		}
	})

	t.Run("returns 404 for unknown", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/intents/nonexistent", nil)
		w := httptest.NewRecorder()

		srv.handleGetIntent(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})
}

func TestOperatorPlanIntent(t *testing.T) {
	srv := newTestOperatorServer()
	srv.intents["int-1"] = &operatorIntent{
		IntentID: "int-1", Type: "deploy", Description: "Deploy v1",
		Status: IntentDraft, CreatedAt: "2026-02-08T00:00:00Z", UpdatedAt: "2026-02-08T00:00:00Z",
	}

	t.Run("generates plan for draft intent", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/intents/int-1/plan", nil)
		w := httptest.NewRecorder()

		srv.handlePlanIntent(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}

		var plan operatorPlan
		_ = json.NewDecoder(w.Body).Decode(&plan)

		if plan.IntentID != "int-1" {
			t.Fatalf("expected intent_id=int-1, got %s", plan.IntentID)
		}
		if plan.Hash == "" {
			t.Fatal("plan hash should not be empty")
		}
		if !plan.DryRunOK {
			t.Fatal("dry_run_ok should be true")
		}
		if len(plan.Steps) == 0 {
			t.Fatal("plan should have steps")
		}
		if plan.RiskLevel != "high" {
			t.Fatalf("expected risk_level=high for deploy, got %s", plan.RiskLevel)
		}
		if plan.ReceiptID == "" {
			t.Fatal("receipt_id should not be empty")
		}

		// Verify intent status changed to planned
		if srv.intents["int-1"].Status != IntentPlanned {
			t.Fatalf("expected intent status=planned, got %s", srv.intents["int-1"].Status)
		}
	})

	t.Run("rejects plan for submitted intent", func(t *testing.T) {
		srv.intents["int-2"] = &operatorIntent{
			IntentID: "int-2", Type: "scale", Description: "Scale up",
			Status: IntentSubmitted, CreatedAt: "2026-02-08T00:00:00Z", UpdatedAt: "2026-02-08T00:00:00Z",
		}

		req := httptest.NewRequest(http.MethodPost, "/api/intents/int-2/plan", nil)
		w := httptest.NewRecorder()

		srv.handlePlanIntent(w, req)

		if w.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d", w.Code)
		}
	})
}

func TestOperatorSubmitIntent(t *testing.T) {
	srv := newTestOperatorServer()
	srv.intents["int-1"] = &operatorIntent{
		IntentID: "int-1", Type: "deploy", Description: "Deploy v1",
		Status: IntentPlanned, PlanHash: "sha256:abc", CreatedAt: "2026-02-08T00:00:00Z", UpdatedAt: "2026-02-08T00:00:00Z",
	}

	t.Run("submits planned intent", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/intents/int-1/submit", nil)
		w := httptest.NewRecorder()

		srv.handleSubmitIntent(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}

		var intent operatorIntent
		_ = json.NewDecoder(w.Body).Decode(&intent)
		if intent.Status != IntentSubmitted {
			t.Fatalf("expected status=submitted, got %s", intent.Status)
		}
	})

	t.Run("rejects submit for draft intent", func(t *testing.T) {
		srv.intents["int-2"] = &operatorIntent{
			IntentID: "int-2", Type: "scale", Description: "Scale up",
			Status: IntentDraft, CreatedAt: "2026-02-08T00:00:00Z", UpdatedAt: "2026-02-08T00:00:00Z",
		}

		req := httptest.NewRequest(http.MethodPost, "/api/intents/int-2/submit", nil)
		w := httptest.NewRecorder()

		srv.handleSubmitIntent(w, req)

		if w.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d", w.Code)
		}
	})
}

func TestOperatorApprove(t *testing.T) {
	srv := newTestOperatorServer()
	srv.intents["int-1"] = &operatorIntent{
		IntentID: "int-1", Type: "deploy", Description: "Deploy v1",
		Status: IntentSubmitted, PlanHash: "sha256:abc", CreatedAt: "2026-02-08T00:00:00Z", UpdatedAt: "2026-02-08T00:00:00Z",
	}

	t.Run("approves submitted intent", func(t *testing.T) {
		body := `{"decision":"approve","reason":"Looks good, plan verified","decided_by":"admin"}`
		req := httptest.NewRequest(http.MethodPost, "/api/approvals/int-1/approve", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		srv.handleApproveIntent(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}

		var approval operatorApproval
		_ = json.NewDecoder(w.Body).Decode(&approval)

		if approval.Decision != DecisionApprove {
			t.Fatalf("expected decision=approve, got %s", approval.Decision)
		}
		if approval.IntentID != "int-1" {
			t.Fatalf("expected intent_id=int-1, got %s", approval.IntentID)
		}
		if approval.ReceiptID == "" {
			t.Fatal("receipt_id should not be empty")
		}

		// Verify intent status changed
		if srv.intents["int-1"].Status != IntentApproved {
			t.Fatalf("expected intent status=approved, got %s", srv.intents["int-1"].Status)
		}
	})

	t.Run("rejects approval without reason", func(t *testing.T) {
		srv.intents["int-2"] = &operatorIntent{
			IntentID: "int-2", Type: "scale", Description: "Scale",
			Status: IntentSubmitted, CreatedAt: "2026-02-08T00:00:00Z", UpdatedAt: "2026-02-08T00:00:00Z",
		}

		body := `{"decision":"approve","reason":""}`
		req := httptest.NewRequest(http.MethodPost, "/api/approvals/int-2/approve", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		srv.handleApproveIntent(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("rejects approval for non-submitted intent", func(t *testing.T) {
		srv.intents["int-3"] = &operatorIntent{
			IntentID: "int-3", Type: "scale", Description: "Scale",
			Status: IntentDraft, CreatedAt: "2026-02-08T00:00:00Z", UpdatedAt: "2026-02-08T00:00:00Z",
		}

		body := `{"decision":"approve","reason":"test"}`
		req := httptest.NewRequest(http.MethodPost, "/api/approvals/int-3/approve", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		srv.handleApproveIntent(w, req)

		if w.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d", w.Code)
		}
	})
}

func TestOperatorCreateRun(t *testing.T) {
	srv := newTestOperatorServer()
	srv.intents["int-1"] = &operatorIntent{
		IntentID: "int-1", Type: "deploy", Description: "Deploy v1",
		Status: IntentApproved, PlanHash: "sha256:abc", CreatedAt: "2026-02-08T00:00:00Z", UpdatedAt: "2026-02-08T00:00:00Z",
	}

	t.Run("creates run from approved intent", func(t *testing.T) {
		body := `{"intent_id":"int-1"}`
		req := httptest.NewRequest(http.MethodPost, "/api/runs", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		srv.handleCreateRun(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
		}

		var run operatorRunState
		_ = json.NewDecoder(w.Body).Decode(&run)

		if run.RunID == "" {
			t.Fatal("run_id should not be empty")
		}
		if run.IntentID != "int-1" {
			t.Fatalf("expected intent_id=int-1, got %s", run.IntentID)
		}
		if run.State != RunStateRunning {
			t.Fatalf("expected state=running, got %s", run.State)
		}
		if run.PlanHash != "sha256:abc" {
			t.Fatalf("expected plan_hash=sha256:abc, got %s", run.PlanHash)
		}
		if run.ReceiptID == "" {
			t.Fatal("receipt_id should not be empty")
		}
	})

	t.Run("rejects run for non-approved intent", func(t *testing.T) {
		srv.intents["int-2"] = &operatorIntent{
			IntentID: "int-2", Type: "scale", Description: "Scale",
			Status: IntentDraft, CreatedAt: "2026-02-08T00:00:00Z", UpdatedAt: "2026-02-08T00:00:00Z",
		}

		body := `{"intent_id":"int-2"}`
		req := httptest.NewRequest(http.MethodPost, "/api/runs", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		srv.handleCreateRun(w, req)

		if w.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d", w.Code)
		}
	})
}

func TestOperatorRunControl(t *testing.T) {
	srv := newTestOperatorServer()

	t.Run("pauses a running run", func(t *testing.T) {
		srv.operatorRuns["run-1"] = &operatorRunState{
			RunID: "run-1", IntentID: "int-1", State: RunStateRunning,
			CreatedAt: "2026-02-08T00:00:00Z", UpdatedAt: "2026-02-08T00:00:00Z",
		}

		body := `{"action":"pause","reason":"Investigating anomaly"}`
		req := httptest.NewRequest(http.MethodPost, "/api/runs/run-1/control", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		srv.handleRunControl(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}

		var run operatorRunState
		_ = json.NewDecoder(w.Body).Decode(&run)
		if run.State != RunStatePaused {
			t.Fatalf("expected state=paused, got %s", run.State)
		}
	})

	t.Run("resumes a paused run", func(t *testing.T) {
		body := `{"action":"resume"}`
		req := httptest.NewRequest(http.MethodPost, "/api/runs/run-1/control", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		srv.handleRunControl(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}

		var run operatorRunState
		_ = json.NewDecoder(w.Body).Decode(&run)
		if run.State != RunStateRunning {
			t.Fatalf("expected state=running, got %s", run.State)
		}
	})

	t.Run("cancels a running run", func(t *testing.T) {
		body := `{"action":"cancel","reason":"No longer needed"}`
		req := httptest.NewRequest(http.MethodPost, "/api/runs/run-1/control", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		srv.handleRunControl(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}

		var run operatorRunState
		_ = json.NewDecoder(w.Body).Decode(&run)
		if run.State != RunStateCancelled {
			t.Fatalf("expected state=canceled, got %s", run.State)
		}
	})

	t.Run("rejects invalid transition", func(t *testing.T) {
		// run-1 is now canceled, can't pause it
		body := `{"action":"pause"}`
		req := httptest.NewRequest(http.MethodPost, "/api/runs/run-1/control", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		srv.handleRunControl(w, req)

		if w.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d", w.Code)
		}
	})

	t.Run("retries a canceled run", func(t *testing.T) {
		body := `{"action":"retry","reason":"Fix applied"}`
		req := httptest.NewRequest(http.MethodPost, "/api/runs/run-1/control", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		srv.handleRunControl(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}

		var run operatorRunState
		_ = json.NewDecoder(w.Body).Decode(&run)
		if run.State != RunStateRunning {
			t.Fatalf("expected state=running, got %s", run.State)
		}
	})

	t.Run("returns 404 for unknown run", func(t *testing.T) {
		body := `{"action":"pause"}`
		req := httptest.NewRequest(http.MethodPost, "/api/runs/nonexistent/control", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		srv.handleRunControl(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})
}

func TestOperatorRunReceipts(t *testing.T) {
	srv := newTestOperatorServer()
	srv.intents["int-1"] = &operatorIntent{
		IntentID: "int-1", Type: "deploy", Description: "Deploy v1",
		Status: IntentExecuting, CreatedAt: "2026-02-08T00:00:00Z", UpdatedAt: "2026-02-08T00:00:00Z",
	}
	srv.operatorRuns["run-1"] = &operatorRunState{
		RunID: "run-1", IntentID: "int-1", State: RunStateRunning,
		CreatedAt: "2026-02-08T00:00:00Z", UpdatedAt: "2026-02-08T00:00:00Z",
	}

	t.Run("returns receipts for run", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/runs/run-1/receipts", nil)
		w := httptest.NewRecorder()

		srv.handleRunReceipts(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}

		var resp map[string]interface{}
		_ = json.NewDecoder(w.Body).Decode(&resp)
		if resp["run_id"] != "run-1" {
			t.Fatalf("expected run_id=run-1, got %v", resp["run_id"])
		}
	})
}

func TestOperatorRunReplay(t *testing.T) {
	srv := newTestOperatorServer()
	srv.intents["int-1"] = &operatorIntent{
		IntentID: "int-1", Type: "deploy", Description: "Deploy v1",
		Status: IntentExecuting, CreatedAt: "2026-02-08T00:00:00Z", UpdatedAt: "2026-02-08T00:00:00Z",
	}
	srv.operatorRuns["run-1"] = &operatorRunState{
		RunID: "run-1", IntentID: "int-1", State: RunStateRunning, PlanHash: "sha256:abc",
		CreatedAt: "2026-02-08T00:00:00Z", UpdatedAt: "2026-02-08T00:00:00Z",
	}

	t.Run("returns replay result with diffs", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/runs/run-1/replay", nil)
		w := httptest.NewRecorder()

		srv.handleRunReplay(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}

		var result replayResult
		_ = json.NewDecoder(w.Body).Decode(&result)

		if result.RunID != "run-1" {
			t.Fatalf("expected run_id=run-1, got %s", result.RunID)
		}
		if result.PlannedSteps == 0 {
			t.Fatal("planned_steps should not be 0")
		}
		if result.MatchPercentage != 100.0 {
			t.Fatalf("expected 100%% match, got %.1f%%", result.MatchPercentage)
		}
		if result.ReceiptID == "" {
			t.Fatal("receipt_id should not be empty")
		}
		if len(result.Diffs) == 0 {
			t.Fatal("diffs should not be empty")
		}
	})

	t.Run("returns 404 for unknown run", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/runs/nonexistent/replay", nil)
		w := httptest.NewRecorder()

		srv.handleRunReplay(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})
}

// TestOperatorFullLoop tests the complete operator lifecycle:
// Create → Plan → Submit → Approve → Run → Receipts → Replay
func TestOperatorFullLoop(t *testing.T) {
	srv := newTestOperatorServer()

	// 1. Create intent
	body := `{"type":"deploy","description":"Deploy trading engine v2.1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/intents", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	srv.handleCreateIntent(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", w.Code)
	}

	var intent operatorIntent
	_ = json.NewDecoder(w.Body).Decode(&intent)
	intentID := intent.IntentID

	// 2. Plan preview
	req = httptest.NewRequest(http.MethodPost, "/api/intents/"+intentID+"/plan", nil)
	w = httptest.NewRecorder()
	srv.handlePlanIntent(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("plan: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var plan operatorPlan
	_ = json.NewDecoder(w.Body).Decode(&plan)
	if plan.Hash == "" {
		t.Fatal("plan hash should not be empty")
	}

	// 3. Submit for approval
	req = httptest.NewRequest(http.MethodPost, "/api/intents/"+intentID+"/submit", nil)
	w = httptest.NewRecorder()
	srv.handleSubmitIntent(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("submit: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// 4. Approve
	approveBody := `{"decision":"approve","reason":"Plan reviewed, risk acceptable","decided_by":"release-captain"}`
	req = httptest.NewRequest(http.MethodPost, "/api/approvals/"+intentID+"/approve", bytes.NewBufferString(approveBody))
	w = httptest.NewRecorder()
	srv.handleApproveIntent(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("approve: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// 5. Create run
	runBody := `{"intent_id":"` + intentID + `"}`
	req = httptest.NewRequest(http.MethodPost, "/api/runs", bytes.NewBufferString(runBody))
	w = httptest.NewRecorder()
	srv.handleCreateRun(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("run: expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var run operatorRunState
	_ = json.NewDecoder(w.Body).Decode(&run)
	runID := run.RunID

	// 6. Get receipts
	req = httptest.NewRequest(http.MethodGet, "/api/runs/"+runID+"/receipts", nil)
	w = httptest.NewRecorder()
	srv.handleRunReceipts(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("receipts: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var receiptsResp map[string]interface{}
	_ = json.NewDecoder(w.Body).Decode(&receiptsResp)
	total := receiptsResp["total"].(float64)
	if total == 0 {
		t.Fatal("expected at least 1 receipt for run")
	}

	// 7. Replay
	req = httptest.NewRequest(http.MethodPost, "/api/runs/"+runID+"/replay", nil)
	w = httptest.NewRecorder()
	srv.handleRunReplay(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("replay: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var replay replayResult
	_ = json.NewDecoder(w.Body).Decode(&replay)

	if replay.MatchPercentage != 100.0 {
		t.Fatalf("expected 100%% match, got %.1f%%", replay.MatchPercentage)
	}

	// Verify audit trail: check that receipts were created for each action
	store := srv.receiptStore.(*inMemoryReceiptStore)
	if len(store.receipts) < 5 {
		t.Fatalf("expected at least 5 audit receipts, got %d", len(store.receipts))
	}

	// Verify receipt tools cover the full loop
	tools := make(map[string]bool)
	for _, r := range store.receipts {
		if tool, ok := r.Metadata["tool"]; ok {
			tools[tool.(string)] = true
		}
	}

	expectedTools := []string{
		"operator_create_intent",
		"operator_plan_intent",
		"operator_submit_intent",
		"operator_approve_intent",
		"operator_create_run",
	}
	for _, tool := range expectedTools {
		if !tools[tool] {
			t.Errorf("missing audit receipt for tool %q", tool)
		}
	}
}
