# TODO 01 — Fix secret error message

**Status:** DONE — merged as PR #36

## Problem

`pkg/onepassword/onepassword_client.go:27` formats the error with `secret` (always
empty at that point) instead of `secretPath`. Every error reads:

```
failed to read secret : <underlying error>
```

The path that failed is invisible.

## Branch

```
fix/secret-error-message
```

## Files

- `pkg/onepassword/onepassword_client.go` (production)
- `pkg/onepassword/onepassword_client_test.go` (test)

---

## RED — Write the failing test first

Add to `TestGetSecret` in `pkg/onepassword/onepassword_client_test.go`:

```go
t.Run("should include the secret path in the error message", func(t *testing.T) {
    client := cliClient{
        runner: createMockOnePasswordCommandRunner(t, nil, errors.New("op failed")),
    }

    _, err := client.GetSecret(testSecretPath)

    assert.ErrorContains(t, err, testSecretPath)
})
```

### Verify RED

```bash
go test -v ./pkg/onepassword/... -run "TestGetSecret/should_include_the_secret_path"
```

Expected failure:

```
Error: "failed to read secret : op failed" does not contain "somepath"
```

If the test passes immediately — the bug is already fixed. Stop and investigate before continuing.

---

## GREEN — Minimal fix

`pkg/onepassword/onepassword_client.go:27` — change `secret` to `secretPath`:

```go
// before
return "", fmt.Errorf("failed to read secret %s: %w", secret, err)

// after
return "", fmt.Errorf("failed to read secret %s: %w", secretPath, err)
```

### Verify GREEN

```bash
go test ./pkg/onepassword/...
```

All tests must pass. No new failures.

---

## REFACTOR

No cleanup needed. One-character fix.

---

## Verification checklist

- [ ] Watched the test fail before touching production code
- [ ] Failure was "does not contain" — not a compile error, not a wrong failure
- [ ] Changed only the format argument
- [ ] All `./pkg/onepassword/...` tests green
- [ ] `make test` still green
