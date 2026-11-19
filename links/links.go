package links

import (
	"regexp"
	"slices"
	"strings"
)

var (
	urlRe    = regexp.MustCompile(`(?:https?|ftp|file):/?//[-A-Za-z0-9+&@#/%?=~_|!:,.;]*[-A-Za-z0-9+&@#/%=~_|]`)
	wwwRe    = regexp.MustCompile(`www\.[a-zA-Z](?:-?[a-zA-Z0-9])+\.[a-zA-Z]{2,}(?:/\S+)*`)
	ipRe     = regexp.MustCompile(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}(?::[0-9]{1,5})?(?:/\S+)*`)
	gitRe    = regexp.MustCompile(`(?:ssh://)?git@(\S+)`)
	ghRe     = regexp.MustCompile(`['"]([A-Za-z0-9-]+/[.A-Za-z0-9-]+)['"]`)
	npmRe    = regexp.MustCompile(`import\s+[^"';]*["']([^.][^"';]*)["']`)
	issueRe  = regexp.MustCompile(`(?:^|[^\w/&#])#(\d+)\b`)
	githubRe = regexp.MustCompile(`github\.com[:/]([^/\s]+/[^/\s]+?)(?:\.git)?/?$`)
)

type match struct {
	pos int
	url string
}

// Extract returns the unique links in text, in the order they appear. repo is
// the GitHub "org/repo" a bare #1234 refers to; with none, those are skipped.
func Extract(text, repo string) []string {
	var found []match
	add := func(re *regexp.Regexp, group int, toURL func(m string) string) {
		for _, loc := range re.FindAllStringSubmatchIndex(text, -1) {
			if url := toURL(text[loc[2*group]:loc[2*group+1]]); url != "" {
				found = append(found, match{loc[0], url})
			}
		}
	}

	add(urlRe, 0, identity)
	add(wwwRe, 0, func(m string) string { return "http://" + m })
	add(ipRe, 0, func(m string) string { return "http://" + m })
	add(gitRe, 1, func(m string) string {
		return "https://" + strings.TrimSuffix(strings.ReplaceAll(m, ":", "/"), ".git")
	})
	add(ghRe, 1, func(m string) string { return "https://github.com/" + m })
	add(npmRe, 1, func(m string) string { return "https://npmjs.com/package/" + m })
	if repo != "" {
		add(issueRe, 1, func(m string) string { return "https://github.com/" + repo + "/issues/" + m })
	}

	found = slices.DeleteFunc(found, func(m match) bool { return insideURL(text, m.pos) })
	slices.SortStableFunc(found, func(a, b match) int { return a.pos - b.pos })

	seen := map[string]bool{}
	var urls []string
	for _, m := range found {
		if !seen[m.url] {
			seen[m.url] = true
			urls = append(urls, m.url)
		}
	}
	return urls
}

func insideURL(text string, pos int) bool {
	start := strings.LastIndexAny(text[:pos], " \t\n\"'`<>()[]") + 1
	return strings.Contains(text[start:pos], "://")
}

// GitHubRepo returns "org/repo" for a GitHub remote URL in https or ssh form.
func GitHubRepo(remote string) (string, bool) {
	m := githubRe.FindStringSubmatch(strings.TrimSpace(remote))
	if m == nil {
		return "", false
	}
	return m[1], true
}

func identity(s string) string { return s }
