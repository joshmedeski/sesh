package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/perf/benchfmt"
)

// A run is one `go test -bench` output file: the measurements it holds, plus
// enough of the environment that produced them to tell whether two runs are
// comparable at all.
type run struct {
	path    string
	env     env
	notes   []string
	order   []string
	benches map[string]*bench
}

// env is the part of a run that has to match for a comparison to mean
// anything. go test writes all four as benchmark-format configuration lines,
// so they come along with the measurements for free.
type env struct {
	goos, goarch, cpu string
}

func (e env) String() string {
	parts := []string{}
	for _, p := range []string{e.goos + "/" + e.goarch, e.cpu} {
		if p != "" && p != "/" {
			parts = append(parts, p)
		}
	}
	return strings.Join(parts, " · ")
}

// bench is every sample of every metric for one benchmark, across the -count
// repetitions in the file.
type bench struct {
	pkg     string
	name    string
	samples map[string][]float64
}

// label is what the chart prints: the package's last element, which is all
// that distinguishes lister from picker here, and the benchmark name with its
// sub-benchmark configuration.
func (b *bench) label() string {
	if b.pkg == "" {
		return b.name
	}
	return path.Base(b.pkg) + " " + b.name
}

// readRun parses a benchmark-format file. Lines go test emits that are not
// benchmark results — "PASS", "ok  \t...", and the "#" provenance comments
// `just bench-baseline` writes — are not part of the format, and benchfmt
// ignores them; the comments are re-read separately by readNotes.
func readRun(p string) (*run, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := &run{path: p, benches: map[string]*bench{}}
	var buf bytes.Buffer
	rd := benchfmt.NewReader(io.TeeReader(f, &buf), p)
	for rd.Scan() {
		res, ok := rd.Result().(*benchfmt.Result)
		if !ok {
			// A SyntaxError: report it and keep reading, since one malformed
			// line should not throw away the rest of the file.
			fmt.Fprintf(os.Stderr, "%s: %v\n", p, rd.Result())
			continue
		}
		r.record(res)
	}
	if err := rd.Err(); err != nil {
		return nil, err
	}
	if len(r.benches) == 0 {
		return nil, fmt.Errorf("%s: no benchmark results", p)
	}
	r.notes = readNotes(&buf)
	return r, nil
}

func (r *run) record(res *benchfmt.Result) {
	r.env.goos = orKeep(r.env.goos, res.GetConfig("goos"))
	r.env.goarch = orKeep(r.env.goarch, res.GetConfig("goarch"))
	r.env.cpu = orKeep(r.env.cpu, res.GetConfig("cpu"))

	pkg := res.GetConfig("pkg")
	name := trimProcs(res.Name.String())
	key := pkg + "\x00" + name
	b, ok := r.benches[key]
	if !ok {
		b = &bench{pkg: pkg, name: name, samples: map[string][]float64{}}
		r.benches[key] = b
		r.order = append(r.order, key)
	}
	for _, v := range res.Values {
		b.samples[v.Unit] = append(b.samples[v.Unit], v.Value)
	}
}

// trimProcs drops the "-16" GOMAXPROCS suffix go test appends to every
// benchmark name. It records the core count of the machine that ran it, so
// leaving it on would make a 16-core baseline share no benchmark names at all
// with a 10-core worktree run.
func trimProcs(name string) string {
	i := strings.LastIndexByte(name, '-')
	if i <= 0 {
		return name
	}
	if _, err := strconv.Atoi(name[i+1:]); err != nil {
		return name
	}
	return name[:i]
}

// orKeep prefers the first value seen. A file that mixes machines is already
// broken; keeping the first keeps the mismatch visible in the header rather
// than letting the last result silently define the environment.
func orKeep(have, next string) string {
	if have != "" || next == "" {
		return have
	}
	return next
}

// readNotes collects the "#" comment lines a generated baseline carries. They
// are not benchmark-format configuration on purpose: a configuration key that
// appears in one file and not the other changes how benchstat groups the two,
// and the issue this tool answers wants the baseline file to stay something
// benchstat can still consume directly.
func readNotes(r io.Reader) []string {
	var notes []string
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if note, ok := strings.CutPrefix(line, "#"); ok {
			notes = append(notes, strings.TrimSpace(note))
		}
	}
	return notes
}

// median is the middle sample, averaging the two middle ones for an even
// count. The chart wants the typical run, not the mean, which one slow
// sample drags around.
func median(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	s := append([]float64(nil), xs...)
	sort.Float64s(s)
	mid := len(s) / 2
	if len(s)%2 == 1 {
		return s[mid]
	}
	return (s[mid-1] + s[mid]) / 2
}

// spread is half the sample range as a fraction of the median: a cheap,
// assumption-free stand-in for the confidence interval benchstat computes.
// It is used only to decide whether a delta is worth colouring, so being
// conservative — a wider spread than a real interval — is the safe direction.
func spread(xs []float64) float64 {
	m := median(xs)
	if len(xs) < 2 || m == 0 {
		return 0
	}
	lo, hi := xs[0], xs[0]
	for _, x := range xs {
		lo, hi = min(lo, x), max(hi, x)
	}
	return (hi - lo) / 2 / m
}
