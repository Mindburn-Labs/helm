#!/bin/bash
set -e

# HELM Distribution Script (SOTA Full Coverage)
# Usage: ./scripts/release/distribute.sh [version]
# Example: ./scripts/release/distribute.sh 0.1.0

VERSION=$1
if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

echo "🚀 Distributing HELM $VERSION across all ecosystems..."

# 1. Go (via Git Tags)
echo "🐹 Tagging Go SDK..."
git tag "sdk/go/v$VERSION"
git push origin "sdk/go/v$VERSION"
echo "✅ Go SDK tagged (v$VERSION)."

# 2. Rust (Crates.io)
echo "🦀 Publishing Rust SDK..."
if [ -z "$CARGO_REGISTRY_TOKEN" ]; then
    echo "⚠️  CARGO_REGISTRY_TOKEN not set. Skipping Rust publish."
else
    cd sdk/rust
    cargo publish --token "$CARGO_REGISTRY_TOKEN"
    cd ../..
    echo "✅ Rust SDK published."
fi

# 3. NPM (TypeScript)
echo "📦 Publishing NPM package..."
if [ -z "$NPM_TOKEN" ]; then
    echo "⚠️  NPM_TOKEN not set. Skipping NPM publish."
else
    cd sdk/ts
    npm version "$VERSION" --no-git-tag-version --allow-same-version
    echo "//registry.npmjs.org/:_authToken=$NPM_TOKEN" > .npmrc
    npm publish --access public
    rm .npmrc
    cd ../..
    echo "✅ NPM package published."
fi

# 4. PyPI (Python)
echo "🐍 Publishing PyPI package..."
if [ -z "$PYPI_TOKEN" ]; then
    echo "⚠️  PYPI_TOKEN not set. Skipping PyPI publish."
else
    cd sdk/python
    pip install -q build twine
    python3 -m build
    twine upload dist/* -u __token__ -p "$PYPI_TOKEN"
    cd ../..
    echo "✅ PyPI package published."
fi

# 5. Docker
echo "🐳 Publishing Docker image..."
if [ -z "$DOCKER_REPO" ]; then
    echo "⚠️  DOCKER_REPO not set. Skipping Docker publish."
else
    docker tag helm:latest "$DOCKER_REPO/helm:v$VERSION"
    docker tag helm:latest "$DOCKER_REPO/helm:latest"
    docker push "$DOCKER_REPO/helm:v$VERSION"
    docker push "$DOCKER_REPO/helm:latest"
    echo "✅ Docker image published."
fi

echo "🎉 Full Distribution complete for version $VERSION!"