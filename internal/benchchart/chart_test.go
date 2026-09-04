package main

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func pair(base, head []float64) row { return row{label: "x", base: base, head: head} }

// steady is a run whose samples agree, so a comparison against another steady
// run turns entirely on the difference between the two.
func steady(v float64) []float64 { return []float64{v, v, v} }

func TestAxis(t *testing.T) {
	// The bar is measured against the section's span rather than against
	// itself, which is what lets two rows be compared to each other.
	const half, scale = 10, 0.5

	tests := []struct {
		name              string
		delta             float64
		both              bool
		wantLeft, wantRig int
	}{
		{name: "the span fills one side", delta: 0.5, both: true, wantRig: 10},
		{name: "half the span fills half of it", delta: 0.25, both: true, wantRig: 5},
		{name: "a fifth of the span is a fifth of the bar", delta: 0.1, both: true, wantRig: 2},
		{name: "an improvement grows the other way", delta: -0.25, both: true, wantLeft: 5},
		{name: "no change draws no bar", delta: 0, both: true},
		{name: "past the span is clamped", delta: 5, both: true, wantRig: 10},
		{name: "a benchmark only one run has draws no bar", delta: 0.5, both: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ansi.Strip(axis(line{delta: tt.delta, both: tt.both, style: worse}, scale, half))
			left, right, ok := strings.Cut(got, "│")
			require.True(t, ok, "every row carries the centre line, so they stack into a spine")

			assert.Equal(t, half, len([]rune(left)), "the centre stays in the same column on every row")
			assert.Equal(t, half, len([]rune(right)))
			assert.Equal(t, tt.wantLeft, strings.Count(left, "█"))
			assert.Equal(t, tt.wantRig, strings.Count(right, "█"))
			if tt.wantLeft > 0 {
				assert.True(t, strings.HasSuffix(left, "█"), "an improvement grows out of the centre, not the edge")
			}
		})
	}
}

func TestNiceScale(t *testing.T) {
	tests := []struct {
		name       string
		peak, want float64
	}{
		{"a small peak still gets a floor", 0.001, 0.05},
		{"rounds up to the next round number", 0.09, 0.1},
		{"an exact step is its own scale", 0.25, 0.25},
		{"just past a step takes the next", 0.26, 0.5},
		{"a peak past every step is its own scale", 25, 25},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, niceScale(tt.peak))
		})
	}
}

func TestPercent(t *testing.T) {
	assert.Equal(t, "50%", percent(0.5))
	assert.Equal(t, "5%", percent(0.05))
	assert.Equal(t, "12.5%", percent(0.125))
}

func TestDelta(t *testing.T) {
	tests := []struct {
		name string
		row  row
		want string
	}{
		{"a regression is signed", pair(steady(100), steady(150)), "+50.0%"},
		{"an improvement is signed", pair(steady(100), steady(50)), "-50.0%"},
		{"a move inside the noise floor is not a number", pair(steady(100), steady(101)), "~"},
		{"a move inside the runs' own spread is not a number", row{base: []float64{50, 100, 150}, head: steady(130)}, "~"},
		{"no baseline means the branch added it", pair(nil, steady(10)), "new"},
		{"no current run means the branch removed it", pair(steady(10), nil), "gone"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, delta(tt.row))
		})
	}
}

func TestStyle(t *testing.T) {
	// Every metric the chart plots is lower-is-better, which is what lets it
	// colour a delta without knowing which one it is looking at.
	assert.Equal(t, better, style(pair(steady(100), steady(50))))
	assert.Equal(t, worse, style(pair(steady(100), steady(200))))
	assert.Equal(t, faint, style(pair(steady(100), steady(100))), "noise gets no colour")
	assert.Equal(t, bold, style(pair(nil, steady(100))), "an added benchmark is neither")
}

func TestRowsFor(t *testing.T) {
	base := &run{order: []string{"a", "b"}, benches: map[string]*bench{
		"a": {name: "A", samples: map[string][]float64{"sec/op": {1}}},
		"b": {name: "B", samples: map[string][]float64{"sec/op": {2}}},
	}}
	head := &run{order: []string{"a", "c"}, benches: map[string]*bench{
		"a": {name: "A", samples: map[string][]float64{"sec/op": {3}}},
		"c": {name: "C", samples: map[string][]float64{"sec/op": {4}}},
	}}

	rows := rowsFor("sec/op", base, head)
	require.Len(t, rows, 3)

	// Baseline order first, so a chart re-read after a change lists the same
	// benchmarks in the same places; anything the branch added is appended.
	assert.Equal(t, []string{"A", "B", "C"}, []string{rows[0].label, rows[1].label, rows[2].label})
	assert.Equal(t, []float64{1}, rows[0].base)
	assert.Equal(t, []float64{3}, rows[0].head)
	assert.Empty(t, rows[1].head, "B is gone from the current run")
	assert.Empty(t, rows[2].base, "C is new in the current run")
}

