package lister

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListTmuxPanesWithoutPaneNavigator(t *testing.T) {
	l := &RealLister{mux: mux.NewMockMux(t)}

	_, err := l.ListTmuxPanes()

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrPanesUnsupported)
}
