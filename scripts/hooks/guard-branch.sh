#!/usr/bin/env bash
# PreToolUse hook: blocks Edit/Write on the main or master branch.

git rev-parse --git-dir >/dev/null 2>&1 || exit 0

BRANCH=$(git branch --show-current 2>/dev/null)

if [ "$BRANCH" = "main" ] || [ "$BRANCH" = "master" ]; then
  jq -n --arg b "$BRANCH" '{
    "hookSpecificOutput": {
      "hookEventName": "PreToolUse",
      "permissionDecision": "deny",
      "permissionDecisionReason": ("On branch " + $b + ". Create a feature branch first.")
    }
  }'
fi

exit 0
