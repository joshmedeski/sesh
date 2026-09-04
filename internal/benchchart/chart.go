package main

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"charm.land/lipgloss/v2"
	"golang.org/x/perf/benchunit"
)

// A metric is one of the columns `-benchmem` reports. All three are
// lower-is-better, which is what lets the chart colour a delta without
// knowing which one it is looking at.
type metric struct {
	unit   string // the benchmark-format unit, as benchfmt tidies it
	title  string // what go test called it, which is what people search for
	suffix string // appended to the SI-scaled value
}

var metrics = map[string]metric{
	"time":   {unit: "sec/op", title: "time/op", suffix: "s"},
	"allocs": {unit: "allocs/op", title: "allocs/op"},
	"bytes":  {unit: "B/op", title: "B/op", suffix: "B"},
}

// metricOrder is the order sections appear in when several are requested,
// independent of the order the flag listed them.
var metricOrder = []string{"time", "allocs", "bytes"}

var (
	faint  = lipgloss.NewStyle().Faint(true)
	bold   = lipgloss.NewStyle().Bold(true)
	better = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(2))
	worse  = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(1))
)

const (
	minAxis  = 13 // narrower than this and the bars stop telling deltas apart
	minLabel = 12
)

// row pairs one benchmark's baseline and current measurements. Either side
// may be absent: a benchmark added on the branch has no baseline, and one
// deleted on the branch has no current run.
type row struct {
	label      string
	base, head []float64
}

// chart renders the whole comparison: a header describing where the two runs
// came from, then one section per metric.
func chart(base, head *run, keys []string, width int, all bool) string {
	var b strings.Builder
	b.WriteString(header(base, head))
	for _, name := range metricOrder {
		m, ok := metrics[name]
		if !ok || !slices.Contains(keys, name) {
			continue
		}
		b.WriteString(section(m, base, head, width, all))
	}
	return b.String()
}

func header(base, head *run) string {
	var b strings.Builder
	b.WriteString(bold.Render("Benchmarks vs. baseline") + "\n\n")
	for _, r := range []struct {
		what string
		run  *run
	}{{"baseline", base}, {"current ", head}} {
		b.WriteString(fmt.Sprintf("  %s  %s\n", bold.Render(r.what), r.run.path))
		b.WriteString("            " + faint.Render(r.run.env.String()) + "\n")
		for _, note := range r.run.notes {
			b.WriteString("            " + faint.Render(note) + "\n")
		}
	}
	// A committed baseline is only a target on a machine like the one that
	// produced it. Saying which parts disagree beats letting someone read a
	// 40% regression off two different laptops.
	if diff := envDiff(base.env, head.env); diff != "" {
		b.WriteString("\n  " + worse.Render("⚠ "+diff) + "\n")
		b.WriteString("  " + faint.Render("These runs are not comparable; regenerate the baseline with `just bench-baseline`.") + "\n")
	}
	b.WriteString("\n  " + faint.Render("bars show the gap from the baseline: faster ◀ │ ▶ slower") + "\n")
	return b.String() + "\n"
}

func envDiff(a, b env) string {
	var diffs []string
	for _, f := range []struct{ what, x, y string }{
		{"OS", a.goos, b.goos},
		{"arch", a.goarch, b.goarch},
		{"CPU", a.cpu, b.cpu},
	} {
		if f.x != "" && f.y != "" && f.x != f.y {
			diffs = append(diffs, fmt.Sprintf("%s (%s vs. %s)", f.what, f.x, f.y))
		}
	}
	if len(diffs) == 0 {
		return ""
	}
	return "baseline and current differ in " + strings.Join(diffs, ", ")
}

// A line is one benchmark reduced to what the chart plots: where it was,
// where it is, and the gap between the two.
type line struct {
	label string
	from  string // the baseline value, SI-scaled
	to    string // the current value, SI-scaled
	note  string // the delta, or "new"/"gone"
	delta float64
	style lipgloss.Style
	both  bool // false for a benchmark only one of the two runs has
}

