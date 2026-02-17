#!/usr/bin/env bash
# HELM Examples Smoke Test
# Validates all example directories have correct structure.
# Optionally builds compilable examples.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
EXAMPLES_DIR="$PROJECT_ROOT/examples"

PASS=0
FAIL=0

check_example() {
    local dir="$1"
    local name
    name="$(basename "$dir")"
    echo -n "  $name ... "

    # Every example must have a README.md
    if [ ! -f "$dir/README.md" ]; then
        echo "❌ FAIL (missing README.md)"
        FAIL=$((FAIL + 1))
        return
    fi

    # Check for main entrypoint
    local has_entrypoint=false
    for f in main.py main.ts main.js main.go main.sh main.rs Main.java; do
        if [ -f "$dir/$f" ]; then
            has_entrypoint=true
            break
        fi
    done

    # Check for pom.xml (Java) or Cargo.toml (Rust) as entrypoint indicator
    if [ -f "$dir/pom.xml" ] || [ -f "$dir/Cargo.toml" ]; then
        has_entrypoint=true
    fi

    if [ "$has_entrypoint" = false ]; then
        echo "❌ FAIL (no main entrypoint found)"
        FAIL=$((FAIL + 1))
        return
    fi

    echo "✅ PASS"
    PASS=$((PASS + 1))
}

echo "HELM Examples Smoke Test"
echo "════════════════════════"
echo ""

for dir in "$EXAMPLES_DIR"/*/; do
    [ -d "$dir" ] && check_example "$dir"
done

# Build-check Go example if go is available
if command -v go &>/dev/null && [ -d "$EXAMPLES_DIR/go_client" ]; then
    echo ""
    echo "  go_client build check ... "
    if (cd "$EXAMPLES_DIR/go_client" && go build -o /dev/null . 2>&1); then
        echo "  go_client build ✅"
    else
        echo "  go_client build ⚠️ (non-fatal)"
    fi
fi

echo ""
echo "════════════════════════"
echo "Results: $PASS passed, $FAIL failed"

if [ "$FAIL" -gt 0 ]; then
    exit 1
fi