func TestRowsForSkipsUnmeasuredMetric(t *testing.T) {
	// A benchmark without ReportAllocs contributes no allocs/op samples; it
	// should drop out of that section rather than take up a row saying
	// nothing.
	base := &run{order: []string{"a"}, benches: map[string]*bench{
		"a": {name: "A", samples: map[string][]float64{"sec/op": {1}}},
	}}
	assert.Empty(t, rowsFor("allocs/op", base, base))
}

func TestEnvDiff(t *testing.T) {
	mac := env{goos: "darwin", goarch: "arm64", cpu: "Apple M4 Max"}

	assert.Empty(t, envDiff(mac, mac))
	assert.Contains(t, envDiff(mac, env{goos: "linux", goarch: "arm64", cpu: "Apple M4 Max"}), "OS (darwin vs. linux)")
	assert.Contains(t, envDiff(mac, env{goos: "darwin", goarch: "arm64", cpu: "Intel"}), "CPU (Apple M4 Max vs. Intel)")
	assert.Empty(t, envDiff(mac, env{goos: "darwin"}), "an unknown field is not a mismatch")
}

// twoRuns builds a baseline and a current run in which exactly one of the
// fixture's two benchmarks moved.
func twoRuns(t *testing.T) (*run, *run) {
	t.Helper()
	base, err := readRun(writeRun(t, fixture))
	require.NoError(t, err)
	head, err := readRun(writeRun(t, strings.ReplaceAll(fixture, "34880 ns/op", "69760 ns/op")))
	require.NoError(t, err)
	return base, head
}

func TestSectionChartsOnlyWhatMoved(t *testing.T) {
	base, head := twoRuns(t)

	out := ansi.Strip(section(metrics["time"], base, head, 100, false))

	// The point of the default: a run covers forty-odd benchmarks, and the
	// one that moved should not have to be found among them.
	assert.Contains(t, out, "picker View/n=10")
	assert.NotContains(t, out, "lister ApplyDedup/n=10")
	assert.Contains(t, out, "1 of 2 moved")
	assert.Contains(t, out, "1 within noise (-all to show)")
	assert.Contains(t, out, "+100.0%")
}

func TestSectionAllRestoresTheRest(t *testing.T) {
	base, head := twoRuns(t)

	out := ansi.Strip(section(metrics["time"], base, head, 100, true))

	assert.Contains(t, out, "picker View/n=10")
	assert.Contains(t, out, "lister ApplyDedup/n=10")
	assert.NotContains(t, out, "-all to show")
	// Moved rows come first, so -all adds context below rather than burying
	// the finding again.
	assert.Less(t, strings.Index(out, "picker View"), strings.Index(out, "lister ApplyDedup"))
}

func TestSectionScalesTheAxisToWhatIsShown(t *testing.T) {
	base, head := twoRuns(t)

	// One benchmark doubled, so the axis spans 100% and that row's bar runs
	// the full half-width.
	out := ansi.Strip(section(metrics["time"], base, head, 100, false))
	assert.Contains(t, out, "axis ±100%")

	// Halve the move and the axis follows it down, so the bar stays long
	// enough to read rather than shrinking into the spine.
	smaller, err := readRun(writeRun(t, strings.ReplaceAll(fixture, "34880 ns/op", "38368 ns/op")))
	require.NoError(t, err)
	out = ansi.Strip(section(metrics["time"], base, smaller, 100, false))
	assert.Contains(t, out, "axis ±10%")
	assert.Contains(t, out, "+10.0%")
}

func TestSectionListsBenchmarksOnlyOneRunHas(t *testing.T) {
	base, err := readRun(writeRun(t, fixture))
	require.NoError(t, err)
	head, err := readRun(writeRun(t, strings.ReplaceAll(fixture, "BenchmarkView", "BenchmarkRenamedView")))
	require.NoError(t, err)

	out := ansi.Strip(section(metrics["time"], base, head, 100, false))
	assert.Contains(t, out, "gone  picker View/n=10")
	assert.Contains(t, out, "new  picker RenamedView/n=10")
	assert.Contains(t, out, "2 added or removed")
}

