# Secret distributor

[![QA](https://github.com/koenighotze/github-distribute-secrets/actions/workflows/qa.yml/badge.svg)](https://github.com/koenighotze/github-distribute-secrets/actions/workflows/qa.yml)
[![Test](https://github.com/koenighotze/github-distribute-secrets/actions/workflows/test.yml/badge.svg)](https://github.com/koenighotze/github-distribute-secrets/actions/workflows/test.yml)

Run `scripts/apply-configuration.sh` to apply the secrets to the repositories.

## Configuration

See [config.yml](./config.yml) for details on how to configure the secrets distribution.

```yaml
# Common secrets shared across multiple projects or environments.
common:
  name-of-the-secret: reference-to-the-1password-value

reposiotory-name:
  name-of-the-secret: reference-to-the-1password-value
```
