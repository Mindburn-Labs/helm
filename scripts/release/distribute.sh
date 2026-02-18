#!/bin/bash
set -e

# HELM Distribution Script
# Usage: ./scripts/release/distribute.sh [version]
# Example: ./scripts/release/distribute.sh v0.1.0-sota

VERSION=$1
if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

echo "🚀 Distributing HELM $VERSION..."

# 1. Verification
echo "🔍 Verifying artifacts..."
if [ ! -f "bin/SHA256SUMS.txt" ]; then
    echo "❌ bin/SHA256SUMS.txt not found. Run 'make release-binaries' first."
    exit 1
fi
shasum -c bin/SHA256SUMS.txt --ignore-missing
echo "✅ Artifacts verified."

# 2. Docker Publish
echo "🐳 Publishing Docker image..."
if [ -z "$DOCKER_REPO" ]; then
    echo "⚠️  DOCKER_REPO not set. Skipping Docker publish."
else
    docker tag helm:latest "$DOCKER_REPO/helm:$VERSION"
    docker tag helm:latest "$DOCKER_REPO/helm:latest"
    docker push "$DOCKER_REPO/helm:$VERSION"
    docker push "$DOCKER_REPO/helm:latest"
    echo "✅ Docker image published."
fi

# 3. NPM Publish
echo "📦 Publishing NPM package..."
if [ -z "$NPM_TOKEN" ]; then
    echo "⚠️  NPM_TOKEN not set. Skipping NPM publish."
else
    cd sdk/ts
    # Ensure version matches
    npm version "$VERSION" --no-git-tag-version --allow-same-version
    echo "//registry.npmjs.org/:_authToken=$NPM_TOKEN" > .npmrc
    npm publish --access public
    rm .npmrc
    cd ../..
    echo "✅ NPM package published."
fi

# 4. PyPI Publish
echo "🐍 Publishing PyPI package..."
if [ -z "$PYPI_TOKEN" ]; then
    echo "⚠️  PYPI_TOKEN not set. Skipping PyPI publish."
else
    cd sdk/python
    # Ensure build tools
    pip install -q build twine
    python3 -m build
    twine upload dist/* -u __token__ -p "$PYPI_TOKEN"
    cd ../..
    echo "✅ PyPI package published."
fi

# 5. GitHub Release
echo "🐙 Creating GitHub Release..."
if ! command -v gh &> /dev/null; then
    echo "⚠️  'gh' CLI not found. Skipping GitHub Release."
else
    # Check if tag exists locally, if not create it
    if ! git rev-parse "$VERSION" >/dev/null 2>&1; then
        echo "Creating git tag $VERSION..."
        git tag "$VERSION"
        git push origin "$VERSION"
    fi

    gh release create "$VERSION" 
        bin/helm-darwin-amd64 
        bin/helm-darwin-arm64 
        bin/helm-linux-amd64 
        bin/helm-linux-arm64 
        bin/SHA256SUMS.txt 
        artifacts/99_final/sbom.json 
        --title "HELM $VERSION (SOTA)" 
        --notes-file RELEASE_NOTES.md 
        --draft

    echo "✅ GitHub Release created (Draft)."
fi

echo "🎉 Distribution complete!"
