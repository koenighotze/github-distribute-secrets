# Watch PR — CI and Review Loop

After opening a PR, wait for CI and Copilot/reviewer feedback, fix warranted issues,
and push until the PR is green and all review findings are addressed.

## Step 1 — Identify the PR

```bash
BRANCH=$(git branch --show-current)
gh pr list --head "$BRANCH" --state open --json number,url
```

If no open PR is found, stop and tell the user.

## Step 2 — Wait for CI

```bash
gh pr checks <number> --watch --interval 15
```

## Step 3 — Handle CI failures

For each failed check:

1. Fetch the failure output from the check run logs via `gh run view`
2. Diagnose the root cause
3. Fix the code
4. Run `make test` to verify locally
5. Push — then go back to Step 2

## Step 4 — Fetch review feedback

```bash
gh pr view <number> --json reviews,comments
```

For each finding:
- Evaluate whether it is correct and warranted
- If yes: apply the fix, run `make test`
- If no: note the reasoning — do not blindly apply every suggestion

## Step 5 — Push and re-check

If any fixes were made:

```bash
git add -u
git commit -m "fix: address CI failures / review findings"
git push
```

Then go back to Step 2 to confirm CI still passes.

## Done

Report the PR as ready to merge when CI is green and all review findings are addressed or
explicitly dismissed with reasoning.
