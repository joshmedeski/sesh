package dashboard

// branchLoadedMsg carries the current git branch for a session path, fetched
// asynchronously by both the sessions and configured sections.
type branchLoadedMsg struct {
	path   string
	branch string
}

// statusLoadedMsg carries the formatted git status for a session path, fetched
// asynchronously by both the sessions and configured sections.
type statusLoadedMsg struct {
	path   string
	status string
}
