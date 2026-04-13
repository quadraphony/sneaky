# CLI Specification

## Purpose

The CLI exists to test and validate Sneaky Core independently of any frontend.

## Command Set

### `sneakycli validate <config-path>`
Validates:
- file existence
- file readability
- config detection
- config validity

### `sneakycli start <config-path>`
Starts the detected adapter using the provided config.

### `sneakycli status`
Reports current runtime state if the process model supports it.

### `sneakycli version`
Prints build and version information.

## CLI Output Rules

- output must be clear
- validation failures must explain why detection failed
- runtime failures must not be hidden
- exit codes must be meaningful
- avoid noisy unreadable logs unless debug mode is explicitly enabled

## CLI Is Not the Product

The CLI is a development and verification tool.
It must not become the center of the architecture.
