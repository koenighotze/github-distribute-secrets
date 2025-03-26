PYTHON_VERSION := pypy3.10
UV := uv
PIP := $(UV) pip
RUN := $(UV) run
VENV := .venv

# Declare all non-file targets as phony to ensure they always run
.PHONY: all install.uv init clean local.setup autoformat lint test qa freeze

# Default target: run all quality assurance checks
all: qa

# Install uv if not present
install.uv:
	@command -v $(UV) >/dev/null 2>&1 || { echo "Installing uv..."; curl -LsSf https://astral.sh/uv/install.sh | sh; }

install.python:
	$(UV) python install $(PYTHON_VERSION)
	$(UV) python pin $(PYTHON_VERSION)

# Initialize the project environment
init:
	$(PIP) install --upgrade pip
	$(PIP) install -r requirements.txt

# Clean up generated files and directories
clean:
	rm -rf 3.10/ $(VENV)/ __pycache__/ bin/ lib/ include/ .pytest_cache/ .coverage .mypy_cache/

# Set up local development environment
local.setup: install.uv
	$(UV) venv
	@echo "Virtual environment created. Activate it with: source $(VENV)/bin/activate"
	@echo "Then run: make init"

# Auto-format code
autoformat: init
	$(RUN) ruff format

lint.code: init
	$(RUN) pydocstyle tests *py

# Check types
lint.types: init
	$(RUN) mypy .

lint.ruff: init
	$(RUN) ruff check

lint.ruff.fix: init
	$(RUN) ruff check --fix

lint.ruff.watch: init
	$(RUN) ruff check --watch

# Run all linters
lint: lint.code lint.types lint.ruff

# Run tests
test: init
	$(RUN) pytest -v tests

# Run all quality assurance checks
qa: lint test

# Freeze dependencies
freeze:
	$(PIP) freeze > requirements.txt

run:
	$(RUN) python main.py config.yml