#!/usr/bin/env bash
# PostToolUse hook: runs make test after any .go file is edited.

FILE=$(jq -r '.tool_input.file_path // empty')

[ -z "$FILE" ] && exit 0
echo "$FILE" | grep -qE '\.go$' || exit 0

cd "$(git rev-parse --show-toplevel)"
make test 2>&1 | tail -20
exit 0
