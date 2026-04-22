package startup

import "testing"

func TestPosixSingleQuote(t *testing.T) {
	cases := []struct {
		name, in, want string
	}{
		{"empty", "", "''"},
		{"plain", "hello", "'hello'"},
		{"spaces", "echo hi there", "'echo hi there'"},
		{"dollar-stays-literal", "$USER", "'$USER'"},
		{"backslash-stays-literal", `a\b`, `'a\b'`},
		{"single-quote", "it's", `'it'\''s'`},
		{"multiple-quotes", "a'b'c", `'a'\''b'\''c'`},
		{"newline", "a\nb", "'a\nb'"},
		{"double-quote-stays-literal", `say "hi"`, `'say "hi"'`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := posixSingleQuote(tc.in)
			if got != tc.want {
				t.Errorf("posixSingleQuote(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
