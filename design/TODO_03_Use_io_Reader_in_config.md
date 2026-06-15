# TODO 03 — Use `io.Reader` in `NewConfigFromReader`

**Status:** DONE — merged as PR #37

## Problem

`NewConfigFromReader` in `internal/config/config.go:44` accepts `*bytes.Reader`
instead of the `io.Reader` interface. The yaml decoder only needs `io.Reader`.
This unnecessarily prevents passing other reader types (files, HTTP responses, strings).

## Branch

```
refactor/io-reader-interface
```

## Files

- `internal/config/config.go`
- `internal/config/config_test.go`

---

## RED — Write the failing test first

Add to `TestNewConfigFromReader` in `internal/config/config_test.go`:

```go
t.Run("should accept any io.Reader, not only bytes.Reader", func(t *testing.T) {
    reader := strings.NewReader(yamlConfigurationCommonOnly)

    result, err := NewConfigFromReader(reader)

    assert.Nil(t, err)
    assert.NotNil(t, result)
})
```

Add `"strings"` to the test file's imports.

### Verify RED

```bash
go test ./internal/config/... 2>&1 | head -20
```

Expected: **compilation error** — `*strings.Reader` cannot be used as `*bytes.Reader`.
That is the RED. Do not proceed until you see it.

---

## GREEN — Widen the parameter type

`internal/config/config.go:44` — change the parameter:

```go
// before
func NewConfigFromReader(reader *bytes.Reader) (config *Configuration, err error) {

// after
func NewConfigFromReader(reader io.Reader) (config *Configuration, err error) {
```

Add `"io"` to imports. The `"bytes"` import stays — it is used elsewhere in the file.

No call-site changes needed. All existing callers pass `*bytes.Reader` which satisfies `io.Reader`.

### Verify GREEN

```bash
go test ./internal/config/...
```

All tests pass.

---

## REFACTOR

No cleanup needed.

---

## Verification checklist

- [ ] Saw compilation failure (`*strings.Reader` incompatible with `*bytes.Reader`)
- [ ] Only changed the parameter type — no other production code altered
- [ ] `"io"` added to imports, `"bytes"` import retained
- [ ] All `./internal/config/...` tests green
- [ ] `make test` green
