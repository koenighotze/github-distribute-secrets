# Secret distributor

[![QA](https://github.com/koenighotze/github-distribute-secrets/actions/workflows/qa.yml/badge.svg)](https://github.com/koenighotze/github-distribute-secrets/actions/workflows/qa.yml)
[![Test](https://github.com/koenighotze/github-distribute-secrets/actions/workflows/test.yml/badge.svg)](https://github.com/koenighotze/github-distribute-secrets/actions/workflows/test.yml)

Run `scripts/apply-configuration.sh` to apply the secrets to the repositories.

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
