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

just bench-baseline  # Record this machine's benchmark baseline (run on main)
just bench-compare   # Chart this worktree against that baseline
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
and compares them against the PR's merge base with
[benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat).

That job compiles both sides up front (`go test -c` at head and at the merge
base) and then **alternates** them, one round each, six rounds. It used to run
the whole suite at head and then the whole suite at the base, which put ten
minutes between the two arms — and a shared runner does not hold still for ten
minutes. Because `-count=6` collects a benchmark's six samples back to back,
each arm was one contiguous block in time, so any drift across the gap
separated the blocks cleanly and every significance test called it real. A pull
request that touched no benchmarked code reported 25 of 96 benchmarks as
"moved", 24 of them in the same direction, with allocations identical
throughout. Alternating spreads drift across both arms instead of confounding
it with the comparison; the same null comparison then reported 1 of 96.

That job leaves its result as a comment on the pull request. It opens with a
one-sentence verdict in a GitHub alert block — the level escalates only past
the thresholds in "Regression thresholds" below — then a chart of what moved,
in a ` ```diff ` fence, which is how a comment gets colour: GitHub renders a
line starting with `-` red and one starting with `+` green, and slower is `+`
and faster is `-`. The benchstat table and the full measurement table are
folded underneath.

Rows are named in English rather than by benchmark name (`FilterSessions/n=1000/plain/empty`
reads as "Filtering 1,000 sessions on a keystroke (nothing typed)"), and ranked
by how much work the change actually costs rather than by percentage — a 50%
swing on a four-nanosecond call is noise wearing a big number, and 8% off half
a millisecond is not. The phrases live in `internal/benchchart/names.go`;
`TestEveryBenchmarkIsNamed` fails when a benchmark is added without one.

The comment is rewritten in place on every push, so a pull request carries one
comment showing the current numbers instead of one per push. The job summary
gets the same body, which is where to look on a pull request from a fork —
`GITHUB_TOKEN` is read-only there, so the comment step skips it.

When adding a benchmark:

- Run it at N = 10, 100, and 1000 sessions (`benchSizes`), reusing the
  package's `benchSessions` fixture rather than a new one.
- Call `b.ReportAllocs()`. On a shared CI runner, allocations per operation are
  the stable signal; a feature that starts copying the whole list per keystroke
  shows up there first and in the timings later.
- Cover the worst case, not the convenient one. For the picker that is a
  one-character query, which matches nearly everything.

### Local before/after chart

CI's benchstat table only exists once a pull request does. For the loop before
that — change something, see whether it cost anything — there is a committed
baseline and a bar chart:

```sh
just bench-baseline   # on main: record this machine's baseline, then commit it
just bench-compare    # in a worktree: run the suite and chart it against that
```

`bench-baseline` writes `testdata/bench/baseline-$GOOS-$GOARCH.txt` — the raw
`go test -bench` output, which `benchstat` can still read directly, prefixed
with `#` comments recording the commit, the Go version, and the flags. Commit
it. `bench-compare` runs the same suite into the gitignored `.bench/` and
charts the difference:

```
time/op  6 of 96 moved · 90 within noise (-all to show) · axis ±25%

  picker FilterSessions/n=10/separator-aware/no-match  1.224µs → 1.127µs        █████████│      -7.9%
  picker FilterSessions/n=1000/plain/empty             46.05µs → 37.13µs   ███████████████│      -19.4%
  picker FilterAliases/n=1000/a                        22.88µs → 23.62µs                  │████  +3.2%
```

Two things about that shape are deliberate:

- **Only what moved.** The suite covers ninety-odd benchmark-and-size
  combinations, and printing all of them buries the handful that changed. A
  benchmark whose delta is inside the two runs' own sample spread is counted
  and dropped. `-all` puts them back, still ordered so the ones that moved come
  first.
- **The bars are deltas, not magnitudes.** Each bar grows out of a centre line
  — left for faster, right for slower — and every bar in a section is measured
  against the same span, named in the header. That is what makes a 19%
  regression look four times worse than a 5% one, which is the whole point of
  drawing this instead of reading the numbers. It also means bar length says
  nothing about whether the benchmark is a nanosecond or a millisecond; the
  columns to the left carry that.

Both recipes use `-benchtime=100ms -count=6`, the same as CI, for the reason the
next section gives: a looser measurement makes the deltas meaningless.

Some caveats the chart itself will remind you of:

- **A baseline belongs to the machine that recorded it.** The file is keyed by
  `GOOS`/`GOARCH` and records the CPU, and the chart warns when the two runs
  disagree on any of them. Someone else's laptop is not a target — record your
  own.
