# Cedar PDP Sidecar

Minimal Cedar policy evaluation sidecar for HELM.

## Quick Start

```bash
docker compose up -d
# POST http://localhost:8182/decide
```

## API

### POST /decide

```json
{
  "principal": "User::\"alice\"",
  "action": "Action::\"read\"",
  "resource": "Resource::\"document\"",
  "context": {}
}
```

Response:
```json
{
  "decision": "Allow",
  "diagnostics": {
    "reason": [],
    "errors": []
  }
}
```

## Policy Files

Place Cedar policies in `policies/` and entities in `entities/`:
- `policies/*.cedar` — Cedar policy files
- `entities/entities.json` — Entity store

## Building

See the Dockerfile or run locally with a Cedar CLI:
```bash
cedar authorize --policies policies/ --entities entities/entities.json
```
