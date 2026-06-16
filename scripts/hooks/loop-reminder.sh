#!/usr/bin/env bash
# Stop hook: injects the next loop step into Claude's context for the next turn.

git rev-parse --git-dir >/dev/null 2>&1 || exit 0

BRANCH=$(git branch --show-current 2>/dev/null)
DIRTY=$(git status --porcelain 2>/dev/null | grep -v '^??' | wc -l | tr -d ' ')
AHEAD=$(git log @{u}.. --oneline 2>/dev/null | wc -l | tr -d ' ')

MSG=""

if [ "$BRANCH" = "main" ] || [ "$BRANCH" = "master" ]; then
  MSG="On main branch — create a feature branch before editing files."
elif [ "$DIRTY" -gt 0 ] && [ "$AHEAD" -eq 0 ]; then
  MSG="Uncommitted changes on '$BRANCH'. Next: /superpowers:verification-before-completion → /commit-commands:commit-push-pr"
elif [ "$AHEAD" -gt 0 ]; then
  MSG="Unpushed commits on '$BRANCH'. Next: /commit-commands:commit-push-pr → /code-review"
fi

[ -z "$MSG" ] && exit 0

jq -n --arg m "$MSG" '{
  "hookSpecificOutput": {
    "hookEventName": "Stop",
    "additionalContext": ("Loop: " + $m)
  }
}'