- **A baseline taken on a branch measures the branch**, so the chart then
  reports "nothing moved" no matter what the branch cost. `bench-baseline`
  asks for confirmation when you are not on a clean `main` for that reason.
  Re-record when `main` moves somewhere the chart should be measuring against.
- **A row can be significant and still meaningless.** When a run's own samples
  spread more than 5% — five times the ±1% `-benchtime=100ms` is documented to
  hold these benchmarks to — the chart reports `unstable ±N%` instead of a
  delta, and leaves that row out of the verdict. A rank test only sees the
  order of the samples, so a run drifting from 110µs to 138µs under its own
  measurement separates perfectly from a steady one and scores `p=0.002`
  against it. That is a fact about the machine, not the code.
- **It is not benchstat.** The noise check is a sample-range comparison, not
  the significance test benchstat runs, and ninety-six simultaneous comparisons
  will turn up the odd false positive whatever the test. Broad strokes: use it
  to notice something, and let the benchstat table and `allocs/op` settle it.

By default the chart plots `time/op` and `allocs/op` — allocations are the more
trustworthy of the two when they disagree, per the section below, and a chart
showing only time invites exactly that misreading. `-metric` picks:

```sh
go run ./internal/benchchart -metric allocs         # allocations alone
go run ./internal/benchchart -metric time,allocs,bytes
go run ./internal/benchchart -all -width 140        # everything, wider
go run ./internal/benchchart -raw-names             # benchmark names, for -bench=
go run ./internal/benchchart -format markdown       # what CI comments

# Re-chart the last run without re-running it. Colour is dropped when stdout
# is not a terminal, so a pager needs CLICOLOR_FORCE.
CLICOLOR_FORCE=1 go run ./internal/benchchart | less -R
```

There are two renderers, and they are deliberately not the same. A terminal has
ANSI colour, two hundred columns and eighth-width block glyphs; a comment has a
`diff` fence, emoji and about a hundred. Making one output satisfy both would
mean settling for what they share. `-raw-names` is the escape hatch when you
want the name to paste into `-bench=`.

### Regression thresholds

The benchmark job is **informational** — it does not fail a PR yet. We are
measuring the run-to-run variance on GitHub's runners before turning a gate on,
because a benchmark gate that cries wolf gets disabled within a month.

When the job does start blocking, these are the thresholds:

| Metric | Threshold | Why |
| --- | --- | --- |
| time/op | more than 25% slower | Generous: shared runners are noisy |
| allocs/op | more than 10% more | Tight: allocation counts are near-deterministic |

Until then, if the benchstat table in your PR shows a regression past those
numbers, treat it as a review comment on your own change: explain it in the PR
or fix it.

Trust `allocs/op` over `time/op` when the two disagree. A timing change with an
unchanged allocation count, in a package the branch never touched, is noise
until something explains it; that pattern accounted for every alert the history
dashboard raised in its first three commits.

Before tightening either threshold, measure what the runner can actually
resolve. The **Benchmark calibration** workflow
(`.github/workflows/bench-calibrate.yml`) runs the suite twice on the same
commit and benchstats the two halves, so every delta it reports is noise by
construction — a threshold tighter than its widest row cannot tell a regression
from a re-run. Start it by pushing a `bench-calibrate/**` branch, or from the
Actions tab.

One lesson from setting these numbers, worth keeping in mind when adding a
`-benchtime`: it takes a duration, and `-benchtime=100x` instead pins each
benchmark to 100 iterations. On a nanosecond-scale benchmark that is
microseconds of total work, too short to reach steady state, and it reports
warmup rather than the code — spreads of ±70% to ±125% and absolute numbers
almost 3x high, even on an idle machine. `-benchtime=100ms` holds the same
benchmarks to ±1%. No threshold is meaningful under a measurement looser than
the threshold itself.

### History dashboard

benchstat only ever compares two commits, so it can tell you *whether* a PR
regressed something but never *when* something got slow. Every push to `main`
therefore appends a point to a chart published by
[github-action-benchmark](https://github.com/benchmark-action/github-action-benchmark)
at `https://joshmedeski.github.io/sesh/dev/bench`. If a benchmark comes in more
than 25% slower than the previous point, the action says so in the job summary —
it never fails the build, since by then the code is already on `main`.

Read that alert as "look at the chart", not as a verdict. It compares a single
point against the single point before it with a fixed ratio and no model of
variance, so it flags runner noise, and it ratchets: one lucky-fast commit
becomes the baseline that makes the next ordinary commit look like a
regression. It used to comment on the commit it flagged, which put three
permanent regression notices on `main` — one of them on a commit that changed
nothing but a CI YAML file. Use the chart for the trend and the PR's benchstat
table, which reports significance, for the verdict.

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