// section renders one metric. It prints the benchmarks that moved and
// collapses the rest to a count: a run covers forty-odd benchmarks, and a
// wall of unchanged rows buries the two that matter. -all opts back in.
func section(m metric, base, head *run, width int, all bool) string {
	var moved, still, missing []line
	for _, r := range rowsFor(m.unit, base, head) {
		l := lineFor(m, r)
		switch {
		case !l.both:
			missing = append(missing, l)
		case significant(r):
			moved = append(moved, l)
		default:
			still = append(still, l)
		}
	}
	if len(moved)+len(still)+len(missing) == 0 {
		return ""
	}

	shown := moved
	if all {
		shown = append(slices.Clone(moved), still...)
	}
	// The axis spans the largest delta on show, so the bars stay comparable
	// to each other within a section and never all run to the edge.
	var peak float64
	for _, l := range shown {
		peak = max(peak, math.Abs(l.delta))
	}
	scale := niceScale(peak)

	var b strings.Builder
	b.WriteString(bold.Render(m.title) + "  " + faint.Render(summary(len(moved), len(still), len(missing), len(shown) > 0, scale, all)) + "\n\n")
	b.WriteString(plot(shown, scale, width))
	for _, l := range missing {
		b.WriteString(fmt.Sprintf("  %s  %s  %s\n", l.style.Render(l.note), l.label, faint.Render(l.from+l.to)))
	}
	return b.String() + "\n"
}

// summary is the one line that answers "is there anything here" without
// reading a single bar.
func summary(moved, still, missing int, plotted bool, scale float64, all bool) string {
	total := moved + still + missing
	parts := []string{fmt.Sprintf("%d of %d moved", moved, total)}
	if still > 0 && !all {
		parts = append(parts, fmt.Sprintf("%d within noise (-all to show)", still))
	}
	if missing > 0 {
		parts = append(parts, fmt.Sprintf("%d added or removed", missing))
	}
	if plotted {
		parts = append(parts, "axis ±"+percent(scale))
	}
	return strings.Join(parts, " · ")
}

func lineFor(m metric, r row) line {
	return line{
		label: r.label,
		from:  value(m, r.base),
		to:    value(m, r.head),
		note:  delta(r),
		delta: relative(r),
		style: style(r),
		both:  len(r.base) > 0 && len(r.head) > 0,
	}
}

func value(m metric, xs []float64) string {
	if len(xs) == 0 {
		return ""
	}
	v := median(xs)
	if v == 0 {
		// benchunit renders an exact zero as "0.000", which reads like a
		// rounded-down measurement rather than the allocation-free path it is.
		return "0"
	}
	return benchunit.Scale(v, benchunit.ClassOf(m.unit)) + m.suffix
}

// plot lays the lines out in columns and draws each delta against the shared
// axis. Measuring every bar against one span is what makes the chart
// readable: a 40% regression is five times the bar of an 8% one, rather than
// both running to the edge of the row as they would if each were scaled to
// itself.
func plot(lines []line, scale float64, width int) string {
	if len(lines) == 0 {
		return ""
	}
	labelCol, fromCol, toCol := 0, 0, 0
	for _, l := range lines {
		labelCol = max(labelCol, lipgloss.Width(l.label))
		fromCol = max(fromCol, lipgloss.Width(l.from))
		toCol = max(toCol, lipgloss.Width(l.to))
	}

	// Names and axis share what the fixed columns leave. The axis is reserved
	// its share first: a benchmark name survives being truncated, and the bar
	// that is the reason to look at this at all does not survive being
	// squeezed into three cells by one long name.
	const deltaCol = 7 // "+123.4%"
	fixed := 2 + 2 + fromCol + 3 + toCol + 2 + 2 + deltaCol
	labelCol = min(labelCol, max(width-fixed-axisShare(width-fixed), minLabel))
	half := max((width-fixed-labelCol-1)/2, minAxis/2)

	var b strings.Builder
	for _, l := range lines {
		b.WriteString(fmt.Sprintf("  %s  %s → %s  %s  %s\n",
			pad(truncate(l.label, labelCol), labelCol, false),
			pad(l.from, fromCol, true),
			pad(l.to, toCol, true),
			axis(l, scale, half),
			l.style.Render(l.note)))
	}
	return b.String()
}

