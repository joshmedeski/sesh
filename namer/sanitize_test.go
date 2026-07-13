package namer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeTitle(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"plain title kept", "warm the status cache", "warm the status cache"},
		{"casing preserved", "Warm The Cache", "Warm The Cache"},
		{"colon replaced with space", "fix: crash on start", "fix crash on start"},
		{"dot replaced with space", "bump v2.0 release", "bump v2 0 release"},
		{"collapses resulting double spaces", "a:  b", "a b"},
		{"trims leading and trailing space", "  hello  ", "hello"},
		{"colon then space stays single space", "feat: thing", "feat thing"},
		{"empty stays empty", "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, SanitizeTitle(c.input))
		})
	}
}
