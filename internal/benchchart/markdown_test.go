package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mdOpts() options {
	return options{metrics: []string{"time", "allocs"}}
}

// at builds one already-classified line, which is all the verdict looks at.
func at(label string, delta float64) line {
	return line{label: label, delta: delta, both: true, note: "x"}
}

func only(m metric, moved ...line) map[string]sectionData {
	return map[string]sectionData{m.title: {metric: m, moved: moved}}
}

func TestVerdictLevel(t *testing.T) {
	time := metrics["time"]
	ms := []metric{time}

	tests := []struct {
		name  string
		moved []line
		want  string
	}{
		{"nothing moved is not news", nil, "> [!NOTE]"},
		{"only improvements", []line{at("A", -0.19)}, "> [!TIP]"},
		{"a regression inside the threshold is worth explaining", []line{at("A", 0.08)}, "> [!WARNING]"},
		{"a regression past CONTRIBUTING's threshold is worth stopping for", []line{at("A", 0.4)}, "> [!CAUTION]"},
		{"a big regression outweighs a big win", []line{at("A", -0.5), at("B", 0.4)}, "> [!CAUTION]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Contains(t, verdict(ms, only(time, tt.moved...)), tt.want)
		})
	}
}

func TestVerdictLede(t *testing.T) {
	time := metrics["time"]

	got := verdict([]metric{time}, only(time, at("Filtering 1,000 sessions", -0.194), at("Redrawing", -0.116)))
	// The top row is the one classify ranked first by how much work the change
	// costs, and it is the sentence someone reads instead of the chart.
	assert.Contains(t, got, "**Filtering 1,000 sessions is 19% faster.**")
	assert.NotContains(t, got, "Redrawing")

	assert.Contains(t, verdict([]metric{time}, only(time)), "No benchmark moved.")
}

func TestCounts(t *testing.T) {
	time, allocs := metrics["time"], metrics["allocs"]
	sections := map[string]sectionData{
		time.title:   {metric: time, moved: []line{at("A", 0.1)}, still: make([]line, 95)},
		allocs.title: {metric: allocs, still: make([]line, 96)},
	}
	got := counts([]metric{time, allocs}, sections)
	assert.Equal(t, "1 of 96 moved on time/op; allocs/op unchanged.", got)
}

func TestDescribe(t *testing.T) {
	assert.Equal(t, "19% faster", describe(line{delta: -0.194, note: "-19.4%"}))
	assert.Equal(t, "3% slower", describe(line{delta: 0.032, note: "+3.2%"}))
	assert.Equal(t, "new", describe(line{note: "new"}))
	assert.Equal(t, "gone", describe(line{note: "gone"}))
}

func TestMark(t *testing.T) {
	// GitHub colours a diff fence by the line's first character, and the
	// mapping happens to be exactly right: slower is "+" and red, faster is
	// "-" and green.
	assert.Equal(t, "+", mark(line{delta: 0.1, both: true}))
	assert.Equal(t, "-", mark(line{delta: -0.1, both: true}))
	assert.Equal(t, " ", mark(line{delta: 0, both: true}))
	assert.Equal(t, "!", mark(line{delta: 0.1}), "a benchmark only one run has is neither")
}

func TestBadge(t *testing.T) {
	assert.Equal(t, "🟥", badge(line{delta: 0.1, both: true, note: "+10.0%"}))
	assert.Equal(t, "🟩", badge(line{delta: -0.1, both: true, note: "-10.0%"}))
	assert.Equal(t, "⬜", badge(line{delta: 0.001, both: true, note: "~"}))
	assert.Equal(t, "🆕", badge(line{note: "new"}))
}

func TestMdBar(t *testing.T) {
	assert.Equal(t, 10, len([]rune(mdBar(line{delta: 0.5, both: true}, 0.5, 10))), "the span fills the bar")
	assert.Equal(t, 5, len([]rune(mdBar(line{delta: -0.25, both: true}, 0.5, 10))), "direction is the mark's job, not the bar's")
	assert.Equal(t, 10, len([]rune(mdBar(line{delta: 5, both: true}, 0.5, 10))), "past the span is clamped")
	assert.Empty(t, mdBar(line{delta: 0.5}, 0.5, 10), "a benchmark only one run has has nothing to draw")
}

func TestClip(t *testing.T) {
	// A phrase leads with what it measures, so the front is the part to keep —
	// the opposite of truncate(), which elides the middle of a benchmark name.
	assert.Equal(t, "Filtering 1,000…", clip("Filtering 1,000 sessions on a keystroke", 16))
	assert.Equal(t, "Redrawing", clip("Redrawing", 20))
}

func TestMarkdown(t *testing.T) {
	base, head := twoRuns(t)
	got := markdown(base, head, mdOpts())

	assert.True(t, strings.HasPrefix(got, commentMarker), "the marker identifies the comment in its own source")
	assert.Contains(t, got, "> [!CAUTION]", "the fixture's benchmark doubled")
	assert.Contains(t, got, "```diff\n", "the fence is what makes GitHub colour the rows")
	assert.Contains(t, got, "@@ time/op — 1 of 2 moved @@")
	assert.Contains(t, got, "+ Redrawing the picker, 10 sessions", "a regression leads with the character GitHub renders red")
	assert.Contains(t, got, "100% slower")
	assert.Contains(t, got, "<details><summary>All 2 time/op measurements</summary>")
	assert.Contains(t, got, "| 🟥 | Redrawing the picker, 10 sessions |")
	assert.Contains(t, got, "<details><summary>How this was measured</summary>")
	assert.Contains(t, got, "main abc1234")
}

func TestMarkdownRowsFitTheCommentWidth(t *testing.T) {
	base, err := readRun(writeRun(t, fixture))
	require.NoError(t, err)
	// A phrase long enough to need clipping, so the budget is actually tested.
	head, err := readRun(writeRun(t, strings.ReplaceAll(fixture, "34880 ns/op", "69760 ns/op")))
	require.NoError(t, err)

	for _, line := range strings.Split(markdown(base, head, mdOpts()), "\n") {
		if !strings.HasPrefix(line, "+ ") && !strings.HasPrefix(line, "- ") {
			continue
		}
		assert.LessOrEqual(t, len([]rune(line)), mdWidth, "row wider than a comment shows: %q", line)
	}
}

func TestMarkdownSkipsAMetricThatDidNotMove(t *testing.T) {
	base, err := readRun(writeRun(t, fixture))
	require.NoError(t, err)

	got := markdown(base, base, mdOpts())
	assert.Contains(t, got, "> [!NOTE]")
	assert.NotContains(t, got, "```diff", "a section with nothing in it is not worth a fence")
	assert.Contains(t, got, "Nothing moved across 2 measurements.")
}

func TestMarkdownWarnsOnADifferentMachine(t *testing.T) {
	base, err := readRun(writeRun(t, fixture))
	require.NoError(t, err)
	head, err := readRun(writeRun(t, strings.ReplaceAll(fixture, "Apple M4 Max", "Apple M1")))
	require.NoError(t, err)

	got := markdown(base, head, mdOpts())
	assert.Contains(t, got, "> [!CAUTION]\n> These runs are not comparable")
	assert.Contains(t, got, "CPU (Apple M4 Max vs. Apple M1)")
}

func TestCompareRejectsAnUnknownFormat(t *testing.T) {
	err := compare("", "", "html", options{metrics: []string{"time"}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), `unknown format "html"`)
}
