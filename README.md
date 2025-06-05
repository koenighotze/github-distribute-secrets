# Secret distributor

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/f90eeb7872aa48d587f95a5375a35bed)](https://app.codacy.com/gh/koenighotze/github-distribute-secrets/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Codacy Coverage Badge](https://app.codacy.com/project/badge/Coverage/f90eeb7872aa48d587f95a5375a35bed)](https://app.codacy.com/gh/koenighotze/github-distribute-secrets/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)
[![Build](https://github.com/koenighotze/github-distribute-secrets/actions/workflows/build.yml/badge.svg)](https://github.com/koenighotze/github-distribute-secrets/actions/workflows/build.yml)

Run `./github-distribute-secrets` to apply the secrets to the repositories. Or using the "scripted" version use `make run.local`.

## Project Structure

The project follows the standard Go project layout:

- `cmd/github-distribute-secrets/`: Main application code
- `internal/`: Internal packages not meant for external use
  - `config/`: Configuration handling
  - `github/`: GitHub API client
  - `onepassword/`: 1Password integration
- `scripts/`: Utility scripts

To build the project, run:

```bash
make build
```

For development, you can use:

```bash
make run.local
```

## Configuration

See [config.yml](./config.yml) for details on how to configure the secrets distribution.

```yaml
# Common secrets shared across multiple projects or environments.
common:
  name-of-the-secret: reference-to-the-1password-value

reposiotory-name:
  name-of-the-secret: reference-to-the-1password-value
```

## TODOS

- extract 1password and github into real go modules
- Replace log.Default() with a structured logging library like zerolog or zap
- Add timeouts for external commands
- Add version information to builds
- Add progress indicators during secret distribution
- Add confirmation question
- Add integration tests
