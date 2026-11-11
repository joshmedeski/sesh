package lister

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/joshmedeski/sesh/v2/tmuxinator"
	"github.com/joshmedeski/sesh/v2/zoxide"
)

// listWithBackends are the mocked backends behind a SourceLister. A mock
// without an expectation fails the test when it is called, so a backend that
// must stay untouched needs no assertion of its own.
type listWithBackends struct {
	tmux       *tmux.MockTmux
	zoxide     *zoxide.MockZoxide
	tmuxinator *tmuxinator.MockTmuxinator
	home       *home.MockHome
}

func newListWithBackends() *listWithBackends {
	return &listWithBackends{
		tmux:       new(tmux.MockTmux),
		zoxide:     new(zoxide.MockZoxide),
		tmuxinator: new(tmuxinator.MockTmuxinator),
		home:       new(home.MockHome),
	}
}

func (b *listWithBackends) lister(config model.Config) SourceLister {
	return NewLister(config, b.home, b.tmux, b.zoxide, b.tmuxinator, nil, nil)
}

func zoxideSource(name, path string, score float64) model.SeshSessions {
	key := "zoxide:" + name
	return model.SeshSessions{
		OrderedIndex: []string{key},
		Directory: model.SeshSessionMap{
			key: {Src: "zoxide", Name: name, Path: path, Score: score},
		},
	}
}

func sessionNames(sessions model.SeshSessions) []string {
	names := make([]string, 0, len(sessions.OrderedIndex))
	for _, key := range sessions.OrderedIndex {
		names = append(names, sessions.Directory[key].Name)
	}
	return names
}

func TestListWith_TakesProvidedSourceWithoutCallingItsBackend(t *testing.T) {
	b := newListWithBackends()
	b.tmux.On("ListSessions").Return([]*model.TmuxSession{{Name: "main", Path: "/m"}}, nil)

	got, err := b.lister(model.Config{}).ListWith(
		ListOptions{Tmux: true, Zoxide: true},
		SourceResults{"zoxide": zoxideSource("code", "/code", 1)},
	)

	require.NoError(t, err)
	assert.Equal(t, []string{"main", "code"}, sessionNames(got))
	b.zoxide.AssertNotCalled(t, "ListResults")
}

func TestListWith_ProvidedEmptySourceIsNotFetched(t *testing.T) {
	b := newListWithBackends()
	b.tmux.On("ListSessions").Return([]*model.TmuxSession{{Name: "main", Path: "/m"}}, nil)

	got, err := b.lister(model.Config{}).ListWith(
		ListOptions{Tmux: true, Zoxide: true},
		SourceResults{"zoxide": {}},
	)

	require.NoError(t, err)
	assert.Equal(t, []string{"main"}, sessionNames(got))
	b.zoxide.AssertNotCalled(t, "ListResults")
}

func TestListWith_IgnoresProvidedSourceThatIsNotRequested(t *testing.T) {
	b := newListWithBackends()
	b.tmux.On("ListSessions").Return([]*model.TmuxSession{{Name: "main", Path: "/m"}}, nil)

	got, err := b.lister(model.Config{}).ListWith(
		ListOptions{Tmux: true},
		SourceResults{"zoxide": zoxideSource("code", "/code", 1)},
	)

	require.NoError(t, err)
	assert.Equal(t, []string{"main"}, sessionNames(got))
}

