#!/usr/bin/env bash
# PreToolUse hook: blocks Edit/Write on main or master branch.

BRANCH=$(git branch --show-current 2>/dev/null)

if [ "$BRANCH" = "main" ]; then
  jq -n --arg b "$BRANCH" '{
    "hookSpecificOutput": {
      "hookEventName": "PreToolUse",
      "permissionDecision": "deny",
      "permissionDecisionReason": ("On branch " + $b + ". Create a feature branch first.")
    }
  }'
fi
