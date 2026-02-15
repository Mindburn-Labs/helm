package api

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Mindburn-Labs/helm/core/pkg/contracts"
)

// ApproveHandler handles POST /api/v1/kernel/approve
// This is the backend half of the HITL bridge. The frontend uses WebCrypto API
// to sign the intent hash, and this handler verifies the signature cryptographically.
type ApproveHandler struct {
	// pendingApprovals maps intent_hash → ApprovalRequest
	pendingApprovals map[string]*contracts.ApprovalRequest
}

// NewApproveHandler creates a new approval handler.
func NewApproveHandler() *ApproveHandler {
	return &ApproveHandler{
		pendingApprovals: make(map[string]*contracts.ApprovalRequest),
	}
}

// RegisterPendingApproval adds an intent to the pending approval queue.
func (h *ApproveHandler) RegisterPendingApproval(req *contracts.ApprovalRequest) {
	h.pendingApprovals[req.IntentHash] = req
}

// HandleApprove processes a cryptographic approval from the operator UI.
//
// Flow:
//  1. Parse ApprovalReceipt from request body
//  2. Verify the intent exists in pending queue
//  3. Decode the approver's Ed25519 public key
//  4. Verify the signature over IntentHash
//  5. Mark the approval as APPROVED with the receipt
//  6. Return 200 with the signed receipt
func (h *ApproveHandler) HandleApprove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var receipt contracts.ApprovalReceipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, fmt.Sprintf("invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if receipt.IntentHash == "" || receipt.PublicKey == "" || receipt.Signature == "" {
		http.Error(w, "missing required fields: intent_hash, public_key, signature", http.StatusBadRequest)
		return
	}

	// Check if the intent exists in pending queue
	pending, ok := h.pendingApprovals[receipt.IntentHash]
	if !ok {
		http.Error(w, "intent not found or already processed", http.StatusNotFound)
		return
	}

	if pending.Status != contracts.ApprovalPending {
		http.Error(w, fmt.Sprintf("intent already %s", pending.Status), http.StatusConflict)
		return
	}

	// Check expiry
	if time.Now().After(pending.ExpiresAt) {
		pending.Status = contracts.ApprovalExpired
		http.Error(w, "approval request has expired", http.StatusGone)
		return
	}

	// Decode public key
	pubKeyBytes, err := hex.DecodeString(receipt.PublicKey)
	if err != nil || len(pubKeyBytes) != ed25519.PublicKeySize {
		http.Error(w, "invalid public key", http.StatusBadRequest)
		return
	}

	// Decode signature
	sigBytes, err := hex.DecodeString(receipt.Signature)
	if err != nil {
		http.Error(w, "invalid signature encoding", http.StatusBadRequest)
		return
	}

	// Verify Ed25519 signature over the intent hash
	pubKey := ed25519.PublicKey(pubKeyBytes)
	if !ed25519.Verify(pubKey, []byte(receipt.IntentHash), sigBytes) {
		http.Error(w, "signature verification failed — approval rejected", http.StatusForbidden)
		return
	}

	// Signature valid — approve the intent
	receipt.Timestamp = time.Now()
	pending.Status = contracts.ApprovalApproved
	pending.Receipt = &receipt

	// Respond with the signed approval
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "APPROVED",
		"intent_hash": receipt.IntentHash,
		"approver_id": receipt.ApproverID,
		"timestamp":   receipt.Timestamp,
	})
}

// GetPendingApprovals returns all pending approval requests.
func (h *ApproveHandler) GetPendingApprovals() []*contracts.ApprovalRequest {
	var pending []*contracts.ApprovalRequest
	for _, req := range h.pendingApprovals {
		if req.Status == contracts.ApprovalPending {
			pending = append(pending, req)
		}
	}
	return pending
}
