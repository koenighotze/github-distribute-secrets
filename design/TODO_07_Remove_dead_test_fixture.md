# TODO 07 — Remove dead test fixture

**Status:** DONE — merged as PR #38

## Problem

`yamlConfigurationOverwrites` in `internal/config/config_test.go` is identical to
`yamlConfigurationFull` and is never referenced in any test. It is dead code that
creates false familiarity — a reader might assume it tests override behavior.

## Branch

```
refactor/remove-dead-test-fixture
```

## Files

- `internal/config/config_test.go`

---

## TDD note

Deleting an unused constant has no observable behavior. There is no failing test to write.
The verification is `make test` passing before and after — same pattern as the typo fix.

If the constant were referenced, the test suite would fail to compile on deletion — that
compilation failure would serve as the RED phase.

---

## Procedure

### 1. Confirm the constant is unused

```bash
grep -rn "yamlConfigurationOverwrites" .
```

Expected: exactly one result — the declaration itself. Zero usage sites.

### 2. Confirm baseline

```bash
go test ./internal/config/...
```

All green.

### 3. Delete the constant

Remove lines 21–28 from `internal/config/config_test.go`:

```go
// delete this entire block
yamlConfigurationOverwrites = `
common:
   KEY0: VAL0
repo1:
   KEY1: VAL1
repo2:
   KEY2: VAL2
`
```

### 4. Verify

```bash
go test ./internal/config/...
```

Same result as baseline. No compile errors, no test failures.

---

## Verification checklist

- [ ] `grep` confirmed zero usage sites before deleting
- [ ] Baseline green before starting
- [ ] Deleted the declaration only — no other changes
- [ ] `go test ./internal/config/...` green after deletion
- [ ] `make test` still green
