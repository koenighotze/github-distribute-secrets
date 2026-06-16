# TODO 06 — Drop `log.Default()` calls

**Status:** DONE — merged as PR #41

## Problem

`log.Default().Printf(...)`, `log.Default().Fatalln(...)`, and `log.Default().Panicf(...)`
are used throughout the codebase. `log.Default()` returns the same default logger as the
package-level functions — it is redundant noise.

## Branch

```
refactor/drop-log-default
```

## Files

- `cmd/github-distribute-secrets/github_distribute_secrets.go`
- `pkg/github/github_client.go`
- `cmd/github-distribute-secrets/main.go`

---

## TDD note

This is a pure cosmetic substitution with no behavior change. No failing test to write.
The verification is `make test` passing before and after — same as the typo fix.

---

## Procedure

### 1. Confirm baseline

```bash
make test
```

All green.

### 2. Replace — do not touch anything else

In each file, replace `log.Default().` with nothing (i.e., use the package-level call directly):

| Before | After |
|--------|-------|
| `log.Default().Printf(...)` | `log.Printf(...)` |
| `log.Default().Fatalln(...)` | `log.Fatalln(...)` |
| `log.Default().Panicf(...)` | `log.Panicf(...)` |

Occurrences by file:

**`github_distribute_secrets.go`**
- `log.Default().Printf("Cannot apply config to repository %s successfully!", ...)` → `log.Printf(...)`

**`github_client.go`**
- `log.Default().Printf("In repository %s. Adding secret with key %s", ...)` → `log.Printf(...)`

**`main.go`**
- `log.Default().Fatalln(...)` → `log.Fatalln(...)`

### 3. Verify nothing broke

```bash
make test
```

Identical to baseline.

---

## Verification checklist

- [ ] Baseline `make test` green before starting
- [ ] `grep -rn "log.Default()" .` returns zero results after change
- [ ] `make test` green after change
- [ ] `make build` compiles cleanly
