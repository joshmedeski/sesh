package forge

import (
	"errors"
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssueState(t *testing.T) {
	assert.Equal(t, IssueState("OPEN"), IssueStateOpen)
	assert.Equal(t, IssueState("CLOSED"), IssueStateClosed)
	assert.Equal(t, IssueState("MERGED"), IssueStateMerged)
	assert.Equal(t, IssueState(""), IssueStateUnknown)
}

func TestHost(t *testing.T) {
	tests := []struct {
		name        string
		cfg         model.WorktreeConfig
		want        string
		wantErrText string
	}{
		{
			name: "custom host takes precedence",
			cfg: model.WorktreeConfig{
				Host:  "gitlab.internal.corp",
				Forge: model.ForgeGitlab,
			},
			want: "gitlab.internal.corp",
		},
		{
			name: "custom host without forge",
			cfg: model.WorktreeConfig{
				Host: "custom.domain.org",
			},
			want: "custom.domain.org",
		},
		{
			name: "default forge resolves to github.com",
			cfg:  model.WorktreeConfig{},
			want: "github.com",
		},
		{
			name: "explicit github forge resolves to github.com",
			cfg: model.WorktreeConfig{
				Forge: model.ForgeGithub,
			},
			want: "github.com",
		},
		{
			name: "explicit gitlab forge resolves to gitlab.com",
			cfg: model.WorktreeConfig{
				Forge: model.ForgeGitlab,
			},
			want: "gitlab.com",
		},
		{
			name: "invalid forge returns error",
			cfg: model.WorktreeConfig{
				Forge: model.Forge("bitbucket"),
			},
			wantErrText: `invalid forge "bitbucket"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Host(tt.cfg)
			if tt.wantErrText != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrText)
				assert.Empty(t, got)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRealResolver(t *testing.T) {
	mockForge := NewMockForge(t)
	expectedErr := errors.New("resolver error")

	t.Run("FromConfig success", func(t *testing.T) {
		cfg := model.WorktreeConfig{Repo: "org/repo"}
		r := NewResolver(
			func(c model.WorktreeConfig) (Forge, error) {
				assert.Equal(t, cfg, c)
				return mockForge, nil
			},
			nil,
		)
		got, err := r.FromConfig(cfg)
		require.NoError(t, err)
		assert.Equal(t, mockForge, got)
	})

	t.Run("FromConfig error from closure", func(t *testing.T) {
		cfg := model.WorktreeConfig{Repo: "org/repo"}
		r := NewResolver(
			func(c model.WorktreeConfig) (Forge, error) {
				return nil, expectedErr
			},
			nil,
		)
		got, err := r.FromConfig(cfg)
		require.ErrorIs(t, err, expectedErr)
		assert.Nil(t, got)
	})

	t.Run("FromConfig nil func returns error", func(t *testing.T) {
		var r RealResolver
		got, err := r.FromConfig(model.WorktreeConfig{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "FromConfigFunc is not implemented")
		assert.Nil(t, got)
	})

	t.Run("FromPath success", func(t *testing.T) {
		path := "/path/to/repo"
		r := NewResolver(
			nil,
			func(p string) (Forge, error) {
				assert.Equal(t, path, p)
				return mockForge, nil
			},
		)
		got, err := r.FromPath(path)
		require.NoError(t, err)
		assert.Equal(t, mockForge, got)
	})

	t.Run("FromPath error from closure", func(t *testing.T) {
		path := "/path/to/repo"
		r := NewResolver(
			nil,
			func(p string) (Forge, error) {
				return nil, expectedErr
			},
		)
		got, err := r.FromPath(path)
		require.ErrorIs(t, err, expectedErr)
		assert.Nil(t, got)
	})

	t.Run("FromPath nil func returns error", func(t *testing.T) {
		var r RealResolver
		got, err := r.FromPath("/path/to/repo")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "FromPathFunc is not implemented")
		assert.Nil(t, got)
	})
}

func TestParseIssueNumber(t *testing.T) {
	cases := []struct {
		branch  string
		want    string
		wantHas bool
	}{
		{"400", "400", true},
		{"feat/400-status-bar", "400", true},
		{"bugfix/issue-12", "12", true},
		{"feat/400-then-401", "400", true},
		{"1234-some-feature", "1234", true},
		{"feat/42-add-login", "42", true},
		{"issue-99", "99", true},
		{"main", "", false},
		{"no-numbers-here", "", false},
		{"", "", false},
	}
	for _, c := range cases {
		t.Run(c.branch, func(t *testing.T) {
			got, has := ParseIssueNumber(c.branch)
			assert.Equal(t, c.wantHas, has)
			assert.Equal(t, c.want, got)
		})
	}
}
