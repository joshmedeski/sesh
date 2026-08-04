package github

import (
	"encoding/json"
	"fmt"
	"maps"
	"sort"
	"strconv"
	"strings"
)

// issueBatchSize caps how many issues go into a single GraphQL query. Each
// alias is one node, so the ceiling is query length rather than GitHub's node
// limit, but chunking keeps a directory with hundreds of worktrees from
// building one enormous query.
const issueBatchSize = 100

// Issues fetches issues by number using one aliased GraphQL query per batch,
// instead of one `gh issue view` subprocess per issue.
func (g *RealGithub) Issues(repo string, numbers []int) (map[int]Issue, []int, error) {
	owner, name, ok := splitRepo(repo)
	if !ok {
		return nil, nil, fmt.Errorf("repo must be in owner/name form, got %q", repo)
	}

	found := make(map[int]Issue)
	var missing []int
	for _, batch := range batchNumbers(numbers, issueBatchSize) {
		batchFound, batchMissing, err := g.issueBatch(owner, name, batch)
		if err != nil {
			return nil, nil, err
		}
		maps.Copy(found, batchFound)
		missing = append(missing, batchMissing...)
	}
	return found, missing, nil
}

func (g *RealGithub) issueBatch(owner, name string, numbers []int) (map[int]Issue, []int, error) {
	query, aliases := buildIssueQuery(numbers)

	// A batch where some numbers do not exist makes gh exit non-zero while
	// still printing a body that holds the issues that did resolve, so this
	// reads stdout rather than trusting the exit code.
	out, runErr := g.shell.CmdCapture(
		"gh", "api", "graphql",
		"-f", "query="+query,
		"-f", "owner="+owner,
		"-f", "name="+name,
	)
	if out == "" {
		if runErr != nil {
			return nil, nil, fmt.Errorf("gh api graphql: %w", runErr)
		}
		return nil, nil, fmt.Errorf("gh api graphql: empty response")
	}

	var resp struct {
		Data struct {
			Repository map[string]*Issue `json:"repository"`
		} `json:"data"`
		Errors []graphqlError `json:"errors"`
	}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		// A body we cannot read is only meaningful alongside why gh failed.
		if runErr != nil {
			return nil, nil, fmt.Errorf("gh api graphql: %w", runErr)
		}
		return nil, nil, fmt.Errorf("decoding gh graphql response: %w", err)
	}

	notFound := notFoundAliases(resp.Errors)

	found := make(map[int]Issue, len(numbers))
	var missing []int
	for alias, number := range aliases {
		if issue := resp.Data.Repository[alias]; issue != nil {
			found[number] = *issue
			continue
		}
		// Only a NOT_FOUND proves the issue is really gone. Any other null —
		// a rate limit, a permissions error — is transient, so leave it out
		// of both results rather than caching an absence that isn't real.
		if notFound[alias] {
			missing = append(missing, number)
		}
	}
	sort.Ints(missing)
	return found, missing, nil
}

// issueSelection is the field set requested per number. issueOrPullRequest is
// used rather than issue because worktrees are keyed on PR numbers too (see the
// prCheckout path in the worktree package), and issue(number:) reports a plain
// NOT_FOUND for a pull request — which would get cached as "does not exist".
// The union needs an inline fragment per member; both expose number/title/state.
const issueSelection = "{ ... on Issue { number title state } ... on PullRequest { number title state } }"

// buildIssueQuery builds a query that aliases each issue number, and returns
// the alias-to-number mapping needed to read the response back.
func buildIssueQuery(numbers []int) (string, map[string]int) {
	aliases := make(map[string]int, len(numbers))

	var fields strings.Builder
	for _, number := range numbers {
		alias := issueAlias(number)
		aliases[alias] = number
		fmt.Fprintf(&fields, " %s: issueOrPullRequest(number: %d) %s", alias, number, issueSelection)
	}

	query := "query($owner: String!, $name: String!) { repository(owner: $owner, name: $name) {" +
		fields.String() + " } }"
	return query, aliases
}

// issueAlias names the field for an issue. GraphQL aliases cannot start with a
// digit, hence the prefix.
func issueAlias(number int) string {
	return "i" + strconv.Itoa(number)
}

// graphqlError is one entry of a GraphQL response's "errors" array. Path
// elements are a mix of strings and ints, hence []any.
type graphqlError struct {
	Type string `json:"type"`
	Path []any  `json:"path"`
}

// notFoundAliases collects the aliases GitHub reported as NOT_FOUND. Paths look
// like ["repository", "i409"], so the alias is the last element.
func notFoundAliases(errs []graphqlError) map[string]bool {
	notFound := make(map[string]bool)
	for _, e := range errs {
		if e.Type != "NOT_FOUND" || len(e.Path) == 0 {
			continue
		}
		if alias, ok := e.Path[len(e.Path)-1].(string); ok {
			notFound[alias] = true
		}
	}
	return notFound
}

// splitRepo splits "owner/name", rejecting anything with a missing or extra
// segment.
func splitRepo(repo string) (owner, name string, ok bool) {
	owner, name, found := strings.Cut(repo, "/")
	if !found || owner == "" || name == "" || strings.Contains(name, "/") {
		return "", "", false
	}
	return owner, name, true
}

// batchNumbers dedupes and sorts numbers, then splits them into chunks of at
// most size. Sorting keeps the generated query deterministic.
func batchNumbers(numbers []int, size int) [][]int {
	seen := make(map[int]bool, len(numbers))
	unique := make([]int, 0, len(numbers))
	for _, n := range numbers {
		if !seen[n] {
			seen[n] = true
			unique = append(unique, n)
		}
	}
	sort.Ints(unique)

	var batches [][]int
	for start := 0; start < len(unique); start += size {
		end := min(start+size, len(unique))
		batches = append(batches, unique[start:end])
	}
	return batches
}
