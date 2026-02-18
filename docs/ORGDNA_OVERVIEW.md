# OrgDNA: Organizational Policy Genome

## Overview

OrgDNA is HELM's mechanism for encoding an organization's policy genome — the complete set of rules, constraints, approval chains, and risk thresholds that govern AI agent execution.

Unlike static config files, OrgDNA packs are:
- **Content-addressed**: every version is identifiable by SHA-256 hash
- **Composable**: packs can extend and override each other
- **Auditable**: changes are tracked with provenance and timestamps
- **Machine-readable**: JSON Schema validated, parseable by policy engines

## Core Concepts

### OrgDNA Pack

A pack is a JSON document conforming to `schemas/orgdna.schema.json`:

```json
{
  "schema_version": "1.0.0",
  "org_id": "acme-corp",
  "pack_id": "customer-support-agents",
  "policies": [
    {
      "id": "p001",
      "name": "Refund Approval",
      "action": "APPROVE_REFUND",
      "constraints": {
        "max_amount_usd": 500,
        "requires_approval": ["manager"],
        "cool_down_minutes": 15
      }
    }
  ],
  "risk_thresholds": {
    "max_tokens_per_request": 4096,
    "max_requests_per_minute": 60,
    "budget_usd_daily": 100
  },
  "approval_chains": [
    {
      "action_pattern": "APPROVE_REFUND",
      "chain": ["agent", "manager"],
      "escalation_timeout_minutes": 30
    }
  ]
}
```

### Content Addressing

Every OrgDNA pack is hashed deterministically:

```bash
helm orgdna hash --pack orgdna.json
# sha256:a1b2c3d4...
```

This hash appears in receipts, allowing verifiers to confirm exactly which organizational policies were in effect during any AI decision.

### Validation

```bash
helm orgdna validate --pack orgdna.json
# ✅ Valid OrgDNA pack: acme-corp/customer-support-agents
# Schema version: 1.0.0
# Policies: 3
# Risk thresholds: 4
# Approval chains: 2
```

## Integration with HELM

```
OrgDNA Pack → PolicyDecisionPoint → Guardian → Receipt → EvidencePack
     │                                            │
     hash ──────────────────────────────────────→ policy_content_hash
```

The `policy_content_hash` field in each `DecisionRecord` references the OrgDNA pack hash, creating a verifiable link between organizational policy and individual AI decisions.

## Example Packs

| Pack | Use Case |
|------|----------|
| `examples/orgdna/saas_support.json` | Customer support agents with refund limits |
| `examples/orgdna/finance_approval.json` | Financial approval chains with budget gates |

## Security Model

- OrgDNA packs MUST be signed by authorized personnel
- Changes to packs MUST produce a new hash
- The hash appears in every receipt created under that policy
- Auditors can verify: *"this decision was made under org policy X at version Y"*
