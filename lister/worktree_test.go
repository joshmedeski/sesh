package lister

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/cache"
	"github.com/joshmedeski/sesh/v2/github"
	"github.com/joshmedeski/sesh/v2/model"
)

func worktreeTestConfig() model.Config {
	return model.Config{
		WorktreeConfigs: []model.WorktreeConfig{{Repo: "Nutiliti/nutiliti", Path: "~/c/nu", WorktreeDir: "w"}},
	}
}

func TestWorktreeResolver(t *testing.T) {
	issues := cache.NewNamespaceInDir[github.Issue](t.TempDir(), "github-issues", 1, time.Hour)
	entries := cache.Entries[github.Issue]{}
	entries.Put(github.IssueKey("Nutiliti/nutiliti", 409), github.Issue{Number: 409, Title: "Fix bill split rounding", State: "OPEN"})
	require.NoError(t, issues.Save(entries))

	resolve := worktreeResolver(worktreeTestConfig(), iconTestHome(t), issues)

	want := &model.WorktreeInfo{Repo: "Nutiliti/nutiliti", Number: 409, Title: "Fix bill split rounding", State: "OPEN"}
	assert.Equal(t, want, resolve("/home/user/c/nu/w/409"))
	assert.Equal(t, want, resolve("/home/user/c/nu/w/409/apps/web"))
	assert.Equal(t, &model.WorktreeInfo{Repo: "Nutiliti/nutiliti", Number: 12}, resolve("/home/user/c/nu/w/12"),
		"an uncached issue still reports its number")
	assert.Nil(t, resolve("/home/user/c/nu"))
	assert.Nil(t, resolve("/home/user/c/nu/w/notes"))
	assert.Nil(t, resolve("/home/user/c/nu/wiki/409"))
	assert.Nil(t, resolve(""))
}

func TestWorktreeResolver_NilWithoutWorktrees(t *testing.T) {
	assert.Nil(t, worktreeResolver(model.Config{}, iconTestHome(t), nil))
}

func TestFormatJsonResolvesWorktree(t *testing.T) {
	l := NewLister(worktreeTestConfig(), iconTestHome(t), nil, nil, nil, nil, nil)
	sessions := model.SeshSessions{
		OrderedIndex: []string{"zoxide:409", "zoxide:other"},
		Directory: model.SeshSessionMap{
			"zoxide:409":   {Src: "zoxide", Name: "nu/w/409", Path: "/home/user/c/nu/w/409"},
			"zoxide:other": {Src: "zoxide", Name: "other", Path: "/home/user/c/other"},
		},
	}

	got, err := l.Format(sessions, ListOptions{Json: true})
	require.NoError(t, err)
	assert.Equal(t, &model.WorktreeInfo{Repo: "Nutiliti/nutiliti", Number: 409}, got.Directory["zoxide:409"].Worktree)
	assert.Nil(t, got.Directory["zoxide:other"].Worktree)
}
