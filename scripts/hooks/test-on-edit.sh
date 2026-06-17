#!/usr/bin/env bash
# PostToolUse hook: runs make test after any .go file is edited.
set -o errexit
set -o nounset
set -o pipefail

FILE=$(jq -r '.tool_input.file_path // empty' 2>/dev/null) || exit 0

[ -z "$FILE" ] && exit 0
echo "$FILE" | grep -qE '\.go$' || exit 0

TOPLEVEL=$(git rev-parse --show-toplevel 2>/dev/null) || exit 0

cd "$TOPLEVEL" || exit 1
make test 2>&1 | tail -20
