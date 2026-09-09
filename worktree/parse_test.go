package worktree

import "testing"

func TestParseGitHubRef(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		repo   string
		number int
		isPR   bool
		ok     bool
	}{
		{"issue", "https://github.com/joshmedeski/sesh/issues/409", "joshmedeski/sesh", 409, false, true},
		{"pull", "https://github.com/joshmedeski/sesh/pull/410", "joshmedeski/sesh", 410, true, true},
		{"trailing slash", "https://github.com/joshmedeski/sesh/issues/409/", "joshmedeski/sesh", 409, false, true},
		{"query string", "https://github.com/joshmedeski/sesh/issues/409?foo=bar", "joshmedeski/sesh", 409, false, true},
		{"fragment", "https://github.com/joshmedeski/sesh/pull/410#issuecomment-1", "joshmedeski/sesh", 410, true, true},
		{"no scheme", "github.com/joshmedeski/sesh/issues/409", "joshmedeski/sesh", 409, false, true},
		{"repo home", "https://github.com/joshmedeski/sesh", "", 0, false, false},
		{"non-github", "https://gitlab.com/joshmedeski/sesh/issues/409", "", 0, false, false},
		{"lookalike host", "https://mygithub.com/joshmedeski/sesh/issues/409", "", 0, false, false},
		{"garbage", "not a url", "", 0, false, false},
		{"empty", "", "", 0, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, number, isPR, ok := parseGitHubRef(tt.url)
			if ok != tt.ok || repo != tt.repo || number != tt.number || isPR != tt.isPR {
				t.Fatalf("parseGitHubRef(%q) = (%q, %d, %v, %v); want (%q, %d, %v, %v)",
					tt.url, repo, number, isPR, ok, tt.repo, tt.number, tt.isPR, tt.ok)
			}
		})
	}
}
