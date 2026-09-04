package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fixture is a trimmed `go test -bench` run over two packages, with the
// non-format lines go test emits around the results and the "#" header
// `just bench-baseline` writes.
const fixture = `# main abc1234
# go1.27.0 darwin/arm64, recorded 2026-09-04T12:00Z
goos: darwin
goarch: arm64
pkg: github.com/joshmedeski/sesh/v2/lister
cpu: Apple M4 Max
BenchmarkApplyDedup/n=10-16     	    2673	      5590 ns/op	    1848 B/op	      14 allocs/op
BenchmarkApplyDedup/n=10-16     	    2662	      3876 ns/op	    1850 B/op	      14 allocs/op
PASS
ok  	github.com/joshmedeski/sesh/v2/lister	0.380s
goos: darwin
goarch: arm64
pkg: github.com/joshmedeski/sesh/v2/picker
cpu: Apple M4 Max
BenchmarkView/n=10-16           	     354	     34880 ns/op	    4520 B/op	     151 allocs/op
PASS
ok  	github.com/joshmedeski/sesh/v2/picker	2.628s
`

func writeRun(t *testing.T, body string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "bench.txt")
	require.NoError(t, os.WriteFile(p, []byte(body), 0o644))
	return p
}

func TestReadRun(t *testing.T) {
	r, err := readRun(writeRun(t, fixture))
	require.NoError(t, err)

	assert.Equal(t, env{goos: "darwin", goarch: "arm64", cpu: "Apple M4 Max"}, r.env)
	assert.Equal(t, []string{"main abc1234", "go1.27.0 darwin/arm64, recorded 2026-09-04T12:00Z"}, r.notes)
	assert.Len(t, r.order, 2, "one entry per benchmark, not per sample")

	dedup := r.benches[r.order[0]]
	assert.Equal(t, "lister ApplyDedup/n=10", dedup.label())
	// ns/op is tidied to sec/op on the way in.
	require.Len(t, dedup.samples["sec/op"], 2)
	assert.InDelta(t, 5.59e-6, dedup.samples["sec/op"][0], 1e-12)
	assert.InDelta(t, 3.876e-6, dedup.samples["sec/op"][1], 1e-12)
	assert.Equal(t, []float64{14, 14}, dedup.samples["allocs/op"])
	assert.Equal(t, []float64{1848, 1850}, dedup.samples["B/op"])

	assert.Equal(t, "picker View/n=10", r.benches[r.order[1]].label())
}

func TestReadRunErrors(t *testing.T) {
	t.Run("missing file reports as not exist", func(t *testing.T) {
		_, err := readRun(filepath.Join(t.TempDir(), "absent.txt"))
		assert.True(t, os.IsNotExist(err), "compare() relies on this to suggest `just bench-baseline`")
	})

	t.Run("a file with no results is an error, not an empty chart", func(t *testing.T) {
		_, err := readRun(writeRun(t, "goos: darwin\nPASS\nok  \tpkg\t0.1s\n"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no benchmark results")
	})
}

func TestTrimProcs(t *testing.T) {
	tests := []struct {
		name, in, want string
	}{
		{"drops the GOMAXPROCS suffix", "ApplyDedup/n=10-16", "ApplyDedup/n=10"},
		{"keeps a name with no suffix", "ApplyDedup", "ApplyDedup"},
		{"keeps a non-numeric suffix", "Filter/query=a-b", "Filter/query=a-b"},
		{"keeps a leading dash", "-16", "-16"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, trimProcs(tt.in))
		})
	}
}

func TestMedian(t *testing.T) {
	tests := []struct {
		name string
		in   []float64
		want float64
	}{
		{"empty", nil, 0},
		{"one sample", []float64{4}, 4},
		{"odd count takes the middle", []float64{9, 1, 5}, 5},
		{"even count averages the two middle", []float64{4, 1, 9, 2}, 3},
		{"an outlier does not drag it", []float64{1, 1, 1, 1000}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, median(tt.in))
		})
	}
}

func TestSpread(t *testing.T) {
	tests := []struct {
		name string
		in   []float64
		want float64
	}{
		{"a single sample says nothing about noise", []float64{5}, 0},
		{"identical samples have no spread", []float64{5, 5, 5}, 0},
		{"half the range over the median", []float64{8, 10, 12}, 0.2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.want, spread(tt.in), 1e-9)
		})
	}
}

func TestBaselinePath(t *testing.T) {
	// Keyed by platform: one file holding a Mac's and a Linux runner's numbers
	// would silently compare them against each other.
	assert.Equal(t, filepath.Join("testdata", "bench", "baseline-darwin-arm64.txt"), baselinePath("darwin", "arm64"))
	assert.NotEqual(t, baselinePath("darwin", "arm64"), baselinePath("linux", "amd64"))
}
