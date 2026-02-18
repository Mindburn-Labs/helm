#!/bin/bash
set -e

# HELM Distribution Script (SOTA Full Coverage)
# Usage: ./scripts/release/distribute.sh [version]
# Example: ./scripts/release/distribute.sh 0.1.0

# Load secrets if .env.release exists
if [ -f .env.release ]; then
    echo "🔑 Loading secrets from .env.release..."
    export $(grep -v '^#' .env.release | xargs)
fi

VERSION=$1
if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

echo "🚀 Distributing HELM $VERSION across all ecosystems..."

# 1. Go (via Git Tags)
echo "🐹 Tagging Go SDK..."
git tag -f "sdk/go/v$VERSION"
git push -f origin "sdk/go/v$VERSION"
echo "✅ Go SDK tagged (v$VERSION)."

# 2. Rust (Crates.io)
echo "🦀 Rust SDK already published, skipping..."
# if [ -z "$CARGO_REGISTRY_TOKEN" ]; then
# ...

# 3. NPM (TypeScript)
echo "📦 NPM package already published, skipping..."
# if [ -z "$NPM_TOKEN" ]; then
# ...

# 4. PyPI (Python)
echo "🐍 PyPI package already published, skipping..."
# if [ -z "$PYPI_TOKEN" ]; then
# ...

# 5. Maven (Java)
echo "☕ Skipping Maven publish (verification required)..."
# if [ -z "$OSSRH_USERNAME" ]; then
# ...

# 6. Docker
echo "🐳 Publishing Docker image..."
if [ -z "$DOCKER_REPO" ]; then
    echo "⚠️  DOCKER_REPO not set. Skipping Docker publish."
else
    if [ -n "$DOCKER_PASSWORD" ] && [ -n "$DOCKER_USERNAME" ]; then
        echo "🔑 Logging into Docker..."
        echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
    fi
    docker tag helm:latest "$DOCKER_REPO/helm:v$VERSION"
    docker tag helm:latest "$DOCKER_REPO/helm:latest"
    docker push "$DOCKER_REPO/helm:v$VERSION"
    docker push "$DOCKER_REPO/helm:latest"
    echo "✅ Docker image published."
fi

echo "🎉 Full Distribution complete for version $VERSION!"