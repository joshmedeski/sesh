package replacer

import (
	"testing"
)

type testCase struct {
	input    string
	expected string
}

func TestReplace(t *testing.T) {
	defaultReplacements := map[string]string{
		"{}": "hello",
		"~":  "/home/test",
	}

	testCases := map[string]testCase{
		"multiple replacements": {
			"~/.local/bin/rat {}{}",
			"/home/test/.local/bin/rat hellohello",
		},
		"single replacement": {
			"/bin/rat {}",
			"/bin/rat hello",
		},
		"no replacement": {
			"/bin/rat",
			"/bin/rat",
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			replacer := NewReplacer()
			result := replacer.Replace(test.input, defaultReplacements)
			if result != test.expected {
				t.Errorf("expected %s, got %s", test.expected, result)
			}
		})
	}
	t.Run("", func(t *testing.T) {
	})
}
