# TODO 02 — Remove panic for expected failures

**Status:** DONE — merged as PR #39

## Problem

`githubSecretDistribution` uses `panic` for expected operational errors (config file
missing, secret push failure). The function always returns `true` or panics — the
`bool` return type is dead code. The `false` branch in `main.go` is unreachable.

Panics produce stack traces for routine CLI errors, which is confusing for end users.

## Branch

```
fix/panic-to-error
```

## Files

- `cmd/github-distribute-secrets/github_distribute_secrets.go`
- `cmd/github-distribute-secrets/main.go`
- `cmd/github-distribute-secrets/main_test.go`
- `cmd/github-distribute-secrets/github_distribute_secrets_test.go`

---

## RED — Write the failing test first

Add to `TestGithubSecretDistribution` in `github_distribute_secrets_test.go`:

```go
t.Run("should return error if reading config fails", func(t *testing.T) {
    configFileReader := &MockConfigFileReader{expectedError: assert.AnError}

    err := githubSecretDistribution(configFileReader, &MockOnePasswordClient{}, &mockGithubClient{}, false)

    assert.ErrorIs(t, err, assert.AnError)
})

t.Run("should return error if applying config fails", func(t *testing.T) {
    configFileReader := &MockConfigFileReader{
        expectedConfig: &config.Configuration{
            RawConfig:    map[string]config.RepositoryConfiguration{"repo1": {"key": "val"}},
            Repositories: []string{"repo1"},
        },
    }
    githubClient := &mockGithubClient{expectedError: assert.AnError}

    err := githubSecretDistribution(configFileReader, &MockOnePasswordClient{}, githubClient, false)

    assert.Error(t, err)
})
```

### Verify RED

```bash
go test ./cmd/github-distribute-secrets/... 2>&1 | head -20
```

Expected: **compilation error** — `githubSecretDistribution` returns `bool`, not `error`.
That compilation failure is the RED phase. Do not proceed until you see it.

---

## GREEN — Change the signature and fix callers

### 1. `github_distribute_secrets.go`

```go
func githubSecretDistribution(
    configFileReader config.ConfigFileReader,
    op onepassword.OnePasswordClient,
    gh github.GithubClient,
    dumpConfig bool,
) error {
    configuration, err := configFileReader.ReadConfiguration("./config.yml")
    if err != nil {
        return fmt.Errorf("failed to read config file: %w", err)
    }

    if dumpConfig {
        fmt.Println(configuration.DumpConfiguration())
    }

    if !applyConfiguration(configuration, op, gh) {
        return fmt.Errorf("configuration was not applied successfully")
    }

    return nil
}
```

### 2. `main.go`

Update the var type and call site:

```go
var (
    // ...
    myGithubSecretDistribution = githubSecretDistribution
)
```

The var type is inferred — it will update automatically once the function signature changes.

Call site:

```go
if err := myGithubSecretDistribution(myNewConfigFileReader(), op, gh, *dumpConfig); err != nil {
    log.Fatalln(err)
}
```

### 3. `main_test.go`

Update the mock closure type from `func(...) bool` to `func(...) error`:

```go
myGithubSecretDistribution = func(
    configFileReader config.ConfigFileReader,
    op onepassword.OnePasswordClient,
    gh github.GithubClient,
    dumpConfig bool,
) error {
    calledGithubSecretDistribution = true
    return nil
}
```

Update all three places in `main_test.go` where the mock is assigned.

### 4. `github_distribute_secrets_test.go`

Replace `assert.Panics` tests with error-return assertions:

```go
// before
t.Run("should panic even if a single application fails", func(t *testing.T) {
    ...
    assert.Panics(t, func() {
        githubSecretDistribution(configFileReader, onePasswordClient, githubClient, false)
    })
})

// after
t.Run("should return error if a single application fails", func(t *testing.T) {
    ...
    err := githubSecretDistribution(configFileReader, onePasswordClient, githubClient, false)
    assert.Error(t, err)
})
```

Same for the "should panic if reading the config failed" test.

### Verify GREEN

```bash
go test ./cmd/github-distribute-secrets/...
```

All tests pass. No panics, no compile errors.

---

## REFACTOR

Remove the now-unused `"fmt"` import from `github_distribute_secrets.go` if it drops out,
or verify it stays (it's used in `fmt.Println`).

---

## Verification checklist

- [ ] Saw compilation failure before touching production code
- [ ] New tests use `assert.Error` / `assert.ErrorIs`, not `assert.Panics`
- [ ] Existing panic tests replaced (not deleted — replaced with error assertions)
- [ ] `main.go` call site updated — no `log.Fatalln` on `false` branch
- [ ] `make test` green
- [ ] `make build` compiles cleanly
