# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Sesh is a smart terminal session manager written in Go that helps users create and manage tmux sessions quickly and easily using zoxide. It's a CLI tool that integrates with tmux and zoxide to provide intelligent session management.

## Core Architecture

- **Module**: `github.com/joshmedeski/sesh/v2`
- **Go Version**: 1.24.2 (toolchain 1.24.4)
- **Main Entry Point**: `main.go` → `seshcli.App()`

### Key Packages

- `seshcli/` - CLI commands and dependency injection (`deps.go` wires all dependencies)
- `lister/` - Lists sessions from various sources (tmux, zoxide, config, tmuxinator)
- `connector/` - Handles session connections with type-specific strategies
- `namer/` - Session naming with strategy pattern (git bare → git → directory fallback)
- `configurator/` - Configuration management (`sesh.toml` in `$XDG_CONFIG_HOME/sesh`)
- `tmux/`, `zoxide/`, `tmuxinator/` - External tool integrations
- `model/` - Shared data structures (Config, SeshSession, ConnectOpts)

### Dependency Injection Pattern

Dependencies are wired in `seshcli/deps.go`:
- `BaseDeps` - Config-free dependencies (exec, os, path, shell wrappers)
- `Deps` - Config-dependent dependencies (lister, connector, namer, etc.)
- `NewBaseDeps()` → `BuildAll(configPath)` creates the full dependency graph

### Naming Strategy

The `namer` package determines session names using a strategy chain:
1. Git bare repository name (looks for `.bare` folder)
2. Git repository name (from remote or directory)
3. Directory name (respects `dir_length` config)

## Common Development Commands

### Using justfile (preferred)
```bash
just build          # Build to $GOPATH/bin/sesh
just build v2.0.0   # Build with specific version
just test           # Run tests with coverage (regenerates mocks first)
just mock           # Generate mocks only
```

### Direct commands
```bash
go test -run TestFunctionName ./package/...  # Run single test
go test -cover -race ./namer/...             # Test specific package
```

### Generate mocks
Mocks are configured in `.mockery.yaml` and placed alongside interfaces:
```bash
just mock
# Or: GOFLAGS="-buildvcs=false" mockery
```

## Configuration

- Config file: `$XDG_CONFIG_HOME/sesh/sesh.toml` or `~/.config/sesh/sesh.toml`
- Custom path: `sesh --config /path/to/sesh.toml <command>`
- Test config examples: `configurator/testdata/`

## Development Notes

- All external dependencies use interfaces for testability (see wrapper packages: `execwrap`, `oswrap`, `pathwrap`)
- Mock files follow pattern `mock_*.go` in the same package as the interface
- Follow existing patterns for error handling and `slog` logging
