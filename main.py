# pylint: disable=invalid-name,missing-docstring,missing-function-docstring,logging-fstring-interpolation
import logging
import subprocess
import sys
from pathlib import Path
from typing import Dict

import yaml

logging.basicConfig(
    level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s"
)


def load_config(config_path: Path) -> Dict[str, Dict[str, str]]:
    try:
        with open(config_path, "r", encoding="utf-8") as file:
            config = yaml.safe_load(file)
            if not isinstance(config, dict):
                raise ValueError(
                    "Configuration file must contain a dictionary at the top level."
                )
            return config
    except FileNotFoundError:
        logging.error("Configuration file %s not found.", config_path)
        sys.exit(1)
    except yaml.YAMLError as e:
        logging.error(f"Error parsing YAML file {config_path}: {e}")
        sys.exit(1)
    except ValueError as e:
        logging.error(f"Invalid configuration format: {e}")
        sys.exit(1)


def get_1password_secret(op_reference: str) -> str:
    try:
        logging.info(f"Reading 1Password secret {op_reference}...")
        result = subprocess.run(
            ["op", "read", op_reference], check=True, capture_output=True, text=True
        )
        return result.stdout.strip()
    except subprocess.CalledProcessError as e:
        logging.error(
            "Failed to read 1Password secret %s. Command: %s. Error: %s",
            op_reference,
            e.cmd,
            e.stderr.strip(),
        )
        sys.exit(1)


def set_github_secret(repo: str, secret_name: str, secret_value: str) -> None:
    try:
        logging.info(f"Setting GitHub secret {secret_name} in {repo}...")
        subprocess.run(
            ["gh", "secret", "set", secret_name, "--repo", f"koenighotze/{repo}"],
            input=secret_value,
            text=True,
            check=True,
        )
        logging.info(f"Set secret {secret_name} in {repo}")
    except subprocess.CalledProcessError as e:
        logging.error(
            f"Failed to set GitHub secret {secret_name} in {repo}. Command: {e.cmd}. Error: {e.stderr.strip()}"
        )
        sys.exit(1)


def exit_with_error(message: str) -> None:
    """Log an error message and exit the program."""
    logging.error(message)
    sys.exit(1)


def process_secrets(
    repo: str, secrets: Dict[str, str], common_secrets: Dict[str, str]
) -> None:
    """Process and set secrets for a given repository."""
    print("Processing secrets for repo", repo)

    secrets = secrets or {}  # Initialize secrets to an empty dictionary if None
    common_secrets = (
        common_secrets or {}
    )  # Initialize secrets to an empty dictionary if None
    combined_secrets = {**common_secrets, **secrets}
    if not combined_secrets:  # Skip processing if no secrets are defined
        logging.info("No secrets to process for repo '%s'. Skipping...", repo)
        return
    for secret_name, op_reference in combined_secrets.items():
        logging.info("Processing secret '%s' for repo '%s'...", secret_name, repo)
        secret_value = get_1password_secret(op_reference)
        set_github_secret(repo, secret_name, secret_value)


def main(config_path: Path) -> None:
    config = load_config(config_path)
    common_secrets = config.pop("common", {})  # Extract "common" section if it exists

    for repo, secrets in config.items():
        process_secrets(repo, secrets, common_secrets)


if __name__ == "__main__":
    if len(sys.argv) != 2:
        exit_with_error(f"Usage: {sys.argv[0]} <config.yaml>")

    config_file = Path(sys.argv[1])
    if not config_file.exists():
        exit_with_error(f"Configuration file {config_file} does not exist.")

    main(config_file)
