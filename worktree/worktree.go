package worktree

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/joshmedeski/sesh/v2/browser"
	"github.com/joshmedeski/sesh/v2/cache"
	"github.com/joshmedeski/sesh/v2/connector"
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/github"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/jj"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/pathwrap"
)

type Worktree interface {
	// Connect attaches to the worktree for an issue/PR, adding it if absent.
	Connect(opts model.WorktreeConnectOpts) (string, error)
	// List returns the worktrees for a repo, each paired with its issue title.
	List(opts model.WorktreeListOpts) ([]model.WorktreeEntry, error)
	// BrowseURL returns the GitHub URL for the worktree containing
	// sessionPath: its issue, or with pr the pull request for its branch.
	BrowseURL(sessionPath string, pr bool) (string, error)
	// RepoAt returns the GitHub repo of the [[worktree]] block whose root
	// contains path.
	RepoAt(path string) (string, bool)
}

type RealWorktree struct {
	config    model.Config
	git       git.Git
	jj        jj.Jj
	github    github.Github
	connector connector.Connector
	browser   browser.Browser
	home      home.Home
	os        oswrap.Os
	path      pathwrap.Path
	// issues caches fetched issue titles so listing worktrees does not hit
	// the network on every keystroke-adjacent invocation.
	issues *cache.Namespace[github.Issue]
}

func NewWorktree(
	config model.Config,
	g git.Git,
	j jj.Jj,
	gh github.Github,
	c connector.Connector,
	b browser.Browser,
	h home.Home,
	os oswrap.Os,
	p pathwrap.Path,
	issues *cache.Namespace[github.Issue],
) Worktree {
	return &RealWorktree{config, g, j, gh, c, b, h, os, p, issues}
}

// connectPlan captures the decisions derived from gh before touching git.
type connectPlan struct {
	key        int  // issue or PR number used for the worktree dir/branch
	prCheckout bool // create detached + `gh pr checkout` (vs. new branch)
	pr         github.PullRequest
}

func (w *RealWorktree) Connect(opts model.WorktreeConnectOpts) (string, error) {
	if opts.FromBrowser {
		resolved, err := w.resolveFromBrowser(opts)
		if err != nil {
			return "", err
		}
		opts = resolved
	}

	cfg, err := w.resolveConfig(opts.Repo)
	if err != nil {
		return "", err
	}

	repoPath, err := w.home.ExpandPath(cfg.Path)
	if err != nil {
		return "", err
	}

	isJj := w.isJj(repoPath)
	if fetchEnabled(cfg) {
		if isJj {
			_, _ = w.jj.Fetch(repoPath)
		} else {
			_, _ = w.git.Fetch(repoPath)
		}
	}

	plan, err := w.planConnect(cfg, opts)
	if err != nil {
		return "", err
	}

	root := w.worktreeRoot(cfg, repoPath)
	target := w.path.Join(root, strconv.Itoa(plan.key))

	created := false
	if _, statErr := w.os.Stat(target); statErr == nil {
		// Existing worktree: refresh a PR checkout, then reconnect.
		if plan.prCheckout && isJj {
			if err := w.jjRefreshPr(repoPath, target, plan.pr); err != nil {
				return "", err
			}
		} else if plan.prCheckout {
			if _, err := w.github.PrCheckout(target, cfg.Repo, opts.Number); err != nil {
				return "", err
			}
			_, _ = w.git.Pull(target)
		}
	} else {
		created = true
		if err := w.os.MkdirAll(root, 0o755); err != nil {
			return "", err
		}
		if isJj {
			if err := w.jjAdd(cfg, repoPath, target, plan); err != nil {
				return "", err
			}
		} else if plan.prCheckout {
			if _, err := w.git.WorktreeAddDetached(repoPath, target, baseBranch(cfg)); err != nil {
				return "", err
			}
			if _, err := w.github.PrCheckout(target, cfg.Repo, opts.Number); err != nil {
				return "", err
			}
		} else {
			branch := renderBranch(cfg, plan.key)
			if _, err := w.git.WorktreeAdd(repoPath, target, branch, baseBranch(cfg)); err != nil {
				return "", err
			}
		}
	}

	sessionName, err := w.connector.Connect(target, model.ConnectOpts{
		Switch:     opts.Switch,
		Command:    connectCommand(cfg, created),
		Background: opts.Background || opts.Command != "",
	})
	if err != nil || opts.Command == "" {
		return sessionName, err
	}
	return w.connector.ConnectWindow(model.WindowConnectOpts{
		Name:       commandProgram(opts.Command),
		New:        true,
		Session:    sessionName,
		Command:    opts.Command,
		Background: opts.Background,
		Switch:     opts.Switch,
	})
}

