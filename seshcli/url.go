package seshcli

import (
	"cmp"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/joshmedeski/sesh/v2/github"
	"github.com/joshmedeski/sesh/v2/links"
	"github.com/joshmedeski/sesh/v2/picker"
	"github.com/joshmedeski/sesh/v2/worktree"
)

const defuddleCommand = "defuddle parse {url} --markdown"

func NewURLCommand(base *BaseDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "url",
		Short: "Pick URLs from the current tmux pane and open them in the browser",
		Long: `List every link in the current tmux pane in a picker and open the chosen ones.

Enter opens the highlighted URL; tab marks several to open at once. A bare #1234
links to an issue in the [[worktree]] repo holding the session or pane, or else the GitHub
origin of the pane's directory. --history includes that many lines of scrollback.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}

			history := deps.Config.URL.History
			if cmd.Flags().Changed("history") {
				history, _ = cmd.Flags().GetInt("history")
			}
			content, err := deps.Tmux.CaptureCurrentPane(history)
			if err != nil {
				return fmt.Errorf("couldn't capture the current pane: %w", err)
			}

			urls := links.Extract(content, issueRepo(deps))
			if filter := deps.Config.URL.ExtraFilter; filter != "" {
				extras, err := runExtraFilter(filter, content)
				if err != nil {
					return fmt.Errorf("extra_filter failed: %w", err)
				}
				urls = appendUnique(urls, extras)
			}
			return pickURLs(deps, urls)
		},
	}

	cmd.Flags().Int("history", 0, "lines of scrollback to search besides the visible screen")

	return cmd
}

func pickURLs(deps *Deps, urls []string) error {
	if len(urls) == 0 {
		fmt.Println("no URLs found")
		return nil
	}
	chosen, err := deps.Picker.PickURLs(urls, urlPreview(deps))
	if err != nil {
		return err
	}
	for _, url := range chosen {
		if err := deps.Browser.Open(url); err != nil {
			return err
		}
	}
	return nil
}

func issueRepo(deps *Deps) string {
	if sessionPath, err := deps.Tmux.FocusedSessionPath(); err == nil {
		if repo, ok := deps.Worktree.RepoAt(sessionPath); ok {
			return repo
		}
	}
	panePath, err := deps.Tmux.CurrentPanePath()
	if err != nil {
		return ""
	}
	if repo, ok := deps.Worktree.RepoAt(panePath); ok {
		return repo
	}
	origin, err := deps.Git.OriginURL(panePath)
	if err != nil {
		return ""
	}
	repo, _ := links.GitHubRepo(origin)
	return repo
}

func runExtraFilter(filter, content string) ([]string, error) {
	c := exec.Command(cmp.Or(os.Getenv("SHELL"), "/bin/sh"), "-c", filter)
	c.Stdin = strings.NewReader(content)
	out, err := c.Output()
	if err != nil {
		return nil, err
	}
	var lines []string
	for line := range strings.SplitSeq(string(out), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			lines = append(lines, line)
		}
	}
	return lines, nil
}

func appendUnique(urls, extras []string) []string {
	seen := make(map[string]bool, len(urls))
	for _, url := range urls {
		seen[url] = true
	}
	for _, extra := range extras {
		if !seen[extra] {
			seen[extra] = true
			urls = append(urls, extra)
		}
	}
	return urls
}

func urlPreview(deps *Deps) picker.PreviewFunc {
	command := deps.Config.URL.PreviewCommand
	if command == "" {
		if _, err := deps.Exec.LookPath("defuddle"); err == nil {
			command = defuddleCommand
		}
	}
	return func(url string) (string, error) {
		if repo, number, _, ok := worktree.ParseGitHubRef(url); ok {
			if detail, err := deps.Github.IssueDetail(repo, number); err == nil {
				return issueMarkdown(detail), nil
			}
			return url, nil
		}
		if command == "" {
			return url, nil
		}
		out, err := deps.Shell.ShellCmd(command, map[string]string{"{url}": url})
		if err != nil || strings.TrimSpace(out) == "" {
			return url, nil
		}
		return out, nil
	}
}

func issueMarkdown(detail github.IssueDetail) string {
	return fmt.Sprintf("# %s\n\n**%s**\n\n%s", detail.Title, detail.State, detail.Body)
}
