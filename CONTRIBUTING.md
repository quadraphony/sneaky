# Contributing to Sneaky Core

Thank you for considering contributing to Sneaky Core. This document explains how to get the repository running and how to submit changes.

Getting started

1. Fork the repository and clone your fork.
2. Install Go 1.26 or compatible (project uses go 1.26).
3. Run `go build ./...` and `go test ./...`.

Development workflow

- Use branches named like `feat/<short-name>` or `fix/<short-name>`.
- Keep commits small and focused. Write descriptive commit messages.
- Include tests for new functionality where appropriate.

Linting and formatting

- Run `gofmt -s -w .` before committing.
- We recommend `golangci-lint` for linting; CI will run the configured checks.

Pull request

- Open a PR against `master` and include a clear description and checklist of changes.
- Link PRs to relevant phase plan items in `docs/04_PHASE_PLAN.md`.
