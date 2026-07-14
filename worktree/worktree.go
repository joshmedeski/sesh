package worktree

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/joshmedeski/sesh/v2/connector"
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/github"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/pathwrap"
)

type Worktree interface {
	// Create adds (or reconnects to) a worktree for an issue/PR and connects.
	Create(opts model.WorktreeCreateOpts) (string, error)
}

type RealWorktree struct {
	config    model.Config
	git       git.Git
	github    github.Github
	connector connector.Connector
	home      home.Home
	os        oswrap.Os
	path      pathwrap.Path
}

func NewWorktree(
	config model.Config,
	g git.Git,
	gh github.Github,
	c connector.Connector,
	h home.Home,
	os oswrap.Os,
	p pathwrap.Path,
) Worktree {
	return &RealWorktree{config, g, gh, c, h, os, p}
}

// createPlan captures the decisions derived from gh before touching git.
type createPlan struct {
	key        int  // issue or PR number used for the worktree dir/branch
	prCheckout bool // create detached + `gh pr checkout` (vs. new branch)
}

func (w *RealWorktree) Create(opts model.WorktreeCreateOpts) (string, error) {
	cfg, err := w.resolveConfig(opts.Repo)
	if err != nil {
		return "", err
	}

	repoPath, err := w.home.ExpandPath(cfg.Path)
	if err != nil {
		return "", err
	}

	if fetchEnabled(cfg) {
		w.git.Fetch(repoPath) // best-effort; a failed fetch shouldn't block
	}

	plan, err := w.planCreate(cfg, opts)
	if err != nil {
		return "", err
	}

	root := w.worktreeRoot(cfg, repoPath)
	target := w.path.Join(root, strconv.Itoa(plan.key))

	if _, statErr := w.os.Stat(target); statErr == nil {
		// Existing worktree: refresh a PR checkout, then reconnect.
		if plan.prCheckout {
			if _, err := w.github.PrCheckout(target, cfg.Name, opts.Number); err != nil {
				return "", err
			}
			w.git.Pull(target) // best-effort ff-only
		}
	} else {
		if err := w.os.MkdirAll(root, 0o755); err != nil {
			return "", err
		}
		if plan.prCheckout {
			if _, err := w.git.WorktreeAddDetached(repoPath, target, baseBranch(cfg)); err != nil {
				return "", err
			}
			if _, err := w.github.PrCheckout(target, cfg.Name, opts.Number); err != nil {
				return "", err
			}
		} else {
			branch := renderBranch(cfg, plan.key)
			if _, err := w.git.WorktreeAdd(repoPath, target, branch, baseBranch(cfg)); err != nil {
				return "", err
			}
		}
	}

	return w.connector.Connect(target, model.ConnectOpts{
		Switch:  opts.Switch,
		Command: cfg.StartupCommand,
	})
}

func (w *RealWorktree) planCreate(cfg model.WorktreeConfig, opts model.WorktreeCreateOpts) (createPlan, error) {
	pr, found, err := w.github.PrView(cfg.Name, opts.Number)
	if err != nil {
		return createPlan{}, err
	}
	isPr := found || opts.Pr
	if !isPr {
		return createPlan{key: opts.Number, prCheckout: false}, nil
	}

	me, err := w.github.CurrentUser()
	if err != nil {
		return createPlan{}, err
	}
	if pr.Author != "" && pr.Author != me {
		// Foreign PR: detached checkout keyed on the PR number.
		return createPlan{key: opts.Number, prCheckout: true}, nil
	}

	// Own PR: prefer the closing issue; else scan title/body for #N.
	issue := 0
	if len(pr.ClosingIssues) > 0 {
		issue = pr.ClosingIssues[0]
	} else {
		issue = firstIssueRef(pr.Title + "\n" + pr.Body)
	}
	if issue > 0 {
		return createPlan{key: issue, prCheckout: false}, nil
	}
	return createPlan{key: opts.Number, prCheckout: true}, nil
}

func (w *RealWorktree) resolveConfig(repoOverride string) (model.WorktreeConfig, error) {
	if repoOverride != "" {
		for _, c := range w.config.WorktreeConfigs {
			if c.Name == repoOverride {
				return c, nil
			}
		}
		return model.WorktreeConfig{}, fmt.Errorf("no [[worktree]] config found for repo %q", repoOverride)
	}

	cwd, err := w.os.Getwd()
	if err != nil {
		return model.WorktreeConfig{}, err
	}
	mainWt, err := w.mainWorktree(cwd)
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
	dir := cfg.WorktreeDir
	if dir == "" {
		dir = ".wk"
	}
	if strings.HasPrefix(dir, "/") {
		return dir
	}
	return w.path.Join(repoPath, dir)
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
