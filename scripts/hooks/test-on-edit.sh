#!/usr/bin/env bash
# PostToolUse hook: runs make test after any .go file is edited.
set -o pipefail

FILE=$(jq -r '.tool_input.file_path // empty')

[ -z "$FILE" ] && exit 0
echo "$FILE" | grep -qE '\.go$' || exit 0

TOPLEVEL=$(git rev-parse --show-toplevel 2>/dev/null)
[ -z "$TOPLEVEL" ] && exit 0

cd "$TOPLEVEL" || exit
make test 2>&1 | tail -20
