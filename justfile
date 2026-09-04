# Generate mocks
# Stale mocks are cleared first: mockery cannot regenerate a mock for a package
# that no longer compiles, so a mock left over from before a signature change
# blocks its own replacement.
mock:
    find . -name 'mock_*.go' -delete
    GOFLAGS="-buildvcs=false" go tool mockery

# Run tests with coverage
test: mock
    go test -cover -race ./... -coverprofile=coverage.out

# The packages that carry benchmarks, and the parameters the numbers are only
# meaningful under: -benchtime=100ms is where these benchmarks hold to +-1%,
# and -count=6 is what benchstat wants. Both match CI. See CONTRIBUTING.md.
bench_pkgs := "./lister/... ./picker/..."
bench_flags := "-run='^$' -bench=. -benchmem -benchtime=100ms -count=6"

# Run benchmarks (no -race: it measures the race detector, not the code)
bench: mock
    go test -run='^$' -bench=. -benchmem {{bench_pkgs}}

# Run the suite into <out>, prefixed with a header recording what produced it.
# The header is "#" comments rather than benchmark-format `key: value` lines on
# purpose: a config key present in one file and absent from the other changes
# how benchstat groups them, and these files are meant to stay something
# benchstat can still read directly.
_bench-run out: mock
    #!/usr/bin/env bash
    set -euo pipefail
    mkdir -p "$(dirname '{{out}}')"
    dirty=""
    git diff --quiet HEAD -- . || dirty=" (uncommitted changes)"
    {
      echo "# $(git rev-parse --abbrev-ref HEAD) $(git rev-parse --short HEAD)${dirty}"
      echo "# $(go version | cut -d\  -f3-4), recorded $(date -u +%Y-%m-%dT%H:%MZ)"
      echo "# go test {{bench_flags}}"
    } > '{{out}}'
    # Redirected rather than tee'd: the wall of ns/op is what the chart exists
    # to replace. On a failure it is the only diagnostic there is, so print it.
    echo "running benchmarks into {{out}} (a few minutes)..."
    if ! go test {{bench_flags}} {{bench_pkgs}} >> '{{out}}' 2>&1; then
      cat '{{out}}' >&2
      exit 1
    fi

# Record this machine's benchmark baseline (run on main, then commit the file:
# a baseline taken on a branch measures the branch)
bench-baseline:
    #!/usr/bin/env bash
    set -euo pipefail
    # Every chart is read against this file, and a baseline recorded on a
    # branch quietly measures that branch: the comparison then reports
    # "nothing moved" however much the branch actually cost. Easy to do by
    # reflex from the worktree you are already in, so ask first.
    branch="$(git rev-parse --abbrev-ref HEAD)"
    dirty=""
    git diff --quiet HEAD -- . || dirty=" with uncommitted changes"
    if [ "$branch" != "main" ] || [ -n "$dirty" ]; then
      echo "This would record the baseline from '${branch}'${dirty}, not from main." >&2
      echo "The chart would then be measuring these changes against themselves." >&2
      read -r -p "record anyway? [y/N] " reply
      case "$reply" in [Yy]*) ;; *) echo "aborted" >&2; exit 1;; esac
    fi
    out="testdata/bench/baseline-$(go env GOOS)-$(go env GOARCH).txt"
    just _bench-run "$out"
    echo "wrote $out - commit it"

# Chart this worktree against the committed baseline for this machine.
bench-compare: (_bench-run ".bench/current.txt")
    go run ./internal/benchchart

# Build sesh binary to GOPATH/bin
build version="dev":
    go build -buildvcs=false -ldflags "-X 'main.version={{version}}'" -o `go env GOPATH`/bin/sesh

# Generate man page
man: build
    mkdir -p share/man/man1
    sesh man > share/man/man1/sesh.1
