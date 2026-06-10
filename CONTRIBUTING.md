# Contributing

Thank you for your interest in contributing to myip!

## Getting Started

1. Fork the repository and clone your fork.
2. Install Go (see `go.mod` for the required version).
3. Run the tests to verify your setup:
   ```sh
   go test ./...
   ```

## Making Changes

- Open an issue to discuss significant changes before starting work.
- Keep pull requests focused on a single concern.
- Add tests for any new or changed behavior.
- Run `go vet ./...` and `go test ./...` locally before submitting.

## Pull Request Process

1. Create a branch from `main` with a descriptive name (e.g. `fix/issue-description`).
2. Write a clear PR title and description explaining the motivation and scope.
3. Ensure all CI checks pass.
4. A maintainer will review and merge your PR.

## Reporting Bugs

Please open a [GitHub Issue](https://github.com/kitsuyui/myip/issues/new/choose) with:
- A clear description of the problem
- Steps to reproduce
- Expected vs. actual behavior
- Environment details (OS, Go version, myip version)

## Security Vulnerabilities

See [SECURITY.md](SECURITY.md) for how to report security issues privately.

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md).