func TestSectionSaysSoWhenNothingMoved(t *testing.T) {
	base, err := readRun(writeRun(t, fixture))
	require.NoError(t, err)

	out := ansi.Strip(section(metrics["time"], base, base, 100, false))
	assert.Contains(t, out, "0 of 2 moved")
	assert.NotContains(t, out, "│", "no rows means no axis")
}

func TestChart(t *testing.T) {
	base, head := twoRuns(t)

	out := ansi.Strip(chart(base, head, []string{"time", "allocs"}, 80, false))

	assert.Contains(t, out, "main abc1234", "the baseline says what produced it")
	assert.Contains(t, out, "darwin/arm64 · Apple M4 Max")
	assert.Contains(t, out, "time/op")
	assert.Contains(t, out, "allocs/op")
	assert.NotContains(t, out, "⚠", "the two runs came off the same machine")

	// The bar rows are the ones the width budget governs; a baseline path is
	// printed whole rather than truncated.
	for _, line := range strings.Split(out, "\n") {
		if !strings.Contains(line, "│") {
			continue
		}
		assert.LessOrEqual(t, len([]rune(line)), 80, "row exceeds the width given: %q", line)
	}
}

func TestChartWarnsOnADifferentMachine(t *testing.T) {
	base, err := readRun(writeRun(t, fixture))
	require.NoError(t, err)
	head, err := readRun(writeRun(t, strings.ReplaceAll(fixture, "Apple M4 Max", "Apple M1")))
	require.NoError(t, err)

	// A committed baseline is only a target on a comparable machine; the chart
	// says so rather than presenting someone else's laptop as a goal.
	out := ansi.Strip(chart(base, head, []string{"time"}, 80, false))
	assert.Contains(t, out, "CPU (Apple M4 Max vs. Apple M1)")
	assert.Contains(t, out, "just bench-baseline")
}

func TestChartOnlyRendersRequestedMetrics(t *testing.T) {
	base, err := readRun(writeRun(t, fixture))
	require.NoError(t, err)

	out := ansi.Strip(chart(base, base, []string{"allocs"}, 80, false))
	assert.Contains(t, out, "allocs/op")
	assert.NotContains(t, out, "time/op")
	assert.NotContains(t, out, "B/op")
}

func TestValue(t *testing.T) {
	// A benchmark that allocates nothing is the point of some of these; the
	// chart should not print it as "0.000", which reads like a rounded-down
	// measurement.
	assert.Equal(t, "0", value(metrics["allocs"], steady(0)))
	assert.Equal(t, "14.00", value(metrics["allocs"], steady(14)))
	assert.Equal(t, "5.590µs", value(metrics["time"], steady(5.59e-6)))
	assert.Empty(t, value(metrics["time"], nil))
}

func TestTruncate(t *testing.T) {
	assert.Equal(t, "abcdef", truncate("abcdef", 6))
	// The middle goes, not the tail: these names are distinguished at both
	// ends, by the package at the front and the sub-benchmark at the back.
	assert.Equal(t, "ab…ef", truncate("abcdef", 5))
	assert.Equal(t, "picker Fi…10/1-char", truncate("picker FilterSessions/n=10/1-char", 19))
	assert.Equal(t, "abcdef", truncate("abcdef", 0), "a nonsense width is not worth a crash")
}

func TestPad(t *testing.T) {
	// Padded by display width, not by bytes: "69.76µs" is seven columns and
	// eight bytes, and fmt's width verb counts the bytes.
	assert.Equal(t, "  69.76µs", pad("69.76µs", 9, true))
	assert.Equal(t, "69.76µs  ", pad("69.76µs", 9, false))
	assert.Equal(t, "abc", pad("abc", 2, true), "never narrows")
}

func TestSplitMetrics(t *testing.T) {
	assert.Equal(t, []string{"time", "allocs"}, splitMetrics("time, allocs"))
	assert.Equal(t, []string{"time"}, splitMetrics("time,,"))
	assert.Nil(t, splitMetrics(""))
}

func TestCompareRejectsUnknownMetrics(t *testing.T) {
	err := compare("", "", []string{"latency"}, 80, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), `unknown metric "latency"`)

	err = compare("", "", nil, 80, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no metrics selected")
}
