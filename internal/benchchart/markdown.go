package main

import (
	"fmt"
	"math"
	"slices"
	"strings"
)

// The terminal chart and this one answer the same question for different
// readers. A terminal has ANSI colour, two hundred columns and eighth-width
// block glyphs; a pull request has none of those, but it does have a `diff`
// fence — where GitHub renders a line starting with "-" red and one starting
// with "+" green. That maps onto what is being plotted exactly: slower is "+"
// and red, faster is "-" and green.
//
// Keeping the two renderers separate is what lets each use what it has,
// instead of both settling for what they share.

// commentMarker identifies this comment in its own source, for whoever finds
// it there. It is not what locates the comment for editing — the workflow
// uses `gh pr comment --edit-last` for that.
const commentMarker = "<!-- benchchart -->"

// mdWidth is about what a comment shows before it scrolls sideways, and
// mdBarWidth is what the bar keeps of it however long the phrases get: the
// bar is the reason this is a chart rather than a table.
const (
	mdWidth    = 108
	mdBarWidth = 14
)

// markdown renders the comparison as a GitHub comment body: a one-sentence
// verdict, a coloured chart of what moved, and the full table folded away
// underneath.
func markdown(base, head *run, opts options) string {
	var b strings.Builder
	b.WriteString(commentMarker + "\n")
	b.WriteString("## Benchmarks 📊\n\n")

	sections := map[string]sectionData{}
	for _, m := range chosen(opts) {
		moved, shaky, still, missing := classify(m, base, head, opts)
		sections[m.title] = sectionData{m, moved, shaky, still, missing}
	}

	b.WriteString(verdict(chosen(opts), sections))
	if diff := envDiff(base.env, head.env); diff != "" {
		b.WriteString("\n> [!CAUTION]\n> These runs are not comparable: " + diff + ".\n")
	}
	b.WriteString("\n")

	for _, m := range chosen(opts) {
		b.WriteString(mdSection(sections[m.title]))
	}
	b.WriteString(mdFooter(base, head))
	return b.String()
}

// sectionData is one metric's benchmarks, already split and ranked.
type sectionData struct {
	metric                       metric
	moved, shaky, still, missing []line
}

func (s sectionData) counted() counted {
	return counted{len(s.moved), len(s.shaky), len(s.still), len(s.missing)}
}

func (s sectionData) total() int { return s.counted().total() }

// charted is every row worth a line in the fence: what moved, what could not
// be measured well enough to say, and what only one of the two runs has.
func (s sectionData) charted() []line {
	return append(append(slices.Clone(s.moved), s.shaky...), s.missing...)
}

// verdict is the sentence someone reads instead of the chart. GitHub renders
// an alert block with a coloured border and an icon, which is the loudest
// thing available in a comment — so the level has to mean something: it
// escalates only on a regression past the thresholds CONTRIBUTING commits to.
func verdict(ms []metric, sections map[string]sectionData) string {
	level, lede := "NOTE", "No benchmark moved."
	if top, ok := headline(ms, sections); ok {
		lede = fmt.Sprintf("**%s is %s.**", top.label, describe(top))
		switch {
		case regressedPastThreshold(ms, sections):
			level = "CAUTION"
		case top.delta > 0:
			level = "WARNING"
		default:
			level = "TIP"
		}
	}
	return fmt.Sprintf("> [!%s]\n> %s %s\n", level, lede, counts(ms, sections))
}

// headline picks the row the verdict is about: the largest real change, in the
// first metric that has one. classify has already ranked each metric's rows by
// how much work the change costs, so the answer is the top of the list.
func headline(ms []metric, sections map[string]sectionData) (line, bool) {
	for _, m := range ms {
		if s := sections[m.title]; len(s.moved) > 0 {
			return s.moved[0], true
		}
	}
	return line{}, false
}

// regressedPastThreshold asks whether anything moved past the numbers
// CONTRIBUTING says it intends to gate on. Below them a regression is worth
// explaining; above them it is worth stopping for.
func regressedPastThreshold(ms []metric, sections map[string]sectionData) bool {
	for _, m := range ms {
		for _, l := range sections[m.title].moved {
			if l.delta > m.regression {
				return true
			}
		}
	}
	return false
}

// counts is the sentence after the verdict: how much of the suite this is a
// verdict over. It leads with a number whenever anything moved, which is also
// what keeps a metric name off the front of a sentence.
func counts(ms []metric, sections map[string]sectionData) string {
	var parts []string
	moved, shaky, total := 0, 0, 0
	for _, m := range ms {
		s := sections[m.title]
		moved += len(s.moved)
		shaky += len(s.shaky)
		total = max(total, s.total())
		if len(s.moved) == 0 {
			parts = append(parts, m.title+" unchanged")
			continue
		}
		parts = append(parts, fmt.Sprintf("%d of %d moved on %s", len(s.moved), s.total(), m.title))
	}
	sentence := strings.Join(parts, "; ") + "."
	if moved == 0 {
		// The lede already said nothing moved; this says over how much.
		sentence = fmt.Sprintf("Compared %s.", plural(total, "measurement"))
	}
	if shaky > 0 {
		// Said out loud rather than dropped: a run this unstable is a fact
		// about the runner that someone re-reading a flat comment should know.
		sentence += fmt.Sprintf(" %d more could not be measured well enough to say.", shaky)
	}
	return sentence
}