func commandProgram(command string) string {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return ""
	}
	return filepath.Base(fields[0])
}

// connectCommand picks the command the new session runs. The two are exclusive:
// a worktree being created is a one-time setup — installing dependencies,
// seeding a .env — that has nothing to do with what should run every time you
// come back to a worktree that already exists. Running both would mean
// create_command could never be the whole of what happens on creation.
//
// Either way the connector only sends it when the tmux session is new, so
// reattaching to a live session runs nothing.
func connectCommand(cfg model.WorktreeConfig, created bool) string {
	if created {
		return cfg.CreateCommand
	}
	return cfg.StartupCommand
}

// resolveFromBrowser reads the active browser tab URL and fills Number/Repo/Pr
// from the GitHub issue/PR it points at.
func (w *RealWorktree) resolveFromBrowser(opts model.WorktreeConnectOpts) (model.WorktreeConnectOpts, error) {
	url, ok, err := w.browser.ActiveTabURL()
	if err != nil {
		return opts, err
	}
	if !ok {
		return opts, fmt.Errorf("browser tab lookup unavailable: set [browser].application (macOS only)")
	}
	repo, number, isPR, ok := ParseGitHubRef(url)
	if !ok {
		return opts, fmt.Errorf("could not parse a GitHub issue/PR from browser URL %q", url)
	}
	if opts.Repo == "" {
		opts.Repo = repo
	}
	opts.Number = number
	opts.Pr = opts.Pr || isPR
	return opts, nil
}

func (w *RealWorktree) planConnect(cfg model.WorktreeConfig, opts model.WorktreeConnectOpts) (connectPlan, error) {
	pr, found, err := w.github.PrView(cfg.Repo, opts.Number)
	if err != nil {
		return connectPlan{}, err
	}
	isPr := found || opts.Pr
	if !isPr {
		return connectPlan{key: opts.Number, prCheckout: false}, nil
	}

	me, err := w.github.CurrentUser()
	if err != nil {
		return connectPlan{}, err
	}
	if pr.Author != "" && pr.Author != me {
		// Foreign PR: detached checkout keyed on the PR number.
		return connectPlan{key: opts.Number, prCheckout: true, pr: pr}, nil
	}

	// Own PR: prefer the closing issue; else scan title/body for #N.
	issue := 0
	if len(pr.ClosingIssues) > 0 {
		issue = pr.ClosingIssues[0]
	} else {
		issue = firstIssueRef(pr.Title + "\n" + pr.Body)
	}
	if issue > 0 {
		return connectPlan{key: issue, prCheckout: false}, nil
	}
	return connectPlan{key: opts.Number, prCheckout: true, pr: pr}, nil
}

func (w *RealWorktree) isJj(dir string) bool {
	info, err := w.os.Stat(w.path.Join(dir, ".jj"))
	return err == nil && info.IsDir()
}

func (w *RealWorktree) jjAdd(cfg model.WorktreeConfig, repoPath, target string, plan connectPlan) error {
	name := strconv.Itoa(plan.key)
	if !plan.prCheckout {
		if _, err := w.jj.WorkspaceAdd(repoPath, name, jjRevset(baseBranch(cfg)), target); err != nil {
			return err
		}
		_, err := w.jj.BookmarkCreate(target, renderBranch(cfg, plan.key), "@")
		return err
	}
	rev, err := w.jjFetchPr(repoPath, plan.pr)
	if err != nil {
		return err
	}
	_, err = w.jj.WorkspaceAdd(repoPath, name, rev, target)
	return err
}

