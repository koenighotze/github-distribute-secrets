# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
make build          # compile binary
make test           # unit tests + coverage
make test.all       # unit + integration tests (requires gh and op CLIs)
make lint           # golangci-lint
make vet            # go fmt + go vet
make run.local      # run without compiling

# Single test
go test -v ./pkg/github/... -run TestAddSecretToRepository
go test -v ./cmd/github-distribute-secrets/... -run TestApplyConfiguration

# Coverage HTML
make test.report
```

Integration tests require working `gh` (GitHub CLI) and `op` (1Password CLI) installs. Tag: `//go:build integration`.

## Non negotiable standards

- No pull request shoul have more than 300 lines of code changes. If you need more, split into multiple PRs and consider refactoring.
- Try hard to keep files below 1000 lines. If you need more, consider refactoring.
- Avoid spaghetti code. If you need more, consider refactoring.
- Always go for simpler, human-readable code!

## Architecture

The tool reads `config.yml`, fetches secret values from 1Password, and writes them as GitHub repository secrets via CLI subprocesses — no direct API clients.

**Data flow:**

```
config.yml
  └─ internal/config → Configuration{common + per-repo maps}
       └─ cmd/.../github_distribute_secrets.go → applyConfiguration()
            ├─ pkg/onepassword → "op read <path>"  (results cached in-memory)
            └─ pkg/github     → "gh secret set <key> --body <val> --repo <repo>"
```

**Key design choices:**

- `pkg/cli.CommandRunner` is the single abstraction over `os/exec`. Everything else depends on it, which is what makes unit testing possible without real CLI tools.
- `pkg/github` has two implementations of `GithubClient`: the real one and `dryRunGithubClient` (validates repos exist but skips writes). Selected at startup via `--dry-run`.
- `internal/config.Configuration.GetConfigurationForRepository()` merges `common` secrets with repo-specific overrides. Repo-specific keys win.
- `main.go` uses package-level function-pointer variables (`myNewGhClient`, `myNewOpClient`, etc.) for dependency injection in tests — not interfaces at the `main` level.

**Runtime requirements:** `gh` and `op` must be authenticated before running.

## config.yml format

```yaml
common:                    # applied to every repository
  KEY: op://vault/item/field

owner/repo:                # repo-specific; merged with common (overrides common on conflict)
  EXTRA_KEY: op://vault/item/field
```

The `common` key is reserved; it is never treated as a repository name.
