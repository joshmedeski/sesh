package tmux

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommand_Run(t *testing.T) {
	c := Command{
		execFunc: func(string, []string) (string, error) {
			return "stub", nil
		},
	}

	res, err := c.Run([]string{"arg1", "arg2"})
	require.NoError(t, err)
	require.Equal(t, "stub", res)
}

func TestCommand_GetSession(t *testing.T) {
	testCases := map[string]struct {
		MockResponse string
		MockError    error
		SessionName  string
		Error        error
	}{
		"happy path": {MockResponse: sessionList, SessionName: "dotfiles"},
		"unhappy path": {
			MockResponse: sessionList,
			SessionName:  "not found",
			Error:        ErrNotFound,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			command := &Command{
				execFunc: func(string, []string) (string, error) {
					return tc.MockResponse, tc.MockError
				},
			}
			session, err := command.GetSession(tc.SessionName)
			require.ErrorIs(t, err, tc.Error)
			if err != nil {
				return
			}
			require.Equal(t, tc.SessionName, session.Name())
		})
	}
}