func (w *RealWorktree) jjRefreshPr(repoPath, target string, pr github.PullRequest) error {
	rev, err := w.jjFetchPr(repoPath, pr)
	if err != nil {
		return err
	}
	_, err = w.jj.New(target, rev)
	return err
}

func (w *RealWorktree) jjFetchPr(repoPath string, pr github.PullRequest) (string, error) {
	if pr.HeadRefName == "" {
		return "", fmt.Errorf("could not find the head branch of the pull request")
	}
	remote := "origin"
	if pr.IsFork {
		remote = pr.HeadOwner
		_, _ = w.jj.RemoteAdd(repoPath, remote, fmt.Sprintf("https://github.com/%s/%s.git", pr.HeadOwner, pr.HeadRepo))
	}
	if _, err := w.jj.Fetch(repoPath, "--remote", remote, "--branch", pr.HeadRefName); err != nil {
		return "", err
	}
	return remoteRevset(pr.HeadRefName, remote), nil
}

func remoteRevset(branch, remote string) string {
	return fmt.Sprintf("%q@%q", branch, remote)
}

func jjRevset(base string) string {
	if strings.ContainsAny(base, "@()") {
		return base
	}
	remote, branch, ok := strings.Cut(base, "/")
	if !ok {
		return base
	}
	return remoteRevset(branch, remote)
}

func (w *RealWorktree) resolveConfig(repoOverride string) (model.WorktreeConfig, error) {
	if repoOverride != "" {
		for _, c := range w.config.WorktreeConfigs {
			if c.Repo == repoOverride {
				return c, nil
			}
		}
		return model.WorktreeConfig{}, fmt.Errorf("no [[worktree]] config found for repo %q", repoOverride)
	}

	cwd, err := w.os.Getwd()
	if err != nil {
		return model.WorktreeConfig{}, err
	}
	mainWt, err := w.repoRoot(cwd)
	if err != nil {
		return model.WorktreeConfig{}, err
	}
	for _, c := range w.config.WorktreeConfigs {
		expanded, err := w.home.ExpandPath(c.Path)
		if err != nil {
			continue
		}
		if w.sameDir(expanded, mainWt) {
			return c, nil
		}
	}
	return model.WorktreeConfig{}, fmt.Errorf("no [[worktree]] config matches current repo %q", mainWt)
}

func (w *RealWorktree) repoRoot(cwd string) (string, error) {
	if isJj, _, repo := w.jj.Roots(cwd); isJj {
		return repo, nil
	}
	return w.mainWorktree(cwd)
}

// mainWorktree returns the primary worktree path (first porcelain entry).
func (w *RealWorktree) mainWorktree(cwd string) (string, error) {
	ok, out, err := w.git.WorktreeList(cwd)
	if err != nil || !ok {
		return "", fmt.Errorf("not in a git repository: %s", cwd)
	}
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "worktree ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "worktree ")), nil
		}
	}
	return "", fmt.Errorf("could not determine main worktree from %q", cwd)
}

func (w *RealWorktree) sameDir(a, b string) bool {
	ea, err1 := w.path.EvalSymlinks(a)
	eb, err2 := w.path.EvalSymlinks(b)
	if err1 == nil && err2 == nil {
		return ea == eb
	}
	return a == b
}

func (w *RealWorktree) worktreeRoot(cfg model.WorktreeConfig, repoPath string) string {
	return cfg.Root(repoPath)
}

func (w *RealWorktree) BrowseURL(sessionPath string, pr bool) (string, error) {
	if sessionPath == "" {
		return "", fmt.Errorf("no focused tmux session")
	}
	cfg, number, dir, ok := w.worktreeAt(sessionPath)
	if !ok {
		return "", fmt.Errorf("%s is not a worktree (no [[worktree]] root contains it)", sessionPath)
	}
	if !pr {
		if cfg.Repo == "" {
			return "", fmt.Errorf(
				"the [[worktree]] block for path %q is missing `repo = \"org/repo\"`", cfg.Path)
		}
		return fmt.Sprintf("https://github.com/%s/issues/%d", cfg.Repo, number), nil
	}
	if w.isJj(dir) {
		return w.jjPrURL(cfg, number)
	}
	url, found, err := w.github.PrURL(dir)
	if err != nil {
		return "", err
	}
	if !found {
		return "", fmt.Errorf("no pull request found for worktree %d", number)
	}
	return url, nil
}

