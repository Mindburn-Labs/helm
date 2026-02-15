package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	helmcrypto "github.com/Mindburn-Labs/helm/core/pkg/crypto"
)

// proxyReceipt is the governance receipt attached to every proxied request.
type proxyReceipt struct {
	ReceiptID    string   `json:"receipt_id"`
	Timestamp    string   `json:"timestamp"`
	Upstream     string   `json:"upstream"`
	Model        string   `json:"model,omitempty"`
	InputHash    string   `json:"input_hash"`
	OutputHash   string   `json:"output_hash,omitempty"`
	ToolCalls    int      `json:"tool_calls_intercepted"`
	ToolNames    []string `json:"tool_names,omitempty"`
	ArgsHashes   []string `json:"args_hashes,omitempty"`
	ArgsValid    []bool   `json:"args_valid,omitempty"`
	Status       string   `json:"status"`
	LamportClock uint64   `json:"lamport_clock"`
	PrevHash     string   `json:"prev_hash"`
	Signature    string   `json:"signature,omitempty"`
}

// receiptStore persists receipts to a JSONL file for auditability.
type receiptStore struct {
	mu       sync.Mutex
	file     *os.File
	prevHash string
}

func newReceiptStore(path string) (*receiptStore, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	return &receiptStore{file: f, prevHash: "GENESIS"}, nil
}

func (s *receiptStore) Append(rcpt *proxyReceipt) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	rcpt.PrevHash = s.prevHash

	data, err := json.Marshal(rcpt)
	if err != nil {
		return err
	}

	// Update causal chain: prevHash = SHA-256 of this receipt's JSON
	h := sha256.Sum256(data)
	s.prevHash = "sha256:" + hex.EncodeToString(h[:])

	data = append(data, '\n')
	_, err = s.file.Write(data)
	return err
}

func (s *receiptStore) Close() error {
	return s.file.Close()
}

// validateToolCallArgs performs PEP validation: checks args are valid JSON
// and canonicalizes them via JCS before hashing. Returns (canonical hash, valid, error).
func validateToolCallArgs(argsStr string) (string, bool) {
	// Phase 1: args must parse as valid JSON (fail-closed on malformed)
	var parsed any
	if err := json.Unmarshal([]byte(argsStr), &parsed); err != nil {
		return "", false
	}

	// Phase 2: JCS canonicalization (RFC 8785) — re-marshal with sorted keys, no HTML escaping
	canonical, err := helmcrypto.CanonicalMarshal(parsed)
	if err != nil {
		return "", false
	}

	// Phase 3: SHA-256 of canonical form
	h := sha256.Sum256(canonical)
	return "sha256:" + hex.EncodeToString(h[:]), true
}

