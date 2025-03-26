# -*- coding: utf-8 -*-
import pytest


class TestMain:
    def test_equals_false(self) -> None:
        assert 1 != 2 - 1


if __name__ == "__main__":
    pytest.main()
