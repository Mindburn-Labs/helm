# Use Cases Directory Consolidation

The use case documentation has been consolidated:

- **`docs/use-cases/`** — Canonical directory containing executable shell scripts (UC-001 through UC-012)
- **`docs/use-cases/`** — Markdown descriptions for each use case

Both directories follow the same UC-NNN numbering. The markdown directory provides human-readable descriptions while the scripts directory provides executable demonstrations.

## Mapping

| ID | Script | Description |
|----|--------|-------------|
| UC-001 | `use_cases/UC-001_pep_allow_safe.sh` | PEP allows safe tool execution |
| UC-002 | `use_cases/UC-002_schema_mismatch.sh` | Schema mismatch rejection |
| UC-003 | `use_cases/UC-003_approval_ceremony.sh` | HITL approval ceremony |
| UC-004 | `use_cases/UC-004_wasi_transform.sh` | WASI sandbox transform |
| UC-005 | `use_cases/UC-005_wasi_gas_exhaustion.sh` | Gas exhaustion trap |
| UC-006 | `use_cases/UC-006_idempotency.sh` | Idempotency enforcement |
| UC-007 | `use_cases/UC-007_proofgraph_export.sh` | ProofGraph export |
| UC-008 | `use_cases/UC-008_replay_verify.sh` | Replay verification |
| UC-009 | `use_cases/UC-009_connector_drift.sh` | Connector drift detection |
| UC-010 | `use_cases/UC-010_trust_key_rotation.sh` | Trust key rotation |
| UC-011 | `use_cases/UC-011_island_mode.sh` | Island mode operation |
| UC-012 | `use_cases/UC-012_openai_proxy.sh` | OpenAI proxy governance |

> [!NOTE]
> Both directories are retained for backward compatibility. The `use_cases/` directory
> with executable scripts is authoritative. The `use-cases/` markdown files serve as
> documentation companions.