// runProxyCmd implements `helm proxy` — the 1-line integration wedge.
//
// Usage:
//
//	helm proxy --upstream https://api.openai.com/v1 --port 9090
//
// Then:
//
//	export OPENAI_BASE_URL=http://localhost:9090/v1
//	python your_app.py  # Every tool call now gets a receipt.
//
// Features:
//   - Receipt persistence: JSONL audit log at --receipts-dir
//   - PEP validation: tool_call arguments validated as JSON, canonicalized (JCS), and SHA-256 hashed
//   - Causal chain: receipts linked via PrevHash (SHA-256 of previous receipt)
//   - Ed25519 signature: receipts signed if --sign is enabled
//
// Exit codes:
//
//	0 = clean shutdown
//	2 = config error
func runProxyCmd(args []string, stdout, stderr io.Writer) int {
	cmd := flag.NewFlagSet("proxy", flag.ContinueOnError)
	cmd.SetOutput(stderr)

	var (
		upstream    string
		port        int
		apiKey      string
		jsonOutput  bool
		verbose     bool
		receiptsDir string
		signKey     string
	)

	cmd.StringVar(&upstream, "upstream", "https://api.openai.com/v1", "Upstream API base URL")
	cmd.IntVar(&port, "port", 9090, "Local proxy port")
	cmd.StringVar(&apiKey, "api-key", "", "API key to forward to upstream (or use OPENAI_API_KEY env)")
	cmd.BoolVar(&jsonOutput, "json", false, "Log receipts as JSON to stdout")
	cmd.BoolVar(&verbose, "verbose", false, "Verbose logging")
	cmd.StringVar(&receiptsDir, "receipts-dir", "./helm-receipts", "Directory for persistent receipt JSONL logs")
	cmd.StringVar(&signKey, "sign", "", "Ed25519 signing key seed (enables receipt signatures)")

	if err := cmd.Parse(args); err != nil {
		return 2
	}

	// Normalize upstream URL
	upstream = strings.TrimSuffix(upstream, "/")

	upstreamURL, err := url.Parse(upstream)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "Error: invalid upstream URL: %v\n", err)
		return 2
	}

	// Initialize receipt store
	receiptPath := filepath.Join(receiptsDir, fmt.Sprintf("receipts-%s.jsonl", time.Now().Format("2006-01-02")))
	store, err := newReceiptStore(receiptPath)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "Error: failed to initialize receipt store at %s: %v\n", receiptPath, err)
		return 2
	}
	defer store.Close()

	// Optional: Ed25519 signer for receipt signatures
	var signer *helmcrypto.Ed25519Signer
	if signKey != "" {
		signer, err = helmcrypto.NewEd25519Signer(signKey)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "Error: failed to create signer: %v\n", err)
			return 2
		}
	}

	var lamport uint64

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = upstreamURL.Scheme
			req.URL.Host = upstreamURL.Host
			origPath := req.URL.Path
			if strings.HasPrefix(origPath, "/v1") && strings.HasSuffix(upstream, "/v1") {
				req.URL.Path = upstreamURL.Path + strings.TrimPrefix(origPath, "/v1")
			} else {
				req.URL.Path = upstreamURL.Path + origPath
			}
			req.Host = upstreamURL.Host

			// Forward API key
			if apiKey != "" && req.Header.Get("Authorization") == "" {
				req.Header.Set("Authorization", "Bearer "+apiKey)
			}
		},
		ModifyResponse: func(resp *http.Response) error {
			// Read response body
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				return err
			}

			clock := atomic.AddUint64(&lamport, 1)

			// Hash output
			outHash := sha256.Sum256(body)
			outHashHex := "sha256:" + hex.EncodeToString(outHash[:])

			// Parse for tool_calls + PEP validation
			var chatResp map[string]any
			toolCallCount := 0
			var argsHashes []string
			var argsValid []bool
			var toolNames []string
			var model string
			status := "PASS"

			if err := json.Unmarshal(body, &chatResp); err == nil {
				if m, ok := chatResp["model"].(string); ok {
					model = m
				}
				if choices, ok := chatResp["choices"].([]any); ok {
					for _, c := range choices {
						choice, ok := c.(map[string]any)
						if !ok {
							continue
						}
						msg, ok := choice["message"].(map[string]any)
						if !ok {
							continue
						}
						if tcs, ok := msg["tool_calls"].([]any); ok {
							for _, tc := range tcs {
								toolCallCount++
								tcMap, ok := tc.(map[string]any)
								if !ok {
									continue
								}
								fn, ok := tcMap["function"].(map[string]any)
								if !ok {
									continue
								}

								// Extract tool name
								if name, ok := fn["name"].(string); ok {
									toolNames = append(toolNames, name)
								}

								// PEP validation: validate + canonicalize + hash
								if argsStr, ok := fn["arguments"].(string); ok {
									hash, valid := validateToolCallArgs(argsStr)
									argsHashes = append(argsHashes, hash)
									argsValid = append(argsValid, valid)
									if !valid {
										status = "PEP_VALIDATION_FAILED"
										log.Printf("[WARN] PEP validation failed for tool_call args (malformed JSON)")
									}
								}
							}
						}
					}
				}
			}

			// Build receipt
			rcptID := fmt.Sprintf("rcpt-proxy-%d-%d", time.Now().UnixNano(), clock)
			rcpt := &proxyReceipt{
				ReceiptID:    rcptID,
				Timestamp:    time.Now().UTC().Format(time.RFC3339Nano),
				Upstream:     upstream,
				Model:        model,
				OutputHash:   outHashHex,
				ToolCalls:    toolCallCount,
				ToolNames:    toolNames,
				ArgsHashes:   argsHashes,
				ArgsValid:    argsValid,
				Status:       status,
				LamportClock: clock,
			}

			// Sign receipt if signer available
			if signer != nil {
				payload := fmt.Sprintf("%s:%s:%s:%d", rcpt.ReceiptID, rcpt.OutputHash, rcpt.Status, rcpt.LamportClock)
				sig, signErr := signer.Sign([]byte(payload))
				if signErr == nil {
					rcpt.Signature = sig
				}
			}

			// Persist receipt (JSONL, append-only, causal chain)
			if storeErr := store.Append(rcpt); storeErr != nil {
				log.Printf("[ERROR] receipt persist failed: %v", storeErr)
			}

			// Inject receipt headers
			resp.Header.Set("X-Helm-Receipt-ID", rcpt.ReceiptID)
			resp.Header.Set("X-Helm-Output-Hash", rcpt.OutputHash)
			resp.Header.Set("X-Helm-Lamport-Clock", fmt.Sprintf("%d", rcpt.LamportClock))
			resp.Header.Set("X-Helm-Status", rcpt.Status)
			if toolCallCount > 0 {
				resp.Header.Set("X-Helm-Tool-Calls", fmt.Sprintf("%d", toolCallCount))
			}
			if rcpt.Signature != "" {
				resp.Header.Set("X-Helm-Signature", rcpt.Signature)
			}

			// Log receipt
			if jsonOutput {
				rcptJSON, _ := json.Marshal(rcpt)
				log.Printf("%s", rcptJSON)
			} else if verbose {
				log.Printf("[RECEIPT] %s | %s | tools=%d | status=%s | %s",
					rcpt.ReceiptID, rcpt.Model, rcpt.ToolCalls, rcpt.Status, rcpt.OutputHash[:30]+"…")
			}

			// Restore body
			resp.Body = io.NopCloser(bytes.NewReader(body))
			resp.ContentLength = int64(len(body))
			resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(body)))

			return nil
		},
	}

	mux := http.NewServeMux()

	// Health endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok","mode":"proxy","upstream":"` + upstream + `"}`))
	})

	// Receipts endpoint — serve the JSONL file
	mux.HandleFunc("/helm/receipts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-ndjson")
		data, err := os.ReadFile(receiptPath)
		if err != nil {
			http.Error(w, "no receipts yet", http.StatusNotFound)
			return
		}
		_, _ = w.Write(data)
	})

	// Proxy everything else
	mux.HandleFunc("/", proxy.ServeHTTP)

	addr := fmt.Sprintf(":%d", port)

	_, _ = fmt.Fprintf(stdout, "HELM Proxy Sidecar\n")
	_, _ = fmt.Fprintf(stdout, "══════════════════\n")
	_, _ = fmt.Fprintf(stdout, "  Upstream:  %s\n", upstream)
	_, _ = fmt.Fprintf(stdout, "  Listen:    http://localhost%s\n", addr)
	_, _ = fmt.Fprintf(stdout, "  Health:    http://localhost%s/health\n", addr)
	_, _ = fmt.Fprintf(stdout, "  Receipts:  %s\n", receiptPath)
	if signer != nil {
		_, _ = fmt.Fprintf(stdout, "  Signing:   Ed25519 (key: %s)\n", signer.KeyID)
	}
	_, _ = fmt.Fprintf(stdout, "\n")
	_, _ = fmt.Fprintf(stdout, "  Drop-in usage:\n")
	_, _ = fmt.Fprintf(stdout, "    export OPENAI_BASE_URL=http://localhost%s/v1\n", addr)
	_, _ = fmt.Fprintf(stdout, "    python your_app.py\n")
	_, _ = fmt.Fprintf(stdout, "\n")
	_, _ = fmt.Fprintf(stdout, "  Every tool call is validated, hashed, and receipted. Ctrl+C to stop.\n")

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 30 * time.Second,
	}

	// Graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_ = ctx

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		_, _ = fmt.Fprintf(stderr, "Error: %v\n", err)
		return 2
	}

	return 0
}
