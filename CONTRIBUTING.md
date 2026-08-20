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
just bench  # Run benchmarks (lister and picker)
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

## Benchmarks

The `lister` and `picker` packages carry benchmarks over the paths whose cost
scales with your data rather than with a constant: listing and deduping N
sessions, and filtering them on every keystroke. Run them with:

```sh
just bench                                      # everything
go test -run='^$' -bench=FilterSessions -benchmem ./picker/   # one of them
```

Benchmarks never run under `-race` — the race detector's overhead is most of
what you would be measuring. `just test` and CI's test job leave `-bench` off
for that reason; a separate CI job runs the benchmarks on a single Linux runner
and posts a [benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat)
comparison against the PR's merge base in the job summary.

When adding a benchmark:

- Run it at N = 10, 100, and 1000 sessions (`benchSizes`), reusing the
  package's `benchSessions` fixture rather than a new one.
- Call `b.ReportAllocs()`. On a shared CI runner, allocations per operation are
  the stable signal; a feature that starts copying the whole list per keystroke
  shows up there first and in the timings later.
- Cover the worst case, not the convenient one. For the picker that is a
  one-character query, which matches nearly everything.

### Regression thresholds

The benchmark job is **informational** — it does not fail a PR yet. We are
watching the run-to-run variance on GitHub's runners before turning a gate on,
because a benchmark gate that cries wolf gets disabled within a month.

When the job does start blocking, these are the thresholds:

| Metric | Threshold | Why |
| --- | --- | --- |
| time/op | more than 25% slower | Generous: shared runners are noisy |
| allocs/op | more than 10% more | Tight: allocation counts are near-deterministic |

Until then, if the benchstat table in your PR shows a regression past those
numbers, treat it as a review comment on your own change: explain it in the PR
or fix it.

### History dashboard

benchstat only ever compares two commits, so it can tell you *whether* a PR
regressed something but never *when* something got slow. Every push to `main`
therefore appends a point to a chart published by
[github-action-benchmark](https://github.com/benchmark-action/github-action-benchmark)
at `https://joshmedeski.github.io/sesh/dev/bench`. If a benchmark comes in more
than 25% slower than the previous point, the action leaves a comment on the
offending commit — it never fails the build, since by then the code is already
on `main`.

The chart takes one sample per benchmark per commit: CI runs `-count=6` for
benchstat's sake, and `.github/scripts/bench-median.py` collapses those to the
median before the data is stored.

The dashboard needs one-time setup from a maintainer, since the action pushes to
`gh-pages` but will not create it. The branch already exists; this is the
recipe, and the empty commit matters — `--orphan` alone leaves an unborn branch
with nothing to push, and it stages the whole tree, so skipping the `git rm`
would commit the entire repo onto `gh-pages`:

```sh
git checkout --orphan gh-pages
git rm -r --cached .                 # unstage; leaves your files on disk
git commit --allow-empty -m "Initialize gh-pages"
git push origin gh-pages
git checkout -f main                 # the files above are untracked now
```

GitHub Pages also has to be enabled for the repository, serving from the
`gh-pages` branch. Until it is, the `Benchmark history` job fails on push to
`main`; the benchstat job on pull requests is unaffected.

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
