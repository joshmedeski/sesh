package links

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarkdown(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/article":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(`<html><head><title>Sesh docs</title></head><body>
<nav><a href="/">Home</a></nav>
<article><h1>Sesh docs</h1>
<p>Sesh is a <strong>smart</strong> terminal session manager that helps you create and manage tmux sessions quickly and easily using zoxide.</p>
<p>Read the <a href="/guide">guide</a> to get started with sessions, windows, worktrees and the picker, which lists everything in one place.</p>
</article></body></html>`))
		case "/json":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	md, err := Markdown(srv.URL + "/article")
	require.NoError(t, err)
	assert.Contains(t, md, "# Sesh docs")
	assert.Contains(t, md, "**smart**")
	assert.Contains(t, md, "[guide]("+srv.URL+"/guide)")

	_, err = Markdown(srv.URL + "/json")
	assert.Error(t, err)
	_, err = Markdown(srv.URL + "/missing")
	assert.Error(t, err)
	_, err = Markdown("ftp://example.com/file")
	assert.Error(t, err)
}
