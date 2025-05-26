# Secret distributor

[![Build](https://github.com/koenighotze/github-distribute-secrets/actions/workflows/build.yml/badge.svg)](https://github.com/koenighotze/github-distribute-secrets/actions/workflows/build.yml)

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
## TODO DUMP

chat.promptFiles (Experimental): Enable or disable reusable prompt files.
chat.promptFilesLocations (Experimental): Specify the location of prompt files. Set to true to use the default location (.github/prompts), or use the { "/path/to/folder": boolean } notation to specify a different path. Relative paths are resolved from the root folder(s) of your workspace.
