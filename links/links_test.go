package links

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtract(t *testing.T) {
	text := `see https://example.com/a?b=1, and https://example.com/a?b=1 again
visit www.sesh.dev/docs or 192.168.1.10:8080/health
remote git@github.com:joshmedeski/sesh.git and ssh://git@gitlab.com/org/repo
plugin "folke/lazy.nvim" fixes #42 not https://x.dev/page#7
import { useState } from "react"
import x from './local'`

	assert.Equal(t, []string{
		"https://example.com/a?b=1",
		"http://www.sesh.dev/docs",
		"http://192.168.1.10:8080/health",
		"https://github.com/joshmedeski/sesh",
		"https://gitlab.com/org/repo",
		"https://github.com/folke/lazy.nvim",
		"https://github.com/joshmedeski/sesh/issues/42",
		"https://x.dev/page#7",
		"https://npmjs.com/package/react",
	}, Extract(text, "joshmedeski/sesh"))
}

func TestExtractSkipsIssuesWithoutRepo(t *testing.T) {
	assert.Empty(t, Extract("fixes #42", ""))
}

func TestGitHubRepo(t *testing.T) {
	for _, remote := range []string{
		"git@github.com:joshmedeski/sesh.git",
		"https://github.com/joshmedeski/sesh.git\n",
		"https://github.com/joshmedeski/sesh",
	} {
		repo, ok := GitHubRepo(remote)
		assert.True(t, ok, remote)
		assert.Equal(t, "joshmedeski/sesh", repo)
	}
	_, ok := GitHubRepo("git@gitlab.com:org/repo.git")
	assert.False(t, ok)
}
