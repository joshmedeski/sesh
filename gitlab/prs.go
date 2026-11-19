package gitlab

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/joshmedeski/sesh/v2/forge"
)

// PrsAsIssues fetches many merge requests from one repo in as few requests as possible,
// returning them as forge.Issue entries. This is intended for legacy GitHub-centric workflows
// that treat PRs as issues.
func (g *RealGitlab) PrsAsIssues(host string, repo string, numbers []int) (map[int]forge.Issue, []int, error) {
	found := make(map[int]forge.Issue)
	var missing []int
	for _, batch := range batchNumbers(numbers, issueBatchSize) {
		batchFound, batchMissing, err := g.mrBatch(host, repo, batch)
		if err != nil {
			return nil, nil, err
		}
		for number, mr := range batchFound {
			found[number] = mr.toIssue()
		}
		missing = append(missing, batchMissing...)
	}
	return found, missing, nil
}

// mrBatch fetches a batch of merge requests by number using one aliased GraphQL query,
// instead of one `glab mr view` subprocess per MR.
func (g *RealGitlab) mrBatch(host string, repo string, numbers []int) (map[int]gitlabMergeRequest, []int, error) {
	query, aliases := buildMrQuery(numbers)

	out, runErr := g.shell.Cmd(
		"glab", "api", "graphql", "--hostname", host,
		"--raw-field", "query="+query, "--raw-field", "fullPath="+repo,
	)
	if runErr != nil {
		return nil, nil, runErr
	}
	var resp struct {
		Data struct {
			Project map[string]*gitlabMergeRequest `json:"project"`
		} `json:"data"`
		Errors []graphqlError `json:"errors"`
	}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		return nil, nil, fmt.Errorf("decoding glab graphql response: %w", err)
	}
	found := make(map[int]gitlabMergeRequest, len(numbers))
	var missing []int
	for alias, number := range aliases {
		if pr := resp.Data.Project[alias]; pr != nil {
			found[number] = *pr
			continue
		}
		missing = append(missing, number)
	}
	sort.Ints(missing)
	return found, missing, nil
}

// mrSelection is the field set requested per merge request number.
const mrSelection = "{ iid author { username } title description state }"

// buildMrQuery builds a query that aliases each merge request number and returns
// the alias-to-number mapping needed to read the response back.
func buildMrQuery(numbers []int) (string, map[string]int) {
	aliases := make(map[string]int, len(numbers))

	var fields strings.Builder
	for _, number := range numbers {
		alias := mrAlias(number)
		aliases[alias] = number
		fmt.Fprintf(&fields, ` %s: mergeRequest(iid: "%d") %s`, alias, number, mrSelection)
	}
	query := "query($fullPath: ID!) { project(fullPath: $fullPath) {" + fields.String() + " } }"
	return query, aliases
}

// mrAlias names the field for a merge request. GraphQL aliases cannot start with a
// digit, hence the prefix.
func mrAlias(number int) string {
	return "mr" + strconv.Itoa(number)
}
