#!/usr/bin/env bash
# Stop hook: injects the next loop step into Claude's context for the next turn.

git rev-parse --git-dir >/dev/null 2>&1 || exit 0

BRANCH=$(git branch --show-current 2>/dev/null)
DIRTY=$(git status --porcelain 2>/dev/null | grep -cv '^??')
AHEAD=$(git log '@{u}..' --oneline 2>/dev/null | wc -l | tr -d ' ')

MSG=""

if [ "$BRANCH" = "main" ]; then
  MSG="On main branch — create a feature branch before editing files."
elif [ "$DIRTY" -gt 0 ] && [ "$AHEAD" -eq 0 ]; then
  MSG="Uncommitted changes on '$BRANCH'. Next: /superpowers:verification-before-completion → /commit-commands:commit-push-pr"
elif [ "$AHEAD" -gt 0 ]; then
  MSG="Unpushed commits on '$BRANCH'. Next: /commit-commands:commit-push-pr → /code-review"
elif PR_NUM=$(timeout 5 gh pr list --head "$BRANCH" --state open --json number --jq '.[0].number' 2>/dev/null) && [ -n "$PR_NUM" ]; then
  MSG="Open PR #$PR_NUM on '$BRANCH'. Next: /watch-pr to monitor CI and reviews"
fi

[ -z "$MSG" ] && exit 0

jq -n --arg m "$MSG" '{
  "hookSpecificOutput": {
    "hookEventName": "Stop",
    "additionalContext": ("Loop: " + $m)
  }
}'
