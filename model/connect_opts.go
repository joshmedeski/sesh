package model

type ConnectOpts struct {
	Command    string
	Switch     bool
	Tmuxinator bool
}

type WorktreeConnectOpts struct {
	Number      int
	Repo        string // GitHub "org/repo" override; empty => detect from cwd
	Pr          bool   // force the pull-request path
	Switch      bool
	FromBrowser bool // resolve Number/Repo/Pr from the active browser tab URL
}

// WorktreeListOpts selects which repo's worktrees to list. Path and Repo are
// two ways to pick the same [[worktree]] block; when both are empty the repo is
// detected from the current directory.
type WorktreeListOpts struct {
	Path string // worktree root, e.g. "~/c/nu/w"
	Repo string // GitHub "org/repo"
	// Refresh refetches every title and state, ignoring how fresh the cached
	// ones are. It is how a title edited or an issue closed since the last
	// listing is picked up before the TTL is up.
	Refresh bool
}

// WorktreeEntry is one worktree directory paired with the issue it maps to.
type WorktreeEntry struct {
	Number int    // issue/PR number, taken from the directory name
	Path   string // absolute path to the worktree
	Title  string // issue title; empty when unknown or not yet fetched
	State  string // "OPEN" | "CLOSED"; empty when unknown
}

// WindowConnectOpts identifies a window to connect to and what to run in it.
// Name is the window's identity: when it matches an existing window that window
// is reused, and when it is empty there is nothing to match on so a window is
// always created.
type WindowConnectOpts struct {
	Name string
	// New skips the match and always creates, for when the name is wanted as a
	// label rather than as an identity to reuse.
	New        bool
	Session    string // target session; empty => the attached session
	Path       string // start directory; empty => the target session's root
	Command    string
	Background bool // create without selecting it or moving the client
	Switch     bool // switch rather than attach, for triggers outside tmux
}
