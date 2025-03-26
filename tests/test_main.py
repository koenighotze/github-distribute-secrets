# -*- coding: utf-8 -*-
from pathlib import Path
from unittest.mock import MagicMock, mock_open, patch

import pytest

from main import get_1password_secret, load_config, process_secrets, set_github_secret


class TestMain:
    @patch(
        "builtins.open",
        new_callable=mock_open,
        read_data="common:\n  SECRET: op://vault/item/field\nrepo1:\n  SECRET1: op://vault/item1/field1",
    )
    def test_load_config(self, mock_file):
        config_path = Path("config.yaml")
        result = load_config(config_path)
        assert result == {
            "common": {"SECRET": "op://vault/item/field"},
            "repo1": {"SECRET1": "op://vault/item1/field1"},
        }

    @patch("main.subprocess.run")
    def test_get_1password_secret(self, mock_run):
        mock_run.return_value = MagicMock(stdout="secret_value")
        result = get_1password_secret("op://vault/item/field")
        assert result == "secret_value"
        mock_run.assert_called_once_with(
            ["op", "read", "op://vault/item/field"],
            check=True,
            capture_output=True,
            text=True,
        )

    @patch("main.subprocess.run")
    def test_set_github_secret(self, mock_run):
        set_github_secret("repo1", "SECRET_NAME", "secret_value")
        mock_run.assert_called_once_with(
            ["gh", "secret", "set", "SECRET_NAME", "--repo", "koenighotze/repo1"],
            input="secret_value",
            text=True,
            check=True,
        )

    @patch("main.get_1password_secret", return_value="secret_value")
    @patch("main.set_github_secret")
    def test_process_secrets(self, mock_set_secret, mock_get_secret):
        process_secrets(
            "repo1",
            {"SECRET1": "op://vault/item1/field1"},
            {"COMMON_SECRET": "op://vault/common/field"},
        )
        mock_get_secret.assert_any_call("op://vault/common/field")
        mock_get_secret.assert_any_call("op://vault/item1/field1")
        mock_set_secret.assert_any_call("repo1", "COMMON_SECRET", "secret_value")
        mock_set_secret.assert_any_call("repo1", "SECRET1", "secret_value")


if __name__ == "__main__":
    pytest.main()
