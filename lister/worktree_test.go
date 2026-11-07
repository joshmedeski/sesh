package lister

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/cache"
	"github.com/joshmedeski/sesh/v2/forge"
	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/github"
	"github.com/joshmedeski/sesh/v2/gitlab"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/shell"
)

func worktreeTestConfig() model.Config {
	return model.Config{
		WorktreeConfigs: []model.WorktreeConfig{{Repo: "Nutiliti/nutiliti", Path: "~/c/nu", WorktreeDir: "w"}},
	}
}

// forgeTestResolver hands the same driver to every config, so a test exercises the real
// key implementation instead of a mock that could agree with itself.
func forgeTestResolver(f forge.Forge) forge.Resolver {
	return forge.NewResolver(
		func(cfg model.WorktreeConfig) (forge.Forge, error) {
			return f, nil
		},
		func(path string) (forge.Forge, error) {
			return f, nil
		},
	)
}

func testGithubForge(t *testing.T) forge.Forge {
	t.Helper()
	return github.NewGithub(shell.NewMockShell(t), git.NewMockGit(t))
}

func testGitlabForge(t *testing.T) forge.Forge {
	t.Helper()
	return gitlab.NewGitlab(shell.NewMockShell(t), git.NewMockGit(t))
}

func TestWorktreeResolver(t *testing.T) {
	issues := cache.NewNamespaceInDir[forge.Issue](t.TempDir(), "github-issues", 1, time.Hour)
	entries := cache.Entries[forge.Issue]{}
	entries.Put("github.com/Nutiliti/nutiliti#409", forge.Issue{Number: 409, Title: "Fix bill split rounding", State: "OPEN"})
	require.NoError(t, issues.Save(entries))

	resolve := WorktreeResolver(worktreeTestConfig(), iconTestHome(t), forgeTestResolver(testGithubForge(t)), issues)

	want := &model.WorktreeInfo{Repo: "Nutiliti/nutiliti", Number: 409, Title: "Fix bill split rounding", State: "OPEN", Name: "409", RefType: model.RefTypeIssue}
	assert.Equal(t, want, resolve("/home/user/c/nu/w/409"))
	assert.Equal(t, want, resolve("/home/user/c/nu/w/409/apps/web"))
	assert.Equal(t, &model.WorktreeInfo{Repo: "Nutiliti/nutiliti", Number: 12, Name: "12", RefType: model.RefTypeIssue}, resolve("/home/user/c/nu/w/12"),
		"an uncached issue still reports its number")
	assert.Nil(t, resolve("/home/user/c/nu"))
	assert.Nil(t, resolve("/home/user/c/nu/w/notes"))
	assert.Nil(t, resolve("/home/user/c/nu/wiki/409"))
	assert.Nil(t, resolve(""))
}

// The worktree command writes its titles under Forge.IssueKey and the lister reads them
// back under the same method. This pins that agreement: a key built anywhere else, or one
// that drops the host, would leave the lister looking up a key nothing ever wrote.
func TestWorktreeResolver_UsesTheForgeDriverKey(t *testing.T) {
	driver := testGithubForge(t)
	issues := cache.NewNamespaceInDir[forge.Issue](t.TempDir(), "github-issues", 1, time.Hour)

	key := driver.IssueKey("github.com", "Nutiliti/nutiliti", 409, false)
	entries := cache.Entries[forge.Issue]{}
	entries.Put(key, forge.Issue{Number: 409, Title: "written by the worktree command", State: "OPEN"})
	require.NoError(t, issues.Save(entries))

	resolve := WorktreeResolver(worktreeTestConfig(), iconTestHome(t), forgeTestResolver(driver), issues)

	info := resolve("/home/user/c/nu/w/409")
	require.NotNil(t, info)
	assert.Equal(t, "written by the worktree command", info.Title,
		"the lister must read the key the driver produces")
	assert.Equal(t, "github.com/Nutiliti/nutiliti#409", key,
		"and that key must carry the host, or two hosts would share an entry")
}

// A GitLab merge request's worktree directory is "mr-101", which an integer parse
// cannot read: the row would be dropped, and its title lives under the merge
// request's own key rather than the issue's.
func TestWorktreeResolver_GitlabMergeRequest(t *testing.T) {
	driver := testGitlabForge(t)
	issues := cache.NewNamespaceInDir[forge.Issue](t.TempDir(), "gitlab-issues", 1, time.Hour)
	entries := cache.Entries[forge.Issue]{}
	entries.Put(driver.IssueKey("gitlab.com", "inkscape/inkscape", 101, true),
		forge.Issue{Number: 101, Title: "Fix bezier rounding", State: "OPEN"})
	require.NoError(t, issues.Save(entries))

	config := model.Config{
		WorktreeConfigs: []model.WorktreeConfig{{
			Repo: "inkscape/inkscape", Path: "~/c/ink", WorktreeDir: "w", Forge: model.ForgeGitlab,
		}},
	}
	resolve := WorktreeResolver(config, iconTestHome(t), forgeTestResolver(driver), issues)

	want := &model.WorktreeInfo{
		Repo:    "inkscape/inkscape",
		Number:  101,
		Name:    "mr-101",
		RefType: model.RefTypePr,
		Title:   "Fix bezier rounding",
		State:   "OPEN",
	}
	assert.Equal(t, want, resolve("/home/user/c/ink/w/mr-101"))
	assert.Equal(t, want, resolve("/home/user/c/ink/w/mr-101/apps/web"))
	assert.Equal(t, &model.WorktreeInfo{Repo: "inkscape/inkscape", Number: 101, Name: "101", RefType: model.RefTypeIssue}, resolve("/home/user/c/ink/w/101"),
		"a plain issue number does not resolve to the merge request")
	assert.Nil(t, resolve("/home/user/c/ink/w/mr-notes"),
		"a name the driver cannot read is not a worktree")
}

func TestWorktreeResolver_NilWithoutWorktrees(t *testing.T) {
	assert.Nil(t, WorktreeResolver(model.Config{}, iconTestHome(t), nil, nil))
}

func TestFormatJsonResolvesWorktree(t *testing.T) {
	l := NewLister(worktreeTestConfig(), iconTestHome(t), nil, nil, nil, nil, forgeTestResolver(testGithubForge(t)), nil)
	sessions := model.SeshSessions{
		OrderedIndex: []string{"zoxide:409", "zoxide:other"},
		Directory: model.SeshSessionMap{
			"zoxide:409":   {Src: "zoxide", Name: "nu/w/409", Path: "/home/user/c/nu/w/409"},
			"zoxide:other": {Src: "zoxide", Name: "other", Path: "/home/user/c/other"},
		},
	}

	got, err := l.Format(sessions, ListOptions{Json: true})
	require.NoError(t, err)
	assert.Equal(t, &model.WorktreeInfo{Repo: "Nutiliti/nutiliti", Number: 409, Name: "409", RefType: model.RefTypeIssue}, got.Directory["zoxide:409"].Worktree)
	assert.Nil(t, got.Directory["zoxide:other"].Worktree)
}
