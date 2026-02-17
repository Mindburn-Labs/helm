# ── Stage 1: Build ─────────────────────────────────────
# SC-004: Base images pinned by digest for supply chain integrity
FROM golang:1.24-alpine@sha256:7772cb5322baa0cee6b21d2b6a97c2a4c2bcf3a3e78e63a0a25d9b5a6a8e4d2f AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /src
COPY core/ ./core/

WORKDIR /src/core
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /helm ./cmd/helm/
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /helm-node ./cmd/helm-node/

# ── Stage 2: Runtime ───────────────────────────────────
FROM gcr.io/distroless/static-debian12:nonroot@sha256:e9ac71e2b8e279a8372741b7a0293afda17650d926900233ec3a7b2b7c22a246

COPY --from=builder /helm /usr/local/bin/helm
COPY --from=builder /helm-node /usr/local/bin/helm-node

EXPOSE 8080 9090

ENTRYPOINT ["helm-node"]