func (w *RealWorktree) jjPrURL(cfg model.WorktreeConfig, number int) (string, error) {
	if cfg.Repo == "" {
		return "", fmt.Errorf(
			"the [[worktree]] block for path %q is missing `repo = \"org/repo\"`", cfg.Path)
	}
	url, found, err := w.github.PrURLForHead(cfg.Repo, renderBranch(cfg, number))
	if err != nil || found {
		return url, err
	}
	pr, found, err := w.github.PrView(cfg.Repo, number)
	if err != nil {
		return "", err
	}
	if !found {
		return "", fmt.Errorf("no pull request found for worktree %d", number)
	}
	return pr.URL, nil
}

// worktreeAt finds the [[worktree]] block whose root contains path, and the
// numbered worktree directory under that root.
func (w *RealWorktree) worktreeAt(path string) (model.WorktreeConfig, int, string, bool) {
	path = w.resolveSymlinks(path)
	for _, cfg := range w.config.WorktreeConfigs {
		repoPath, err := w.home.ExpandPath(cfg.Path)
		if err != nil {
			continue
		}
		root := w.resolveSymlinks(w.worktreeRoot(cfg, repoPath))
		rel, err := w.path.Rel(root, path)
		if err != nil || rel == "." || strings.HasPrefix(rel, "..") {
			continue
		}
		key := strings.Split(w.path.ToSlash(rel), "/")[0]
		number, err := strconv.Atoi(key)
		if err != nil {
			continue
		}
		return cfg, number, w.path.Join(root, key), true
	}
	return model.WorktreeConfig{}, 0, "", false
}

func (w *RealWorktree) resolveSymlinks(path string) string {
	if resolved, err := w.path.EvalSymlinks(path); err == nil {
		return resolved
	}
	return path
}

var issueRefRe = regexp.MustCompile(`#(\d+)`)

func firstIssueRef(text string) int {
	m := issueRefRe.FindStringSubmatch(text)
	if len(m) < 2 {
		return 0
	}
	n, _ := strconv.Atoi(m[1])
	return n
}

// githubRefRe matches a GitHub issue or PR URL, capturing org/repo, the kind
// ("issues"|"pull"), and the number. Tolerant of scheme, trailing path,
// query, and fragment.
var githubRefRe = regexp.MustCompile(`(?:^|https?://)(?:www\.)?github\.com/([^/]+/[^/]+)/(issues|pull)/(\d+)`)

// ParseGitHubRef extracts repo, number, and kind from a GitHub issue/PR URL.
// ok is false for any URL that is not a recognizable issue/PR reference.
func ParseGitHubRef(rawURL string) (repo string, number int, isPR bool, ok bool) {
	m := githubRefRe.FindStringSubmatch(rawURL)
	if len(m) < 4 {
		return "", 0, false, false
	}
	n, err := strconv.Atoi(m[3])
	if err != nil {
		return "", 0, false, false
	}
	return m[1], n, m[2] == "pull", true
}

func fetchEnabled(cfg model.WorktreeConfig) bool {
	if cfg.Fetch == nil {
		return true
	}
	return *cfg.Fetch
}

func baseBranch(cfg model.WorktreeConfig) string {
	if cfg.BaseBranch == "" {
		return "origin/main"
	}
	return cfg.BaseBranch
}

func renderBranch(cfg model.WorktreeConfig, key int) string {
	tmpl := cfg.BranchTemplate
	if tmpl == "" {
		tmpl = "{number}"
	}
	return strings.ReplaceAll(tmpl, "{number}", strconv.Itoa(key))
}

func (w *RealWorktree) RepoAt(path string) (string, bool) {
	cfg, _, _, ok := w.worktreeAt(path)
	return cfg.Repo, ok && cfg.Repo != ""
}
