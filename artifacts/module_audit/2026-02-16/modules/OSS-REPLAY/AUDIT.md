# OSS-REPLAY Audit

## Scope
Offline replay of VCR tapes. Tape integrity verification. Effect re-execution with deterministic driver. Disconnected machine verification.

## Reality
- **replay_cmd.go** (145 lines): `helm replay --evidence --verify`. Loads tapes, checks sequence/data_class/manifest. [CODE]
- **replay.go**: VCR tape format, TapeEntry struct. [CODE]
- **Quality: Partial** — Tape verification works. No effect re-execution.

## Gaps
| # | Gap | Severity |
|---|-----|----------|
| 1 | No effect re-execution | P1 |
| 2 | No deterministic replay driver | P1 |

## Recommendations
1. Add replay driver that re-executes effects against tape expectations.
2. Hash comparison: replayed output vs recorded output.
