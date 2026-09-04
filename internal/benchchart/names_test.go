package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPhrase(t *testing.T) {
	tests := []struct {
		name, in, want string
	}{
		{
			"folds the size in and drops a default condition",
			"FilterSessions/n=1000/plain/empty",
			"Filtering 1,000 sessions on a keystroke (nothing typed)",
		},
		{
			"keeps the conditions that distinguish siblings",
			"FilterSessions/n=10/separator-aware/no-match",
			"Filtering 10 sessions on a keystroke (separator-aware, no matches)",
		},
		{
			"a benchmark with no conditions is just the sentence",
			"View/n=1000",
			"Redrawing the picker, 1,000 sessions",
		},
		{
			"the same token can mean different things in different benchmarks",
			"BuildItems/n=100/plain",
			"Building picker rows for 100 sessions (plain rows)",
		},
		{
			"an unnamed benchmark still reads as words",
			"SomeNewThing/n=10",
			"Some new thing, 10 sessions",
		},
		{
			"an unknown condition is passed through rather than dropped",
			"View/n=10/wibble",
			"Redrawing the picker, 10 sessions (wibble)",
		},
		{
			"a benchmark with no size still gets a subject",
			"View",
			"Redrawing the picker, the session list",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, phrase(tt.in))
		})
	}
}

// TestEveryBenchmarkIsNamed is what keeps the phrases from rotting. The
// fallback means an unnamed benchmark renders as de-camel-cased soup rather
// than breaking, which is exactly the kind of failure nobody notices — so
// adding a benchmark without a phrase fails here instead.
func TestEveryBenchmarkIsNamed(t *testing.T) {
	baselines, err := filepath.Glob(filepath.Join("..", "..", "testdata", "bench", "*.txt"))
	require.NoError(t, err)
	require.NotEmpty(t, baselines, "no committed baseline to check the phrases against")

	line := regexp.MustCompile(`(?m)^Benchmark(\S+?)(-\d+)?\s`)
	for _, path := range baselines {
		body, err := os.ReadFile(path)
		require.NoError(t, err)

		for _, m := range line.FindAllStringSubmatch(string(body), -1) {
			base, parts := splitName(m[1])
			assert.Contains(t, what, base,
				"%s: benchmark %q has no phrase in names.go", filepath.Base(path), base)

			_, conds := countOf(parts)
			for _, c := range conds {
				_, qualified := conditions[base+"/"+c]
				_, plain := conditions[c]
				assert.True(t, qualified || plain,
					"%s: %s has no phrase for the condition %q", filepath.Base(path), base, c)
			}
		}
	}
}

func TestCommas(t *testing.T) {
	tests := []struct{ in, want string }{
		{"10", "10"},
		{"100", "100"},
		{"1000", "1,000"},
		{"1000000", "1,000,000"},
		{"", ""},
		{"abc", "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			assert.Equal(t, tt.want, commas(tt.in))
		})
	}
}

func TestDeCamel(t *testing.T) {
	assert.Equal(t, "Find tmux session by base", deCamel("FindTmuxSessionByBase"))
	assert.Equal(t, "View", deCamel("View"))
	assert.Equal(t, "", deCamel(""))
}

func TestCountOf(t *testing.T) {
	sessions, rest := countOf([]string{"n=1000", "plain", "empty"})
	assert.Equal(t, "1,000 sessions", sessions)
	assert.Equal(t, []string{"plain", "empty"}, rest)

	sessions, rest = countOf(nil)
	assert.Equal(t, "the session list", sessions, "a benchmark with no size still needs a subject")
	assert.Empty(t, rest)
}

func TestPhrasesAreShortEnoughToRead(t *testing.T) {
	// The markdown chart clips a phrase that will not fit, and a chart of
	// clipped phrases is a chart of "Filtering 1,000 sessions on a keystro…".
	// Nothing in the committed vocabulary should be near that.
	for base, head := range what {
		assert.LessOrEqual(t, len(strings.ReplaceAll(head, "{n}", "1,000 sessions")), 48,
			"the phrase for %s leaves no room for its conditions", base)
	}
	for token, cond := range conditions {
		assert.LessOrEqual(t, len(cond), 28, "the condition phrase for %s is too long", token)
	}
}
