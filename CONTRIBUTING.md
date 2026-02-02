# Contributing to Sesh

Thank you for your interest in contributing to sesh! This project exists because of contributors like you, and we appreciate every pull request, bug report, and feature suggestion.

## Project Vision

Sesh has a specific vision: **stay simple, do one thing well**. We aim to be a focused, reliable terminal session manager that integrates seamlessly with tmux and zoxide.

**Before starting work on new features or significant refactoring**, please open an issue or discussion first. This ensures your contribution aligns with the project's goals and avoids wasted effort on both sides. Bug fixes and documentation improvements can typically proceed directly.

## Development Setup

### Prerequisites

- **Go 1.24+** - [Installation guide](https://golang.org/doc/install)
- **tmux** - Terminal multiplexer
- **zoxide** - Smart directory jumper
- **just** - Command runner ([installation](https://github.com/casey/just#installation))

### Getting Started

1. Fork the repository on GitHub
2. Clone your fork:
   ```sh
   git clone https://github.com/YOUR_USERNAME/sesh.git
   cd sesh
   ```
3. Verify setup by running tests:
   ```sh
   just test
   ```

### Available Commands

```sh
just mock   # Generate mocks (required before testing if interfaces changed)
just test   # Run tests with coverage
just build  # Build to $GOPATH/bin/sesh
```

## Making Changes

1. **Create a branch** for your changes:
   ```sh
   git checkout -b your-feature-name
   ```

2. **Write code** following existing patterns in the codebase

3. **Add tests** for new functionality. Interfaces use [mockery](https://github.com/vektra/mockery) for mock generation - mocks are auto-generated via `just mock`, so you don't need to write them manually.

4. **Run tests** before submitting:
   ```sh
   just test
   ```

## Code Guidelines

- **Interface-based design** - All external dependencies use interfaces for testability (see wrapper packages: `execwrap`, `oswrap`, `pathwrap`)
- **Mocks** - Generated automatically by mockery. Run `just mock` after modifying interfaces.
- **Logging** - Use `log/slog` for structured logging
- **Error handling** - Follow existing patterns in the codebase

## Pull Request Process

1. **Describe your changes clearly** in the PR description
2. **Ensure CI passes** - All tests must be green
3. **Be responsive to feedback** - We may request changes or clarifications

## Reporting Issues

Found a bug or have a feature request? Please [open an issue](https://github.com/joshmedeski/sesh/issues/new) with as much detail as possible.

For bugs, include:
- Steps to reproduce
- Expected vs actual behavior
- Your environment (OS, Go version, tmux version)

---

Thanks again for contributing!
