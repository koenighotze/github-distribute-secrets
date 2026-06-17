# Ship — Full Delivery Loop

Commit staged work, push, open a PR, then run the full CI + review loop until the PR is green and all findings are addressed or dismissed.

## Phase 1 — Commit & PR

### 1.1 Check state

```bash
git status --porcelain | grep -v '^??'
git branch --show-current
```

If no tracked changes and no unpushed commits, check for an already-open PR and skip to Phase 2.

### 1.2 Stage and commit

```bash
git add -u
```

Write a conventional commit message (`feat:`, `fix:`, `chore:`, etc.) that describes the why, not the what.

```bash
git commit -m "<type>: <message>"
```

### 1.3 Push

```bash
git push -u origin HEAD
```

### 1.4 Create PR (or detect existing)

```bash
PR_NUM=$(gh pr list --head "$(git branch --show-current)" --state open --json number --jq '.[0].number' 2>/dev/null)
```

If `PR_NUM` is empty, create one:

```bash
gh pr create --title "<title>" --body "<summary + test plan>"
PR_NUM=$(gh pr list --head "$(git branch --show-current)" --state open --json number --jq '.[0].number')
```

Capture `PR_NUM` — it is used in every subsequent step.

---

## Phase 2 — CI Loop

Repeat until all checks pass.

### 2.1 Wait for checks

```bash
gh pr checks "$PR_NUM" --watch --interval 15
```

### 2.2 Handle failures

For each failed check:

1. Find the run ID: `gh pr checks "$PR_NUM" --json name,link`
2. Fetch failure logs: `gh run view <run-id> --log-failed`
3. Diagnose root cause
4. Fix the code
5. `make test` — must pass locally before pushing
6. Commit and push:
   ```bash
   git add -u
   git commit -m "fix: <what was wrong>"
   git push
   ```
7. Go back to 2.1

---

## Phase 3 — Review Loop

Run once CI is fully green.

### 3.1 Fetch inline comments (Copilot, reviewers)

```bash
REPO=$(gh repo view --json nameWithOwner --jq '.nameWithOwner')
gh api "repos/$REPO/pulls/$PR_NUM/comments" --jq '.[] | {path: .path, line: .line, body: .body}'
```

### 3.2 Fetch review summaries

```bash
gh pr view "$PR_NUM" --json reviews --jq '.reviews[] | {author: .author.login, state: .state, body: .body}'
```

### 3.3 Evaluate each finding

For each comment:
- **Warranted**: apply fix, run `make test`
- **Not warranted**: note the reasoning explicitly — do not apply blindly

### 3.4 If fixes were made

```bash
git add -u
git commit -m "fix: address review findings"
git push
```

Then go back to Phase 2 to confirm CI still passes with the fixes.

---

## Done

Report the PR as ready to merge when:
- All CI checks are green
- All review findings are addressed or explicitly dismissed with reasoning

```bash
gh pr view "$PR_NUM" --json url,title,state
```
