package lister

import (
	"testing"
)

func TestIsBlacklisted(t *testing.T) {
	tests := []struct {
		name      string
		blacklist []string
		input     string
		expected  bool
	}{
		{
			name:      "matches ~/.config pattern",
			blacklist: []string{`^~/\.config(?:/.*)?$`},
			input:     "~/.config/opencode",
			expected:  true,
		},
		{
			name:      "does not match ~/.config subdirectory",
			blacklist: []string{`^/Users/[^/]+/\.config$`},
			input:     "/Users/username/.config/sesh",
			expected:  false,
		},
		{
			name:      "matches exact string",
			blacklist: []string{"test"},
			input:     "test",
			expected:  true,
		},
		{
			name:      "does not match different string",
			blacklist: []string{"test"},
			input:     "other",
			expected:  false,
		},
		{
			name:      "empty blacklist",
			blacklist: []string{},
			input:     "anything",
			expected:  false,
		},
		{
			name:      "multiple patterns, first matches",
			blacklist: []string{"test", "other"},
			input:     "test",
			expected:  true,
		},
		{
			name:      "multiple patterns, second matches",
			blacklist: []string{"test", "other"},
			input:     "other",
			expected:  true,
		},
		{
			name:      "multiple patterns, none match",
			blacklist: []string{"test", "other"},
			input:     "different",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiled := compileBlacklist(tt.blacklist)
			result := isBlacklisted(compiled, tt.input)
			if result != tt.expected {
				t.Errorf("isBlacklisted(%v, %q) = %v, want %v", tt.blacklist, tt.input, result, tt.expected)
			}
		})
	}
}