// A provided zoxide result that opts do not request must not feed the scores
// a merged sort_order group is ordered by.
func TestListWith_UnrequestedSourceDoesNotScoreAMergedGroup(t *testing.T) {
	config := model.Config{
		SessionConfigs: []model.SessionConfig{{Name: "project", Path: "/p"}},
		SortOrder:      model.SortOrder{[]string{"tmux", "config"}},
	}
	opts := ListOptions{Tmux: true, Config: true}
	setup := func() *listWithBackends {
		b := newListWithBackends()
		b.tmux.On("ListSessions").Return([]*model.TmuxSession{{Name: "main", Path: "/m"}}, nil)
		b.home.On("ExpandPath", "/p").Return("/p", nil)
		return b
	}

	want, err := setup().lister(config).List(opts)
	require.NoError(t, err)
	require.Equal(t, []string{"main", "project"}, sessionNames(want))

	// The provided score would put "project" ahead of "main" if it counted.
	provided := model.SeshSessions{
		OrderedIndex: []string{"zoxide:p", "zoxide:m"},
		Directory: model.SeshSessionMap{
			"zoxide:p": {Src: "zoxide", Name: "p", Path: "/p", Score: 100},
			"zoxide:m": {Src: "zoxide", Name: "m", Path: "/m", Score: 1},
		},
	}
	got, err := setup().lister(config).ListWith(opts, SourceResults{"zoxide": provided})

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestListWith_StillFetchesWindowNamesLive(t *testing.T) {
	b := newListWithBackends()
	b.tmux.On("ListSessions").Return([]*model.TmuxSession{{Name: "main", Path: "/m"}}, nil)
	b.tmux.On("ListAllWindowNames", "").Return(map[string][]string{"main": {"editor"}}, nil).Once()
	config := model.Config{}
	config.TUI.ShowWindows = true

	got, err := b.lister(config).ListWith(
		ListOptions{Tmux: true, Zoxide: true},
		SourceResults{"zoxide": {}},
	)

	require.NoError(t, err)
	assert.Equal(t, []string{"editor"}, got.Directory[got.OrderedIndex[0]].WindowNames)
	b.tmux.AssertNumberOfCalls(t, "ListAllWindowNames", 1)
}

// The window name lookup runs next to the source fetch. A failing source must
// not return before that lookup has finished, or the lookup outlives the call.
func TestListWith_SourceErrorWaitsForWindowNames(t *testing.T) {
	b := newListWithBackends()
	b.tmux.On("ListSessions").Return([]*model.TmuxSession{{Name: "main", Path: "/m"}}, nil)
	entered := make(chan struct{})
	release := make(chan struct{})
	var releaseOnce sync.Once
	releaseLookup := func() { releaseOnce.Do(func() { close(release) }) }
	t.Cleanup(releaseLookup)
	b.tmux.On("ListAllWindowNames", "").
		Run(func(mock.Arguments) {
			close(entered)
			<-release
		}).
		Return(map[string][]string{"main": {"editor"}}, nil).Once()
	b.zoxide.On("ListResults").Return(nil, errors.New("zoxide failed"))
	config := model.Config{}
	config.TUI.ShowWindows = true

	done := make(chan error, 1)
	go func() {
		_, err := b.lister(config).ListWith(ListOptions{Tmux: true, Zoxide: true}, nil)
		done <- err
	}()

	// The zoxide failure is immediate; the lookup is held open.
	<-entered
	select {
	case <-done:
		t.Fatal("ListWith returned while the window name lookup was still running")
	case <-time.After(100 * time.Millisecond):
	}
	releaseLookup()
	err := <-done

	require.Error(t, err)
	assert.Contains(t, err.Error(), "zoxide failed")
	b.tmux.AssertNumberOfCalls(t, "ListAllWindowNames", 1)
}

func TestListWith_FetchesMissingSourceFromItsBackend(t *testing.T) {
	b := newListWithBackends()
	b.tmux.On("ListSessions").Return([]*model.TmuxSession{{Name: "main", Path: "/m"}}, nil)
	b.zoxide.On("ListResults").Return([]*model.ZoxideResult{{Path: "/code", Score: 1}}, nil)
	b.home.On("ShortenHome", "/code").Return("code", nil)

	got, err := b.lister(model.Config{}).ListWith(ListOptions{Tmux: true, Zoxide: true}, SourceResults{})

	require.NoError(t, err)
	assert.Equal(t, []string{"main", "code"}, sessionNames(got))
}

func TestListSources_FetchesOnlyTheNamedSources(t *testing.T) {
	b := newListWithBackends()
	b.zoxide.On("ListResults").Return([]*model.ZoxideResult{{Path: "/code", Score: 1}}, nil)
	b.home.On("ShortenHome", "/code").Return("code", nil)
	b.tmuxinator.On("List").Return([]*model.TmuxinatorConfig{{Name: "dev"}}, nil)

	got, err := b.lister(model.Config{}).ListSources([]string{"zoxide", "tmuxinator"})

	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, []string{"zoxide:code"}, got["zoxide"].OrderedIndex)
	assert.Equal(t, []string{"tmuxinator:dev"}, got["tmuxinator"].OrderedIndex)
	b.tmux.AssertNotCalled(t, "ListSessions")
}

func TestListSources_ReturnsASourceError(t *testing.T) {
	b := newListWithBackends()
	b.zoxide.On("ListResults").Return(nil, errors.New("zoxide failed"))
	b.tmuxinator.On("List").Return([]*model.TmuxinatorConfig{}, nil)

	_, err := b.lister(model.Config{}).ListSources([]string{"zoxide", "tmuxinator"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "zoxide failed")
}

func TestListSources_RejectsAnUnknownSourceBeforeFetchingAny(t *testing.T) {
	b := newListWithBackends()

	_, err := b.lister(model.Config{}).ListSources([]string{"zoxide", "nope"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), `unknown source "nope"`)
	b.zoxide.AssertNotCalled(t, "ListResults")
}
