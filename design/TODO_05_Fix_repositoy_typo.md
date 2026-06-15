# TODO 05 — Fix `repositoy` typo

**Status:** DONE — merged as PR #40

## Problem

The parameter name `repositoy` (missing an `r`) appears in the `GithubClient` interface,
both implementations, and the `applyConfigurationToRepository` function. Parameter names
in Go interfaces don't affect callers, so this is not a breaking change.

## Branch

```
refactor/fix-repositoy-typo
```

## Files

- `pkg/github/github_client.go`
- `pkg/github/github_dry_run_client.go`
- `cmd/github-distribute-secrets/github_distribute_secrets.go`
- `cmd/github-distribute-secrets/github_distribute_secrets_test.go`

---

## TDD note

This is a pure rename with no behavior change. There is no failing test to write —
the behavior is identical before and after. The "test" is the full test suite
passing after the rename.

**Rule:** Do not rename anything until you have confirmed the full suite is green on
`main`. That green suite is your baseline. Any failure after the rename is a regression.

---

## Procedure

### 1. Confirm baseline

```bash
make test
```

All tests green. Record this as the baseline.

### 2. Rename — do not touch anything else

In every file listed above, replace every occurrence of `repositoy` with `repository`.

Occurrences to rename (parameter names only — not strings or comments):

| File | Location |
|------|----------|
| `pkg/github/github_client.go` | interface method signature (line ~10) |
| `pkg/github/github_client.go` | `cliGithubClient.AddSecretToRepository` param (line ~18) |
| `pkg/github/github_dry_run_client.go` | `dryRunGithubClient.AddSecretToRepository` param (line ~14) |
| `cmd/.../github_distribute_secrets.go` | `applyConfigurationToRepository` param + internal uses |
| `cmd/.../github_distribute_secrets_test.go` | `mockGithubClient.AddSecretToRepository` param |

### 3. Verify nothing broke

```bash
make test
```

Must be identical to baseline. Any new failure is a bug in the rename — fix it before
proceeding.

---

## Verification checklist

- [ ] Baseline `make test` green before starting
- [ ] Used find-replace, not manual edits — avoid missing an occurrence
- [ ] `grep -rn "repositoy" .` returns zero results after rename
- [ ] `make test` green after rename — same result as baseline
- [ ] `make build` compiles cleanly
