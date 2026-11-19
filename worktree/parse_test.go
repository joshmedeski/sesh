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
			repo, number, isPR, ok := ParseGitHubRef(tt.url)
			if ok != tt.ok || repo != tt.repo || number != tt.number || isPR != tt.isPR {
				t.Fatalf("ParseGitHubRef(%q) = (%q, %d, %v, %v); want (%q, %d, %v, %v)",
					tt.url, repo, number, isPR, ok, tt.repo, tt.number, tt.isPR, tt.ok)
			}
		})
	}
}

func TestParseForgeRef(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		repo   string
		number int
		isPR   bool
		ok     bool
	}{
		{"github issue", "https://github.com/joshmedeski/sesh/issues/409", "joshmedeski/sesh", 409, false, true},
		{"github pull", "https://github.com/joshmedeski/sesh/pull/410/files", "joshmedeski/sesh", 410, true, true},
		{"github no scheme", "github.com/joshmedeski/sesh/issues/409", "joshmedeski/sesh", 409, false, true},
		{"github enterprise", "https://github.example.com/enterprise/repo/pull/55", "enterprise/repo", 55, true, true},
		{"gitlab issue", "https://gitlab.com/inkscape/inkscape/-/issues/99", "inkscape/inkscape", 99, false, true},
		{"gitlab merge request", "https://gitlab.com/inkscape/inkscape/-/merge_requests/101/diffs", "inkscape/inkscape", 101, true, true},
		{"gitlab subgroup", "https://gitlab.example.com/group/sub/project/-/merge_requests/7#note_1", "group/sub/project", 7, true, true},
		{"gitlab legacy path", "https://gitlab.com/inkscape/inkscape/issues/99", "inkscape/inkscape", 99, false, true},
		{"repo home", "https://gitlab.com/inkscape/inkscape", "", 0, false, false},
		{"garbage", "not a url", "", 0, false, false},
		{"empty", "", "", 0, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, number, isPR, ok := parseForgeRef(tt.url)
			if ok != tt.ok || repo != tt.repo || number != tt.number || isPR != tt.isPR {
				t.Fatalf("parseForgeRef(%q) = (%q, %d, %v, %v); want (%q, %d, %v, %v)",
					tt.url, repo, number, isPR, ok, tt.repo, tt.number, tt.isPR, tt.ok)
			}
		})
	}
}
