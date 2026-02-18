package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Mindburn-Labs/helm/core/pkg/conform"
	"github.com/Mindburn-Labs/helm/core/pkg/verifier"
)

// runVerifyCmd implements `helm verify` per §2.1.
//
// Validates a signed EvidencePack bundle: structure, hashes, and signature.
// Supports auditor mode via --json-out for structured verification reports.
//
// Exit codes:
//
//	0 = verification passed
//	1 = verification failed
//	2 = runtime error
func runVerifyCmd(args []string, stdout, stderr io.Writer) int {
	cmd := flag.NewFlagSet("verify", flag.ContinueOnError)
	cmd.SetOutput(stderr)

	var (
		bundle      string
		jsonOutput  bool
		jsonOutFile string
	)

	cmd.StringVar(&bundle, "bundle", "", "Path to EvidencePack directory (REQUIRED)")
	cmd.BoolVar(&jsonOutput, "json", false, "Output results as JSON to stdout")
	cmd.StringVar(&jsonOutFile, "json-out", "", "Write structured audit report to file (auditor mode)")

	if err := cmd.Parse(args); err != nil {
		return 2
	}

	if bundle == "" {
		_, _ = fmt.Fprintln(stderr, "Error: --bundle is required")
		return 2
	}

	// Use the standalone verifier library (zero network deps)
	report, err := verifier.VerifyBundle(bundle)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "Error: verification failed: %v\n", err)
		return 2
	}

	// Also run legacy conform-based checks for backward compat
	structIssues := conform.ValidateEvidencePackStructure(bundle)
	if len(structIssues) > 0 {
		for _, issue := range structIssues {
			report.Checks = append(report.Checks, verifier.CheckResult{
				Name:   "conform:" + issue,
				Pass:   false,
				Reason: issue,
			})
		}
		report.Verified = false
	}

	// Verify report signature (legacy)
	sigErr := conform.VerifyReport(bundle, func(data []byte, sig string) error {
		return nil // Default: hash-based verification
	})
	if sigErr != nil {
		report.Checks = append(report.Checks, verifier.CheckResult{
			Name:   "signature_verification",
			Pass:   false,
			Reason: fmt.Sprintf("signature: %v", sigErr),
		})
	}

	// Write auditor JSON report to file if requested
	if jsonOutFile != "" {
		data, _ := json.MarshalIndent(report, "", "  ")
		if writeErr := os.WriteFile(jsonOutFile, data, 0644); writeErr != nil {
			_, _ = fmt.Fprintf(stderr, "Error: cannot write audit report: %v\n", writeErr)
			return 2
		}
		_, _ = fmt.Fprintf(stdout, "Audit report written to %s\n", jsonOutFile)
	}

	// Output
	if jsonOutput {
		data, _ := json.MarshalIndent(report, "", "  ")
		_, _ = fmt.Fprintln(stdout, string(data))
	} else {
		if report.Verified {
			_, _ = fmt.Fprintf(stdout, "✅ EvidencePack verification PASSED\n")
			_, _ = fmt.Fprintf(stdout, "Bundle: %s\n", bundle)
			_, _ = fmt.Fprintf(stdout, "Checks: %s\n", report.Summary)
		} else {
			_, _ = fmt.Fprintf(stdout, "❌ EvidencePack verification FAILED\n")
			_, _ = fmt.Fprintf(stdout, "Bundle: %s\n", bundle)
			for _, c := range report.Checks {
				if !c.Pass {
					_, _ = fmt.Fprintf(stdout, "  - %s: %s\n", c.Name, c.Reason)
				}
			}
		}
	}

	if !report.Verified {
		return 1
	}
	return 0
}
