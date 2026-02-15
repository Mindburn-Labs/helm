# ── Stage 1: Build ─────────────────────────────────────
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /src
COPY core/ ./core/

WORKDIR /src/core
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /helm ./cmd/helm/
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /helm-node ./cmd/helm-node/

# ── Stage 2: Runtime ───────────────────────────────────
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /helm /usr/local/bin/helm
COPY --from=builder /helm-node /usr/local/bin/helm-node

EXPOSE 8080 9090

ENTRYPOINT ["helm-node"]
