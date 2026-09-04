// Command benchchart charts a `go test -bench` run against the baseline
// committed under testdata/bench, so a regression is visible as a shorter or
// longer bar rather than as a column of nanosecond counts.
//
// It is a development tool, not part of the sesh binary. Reach for it through
// the justfile:
//
//	just bench-baseline   # on main: record the baseline for this machine
//	just bench-compare    # in a worktree: run the suite and chart it
//
// It is informational, like the benchmark job in CI: it never exits non-zero
// on a regression, only on a failure to read one of the two runs.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/colorprofile"
	"golang.org/x/term"
)

func main() {
	baseline := flag.String("baseline", "", "baseline file (default testdata/bench/baseline-GOOS-GOARCH.txt)")
	current := flag.String("current", ".bench/current.txt", "the run to compare against the baseline")
	list := flag.String("metric", "time,allocs", "metrics to chart: time, allocs, bytes")
	width := flag.Int("width", 0, "chart width in columns (default: the terminal's)")
	all := flag.Bool("all", false, "chart every benchmark, not only the ones that moved")
	raw := flag.Bool("raw-names", false, "label rows by benchmark name rather than in English")
	format := flag.String("format", "text", "text for a terminal, markdown for a pull-request comment")
	flag.Parse()

	if *baseline == "" {
		*baseline = baselinePath(runtime.GOOS, runtime.GOARCH)
	}
	opts := options{metrics: splitMetrics(*list), width: chartWidth(*width), all: *all, raw: *raw}
	if err := compare(*baseline, *current, *format, opts); err != nil {
		fmt.Fprintln(os.Stderr, "benchchart:", err)
		os.Exit(1)
	}
}

// baselinePath keys the baseline by platform. One file holding both a Mac and
// a Linux runner's numbers would compare them against each other, which is
// worse than having no baseline at all.
func baselinePath(goos, goarch string) string {
	return filepath.Join("testdata", "bench", fmt.Sprintf("baseline-%s-%s.txt", goos, goarch))
}

func compare(baseline, current, format string, opts options) error {
	render, ok := renderers[format]
	if !ok {
		return fmt.Errorf("unknown format %q; -format takes text or markdown", format)
	}
	if len(opts.metrics) == 0 {
		return fmt.Errorf("no metrics selected; -metric takes some of time, allocs, bytes")
	}
	for _, k := range opts.metrics {
		if _, ok := metrics[k]; !ok {
			return fmt.Errorf("unknown metric %q; -metric takes some of time, allocs, bytes", k)
		}
	}

	base, err := readRun(baseline)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no baseline at %s; record one on main with `just bench-baseline`", baseline)
		}
		return err
	}
	head, err := readRun(current)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no run at %s; produce one with `just bench-compare`", current)
		}
		return err
	}

	// Downsamples the chart's colour to whatever stdout can actually take,
	// which for a redirect to a file or a pipe into a pager means stripping it
	// rather than writing escape codes into the output. It honours NO_COLOR
	// and CLICOLOR_FORCE on the way. The markdown renderer emits no ANSI at
	// all, so this is a no-op for it.
	out := colorprofile.NewWriter(os.Stdout, os.Environ())
	fmt.Fprint(out, render(base, head, opts))
	return nil
}

// renderers are the two audiences: a terminal, and a pull request.
var renderers = map[string]func(base, head *run, opts options) string{
	"text":     chart,
	"markdown": markdown,
}

// chartWidth prefers an explicit -width, then the terminal's, then 80. The
// output is routinely piped to a pager, where there is no terminal to ask.
func chartWidth(want int) int {
	if want > 0 {
		return want
	}
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		return w
	}
	return 80
}

func splitMetrics(list string) []string {
	var keys []string
	for _, k := range strings.Split(list, ",") {
		if k = strings.TrimSpace(k); k != "" {
			keys = append(keys, k)
		}
	}
	return keys
}