// describe says the delta the way a person would.
func describe(l line) string {
	switch {
	case l.note == "new", l.note == "gone":
		return l.note
	case l.shaky:
		return l.note
	}
	pct := l.delta * 100
	dir := "slower"
	if pct < 0 {
		pct, dir = -pct, "faster"
	}
	return fmt.Sprintf("%.0f%% %s", pct, dir)
}

// mdSection is one metric's chart: a diff fence, so GitHub colours the rows,
// holding only what moved. The rest is in the table below it.
func mdSection(s sectionData) string {
	rows := s.charted()
	if len(rows) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("```diff\n")
	b.WriteString(fmt.Sprintf("@@ %s — %s @@\n", s.metric.title, summary(s.counted(), false, 0, true)))

	peak := 0.0
	for _, l := range s.moved {
		peak = max(peak, math.Abs(l.delta))
	}
	scale := niceScale(peak)

	labelCol, valueCol := 0, 0
	for _, l := range rows {
		labelCol = max(labelCol, len([]rune(l.label)))
		valueCol = max(valueCol, len([]rune(l.from)), len([]rune(l.to)))
	}
	// "- ", the label, the two values and their arrow, and the verdict; the
	// bar takes what is left.
	fixed := 2 + 2 + valueCol + 3 + valueCol + 2 + 2 + 16
	labelCol = min(labelCol, max(mdWidth-fixed-mdBarWidth, 20))
	barCol := max(mdWidth-fixed-labelCol, mdBarWidth)

	for _, l := range rows {
		b.WriteString(fmt.Sprintf("%s %s  %s → %s  %s  %s\n",
			mark(l),
			pad(clip(l.label, labelCol), labelCol, false),
			pad(l.from, valueCol, true),
			pad(l.to, valueCol, true),
			pad(mdBar(l, scale, barCol), barCol, false),
			describe(l)))
	}
	b.WriteString("```\n\n")
	b.WriteString(mdTable(s))
	return b.String()
}

// clip shortens a phrase from the end. A phrase leads with what it measures
// and trails with the conditions, so the front is the part worth keeping —
// the opposite of a benchmark name, where both ends carry meaning.
func clip(s string, n int) string {
	r := []rune(s)
	if len(r) <= n || n < 2 {
		return s
	}
	return strings.TrimRight(string(r[:n-1]), " ") + "…"
}

// mark is the character GitHub's diff highlighter colours the line by.
func mark(l line) string {
	switch {
	case !l.both, l.shaky:
		return "!"
	case l.delta > 0:
		return "+"
	case l.delta < 0:
		return "-"
	}
	return " "
}

// mdBar is one-sided: the "+"/"-" already says which way, so the bar is free
// to spend its whole width on how far.
func mdBar(l line, scale float64, width int) string {
	if scale <= 0 || !l.both || l.shaky {
		return ""
	}
	return strings.Repeat("█", min(int(math.Abs(l.delta)/scale*float64(width)+0.5), width))
}

// mdTable is every benchmark, folded away. The chart above answers "did
// anything move"; this answers "what about the one I care about", which is
// the question the summary line cannot.
func mdTable(s sectionData) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("<details><summary>All %d %s measurements</summary>\n\n", s.total(), s.metric.title))
	b.WriteString("| | What it measures | Before | After | Change |\n|---|---|---|---|---|\n")
	for _, group := range [][]line{s.moved, s.shaky, s.missing, s.still} {
		for _, l := range group {
			b.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				badge(l), l.label, orDash(l.from), orDash(l.to), tableChange(l)))
		}
	}
	b.WriteString("\n</details>\n\n")
	return b.String()
}

// badge is the colour that survives outside a code fence: a table cell has no
// syntax highlighting, but an emoji is an emoji anywhere.
func badge(l line) string {
	switch {
	case !l.both:
		return "🆕"
	case l.shaky:
		return "⚠️"
	case l.note == "~":
		return "⬜"
	case l.delta > 0:
		return "🟥"
	}
	return "🟩"
}

func tableChange(l line) string {
	switch {
	case !l.both, l.shaky:
		return l.note
	case l.note == "~":
		return "no real change"
	}
	return "**" + describe(l) + "**"
}

// mdFooter folds the provenance away. It is the answer to "is this comparison
// even valid", which matters when it is asked and is noise the rest of the
// time.
func mdFooter(base, head *run) string {
	var b strings.Builder
	b.WriteString("<details><summary>How this was measured</summary>\n\n")
	for _, r := range []struct {
		what string
		run  *run
	}{{"Before", base}, {"After", head}} {
		b.WriteString(fmt.Sprintf("**%s** `%s`\n", r.what, r.run.path))
		for _, note := range r.run.notes {
			b.WriteString("- " + note + "\n")
		}
		b.WriteString("\n")
	}
	b.WriteString(head.env.String() + "\n\n</details>\n")
	return b.String()
}

func plural(n int, noun string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, noun)
	}
	return fmt.Sprintf("%d %ss", n, noun)
}

func orDash(s string) string {
	if s == "" {
		return "–"
	}
	return s
}

func capitalise(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
