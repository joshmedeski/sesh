package gitlab

import (
	"encoding/json"
	"fmt"
	"maps"
	"sort"
	"strconv"
	"strings"

	"github.com/joshmedeski/sesh/v2/forge"
)

// issueBatchSize caps how many issues/MRs go into a single GraphQL query. Each
// alias is one node, so the ceiling is query length rather than GitLab's node
// limit, but chunking keeps a directory with hundreds of worktrees from
// building one enormous query.
const issueBatchSize = 100

// Issues fetches issues by number using one aliased GraphQL query per batch,
// instead of one `glab issue view` subprocess per issue.
func (g *RealGitlab) Issues(host string, repo string, numbers []int) (map[int]forge.Issue, []int, error) {
	found := make(map[int]forge.Issue)
	var missing []int
	for _, batch := range batchNumbers(numbers, issueBatchSize) {
		batchFound, batchMissing, err := g.issueBatch(host, repo, batch)
		if err != nil {
			return nil, nil, err
		}
		maps.Copy(found, batchFound)
		missing = append(missing, batchMissing...)
	}
	return found, missing, nil
}

func (g *RealGitlab) issueBatch(host string, repo string, numbers []int) (map[int]forge.Issue, []int, error) {
	query, aliases := buildIssueQuery(numbers)

	out, runErr := g.shell.Cmd(
		"glab", "api", "graphql", "--hostname", host,
		"--raw-field", "query="+query, "--raw-field", "fullPath="+repo,
	)
	if runErr != nil {
		return nil, nil, runErr
	}
	var resp struct {
		Data struct {
			Project map[string]*gitlabIssue `json:"project"`
		} `json:"data"`
		Errors []graphqlError `json:"errors"`
	}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		return nil, nil, fmt.Errorf("decoding glab graphql response: %w", err)
	}
	found := make(map[int]forge.Issue, len(numbers))
	var missing []int
	for alias, number := range aliases {
		if issue := resp.Data.Project[alias]; issue != nil {
			found[number] = issue.toIssue()
			continue
		}
		missing = append(missing, number)
	}
	sort.Ints(missing)
	return found, missing, nil
}

// issueSelection is the field set requested per number.
const issueSelection = "{ iid title state }"

// buildIssueQuery builds a query that aliases each issue number and returns
// the alias-to-number mapping needed to read the response back.
func buildIssueQuery(numbers []int) (string, map[string]int) {
	aliases := make(map[string]int, len(numbers))

	var fields strings.Builder
	for _, number := range numbers {
		alias := issueAlias(number)
		aliases[alias] = number
		fmt.Fprintf(&fields, ` %s: issue(iid: "%d") %s`, alias, number, issueSelection)
	}
	query := "query($fullPath: ID!) { project(fullPath: $fullPath) {" + fields.String() + " } }"
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