// axis draws one delta as a bar growing out of a centre line: left and green
// for faster, right and red for slower, both measured against scale. The
// centre lines stack into a spine, so which side a row falls on reads at a
// glance.
func axis(l line, scale float64, half int) string {
	left, right := strings.Repeat(" ", half), strings.Repeat(" ", half)
	cells := 0
	if scale > 0 && l.both {
		cells = min(int(math.Round(math.Abs(l.delta)/scale*float64(half))), half)
	}
	switch {
	case cells == 0:
	case l.delta < 0:
		left = strings.Repeat(" ", half-cells) + l.style.Render(strings.Repeat("█", cells))
	default:
		right = l.style.Render(strings.Repeat("█", cells)) + strings.Repeat(" ", half-cells)
	}
	return left + faint.Render("│") + right
}

// niceScale rounds the largest delta in a section up to a round number, so
// the axis reads in steps someone can hold in their head. The floor keeps a
// section where nothing moved more than 1% from magnifying that into a full
// bar.
func niceScale(peak float64) float64 {
	for _, s := range []float64{0.05, 0.1, 0.25, 0.5, 1, 2, 5, 10} {
		if peak <= s {
			return s
		}
	}
	return peak
}

func percent(f float64) string {
	return strings.TrimSuffix(strings.TrimRight(fmt.Sprintf("%.1f", f*100), "0"), ".") + "%"
}

// relative is the delta as a fraction of the baseline: what the bars are
// drawn from.
func relative(r row) float64 {
	if len(r.base) == 0 || len(r.head) == 0 {
		return 0
	}
	b, h := median(r.base), median(r.head)
	if b == 0 {
		return 0
	}
	return (h - b) / b
}

// delta is the headline number. A move smaller than the two runs' own sample
// spread is reported as "~", because at -count=6 that is all the data can
// support.
func delta(r row) string {
	switch {
	case len(r.base) == 0:
		return "new"
	case len(r.head) == 0:
		return "gone"
	case !significant(r):
		return "~"
	}
	return fmt.Sprintf("%+.1f%%", relative(r)*100)
}

// significant asks whether a delta clears the noise of the runs it came from.
// It is what decides whether a benchmark is worth a row at all, so the floor
// matters: without it, a pair of unusually steady runs reports a 0.3%
// difference as a finding.
func significant(r row) bool {
	if len(r.base) == 0 || len(r.head) == 0 || median(r.base) == 0 {
		return false
	}
	return math.Abs(relative(r)) > max(spread(r.base)+spread(r.head), 0.02)
}

func style(r row) lipgloss.Style {
	switch {
	case len(r.base) == 0 || len(r.head) == 0:
		return bold
	case !significant(r):
		return faint
	case median(r.head) < median(r.base):
		return better
	}
	return worse
}

// rowsFor pairs the two runs up, baseline order first so a chart re-read
// after a change lists the same benchmarks in the same places, with anything
// the branch added appended.
func rowsFor(unit string, base, head *run) []row {
	var rows []row
	seen := map[string]bool{}
	for _, key := range append(slices.Clone(base.order), head.order...) {
		if seen[key] {
			continue
		}
		seen[key] = true
		b, h := base.benches[key], head.benches[key]
		r := row{}
		if b != nil {
			r.label, r.base = b.label(), b.samples[unit]
		}
		if h != nil {
			r.label, r.head = h.label(), h.samples[unit]
		}
		if len(r.base) == 0 && len(r.head) == 0 {
			continue
		}
		rows = append(rows, r)
	}
	return rows
}

// pad widens s to n columns, measured by display width: "69.76µs" is seven
// columns and eight bytes, and fmt's width verb counts the bytes.
func pad(s string, n int, right bool) string {
	fill := strings.Repeat(" ", max(n-lipgloss.Width(s), 0))
	if right {
		return fill + s
	}
	return s + fill
}

// axisShare is how much of the free space the axis claims.
func axisShare(free int) int {
	return max(free*2/5, minAxis)
}

// truncate elides the middle rather than the tail. These names carry what
// distinguishes them at both ends — the package at the front, the
// sub-benchmark configuration at the back — and cutting the tail turns a
// screen of them into a column of "lister CachingListerApplyF…".
func truncate(s string, n int) string {
	r := []rune(s)
	if lipgloss.Width(s) <= n || n < 3 {
		return s
	}
	head := (n - 1) / 2
	return string(r[:head]) + "…" + string(r[len(r)-(n-1-head):])
}
