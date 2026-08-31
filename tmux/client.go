package tmux

import (
	"strconv"
	"strings"

	"github.com/joshmedeski/sesh/v2/model"
)

// seshClientEnv names the tmux client explicitly, for the cases sesh cannot
// work out on its own: a wrapper that strips $TMUX, or a tmux binding that
// already knows the answer.
//
//	bind-key s run-shell 'SESH_CLIENT=#{client_name} sesh connect foo'
const seshClientEnv = "SESH_CLIENT"

const clientFormat = "#{client_name}\t#{client_tty}\t#{session_id}\t#{client_activity}"

func (t *RealTmux) ListClients() ([]model.TmuxClient, error) {
	lines, err := t.shell.ListCmd(t.bin, "list-clients", "-F", clientFormat)
	if err != nil {
		return nil, err
	}
	clients := make([]model.TmuxClient, 0, len(lines))
	for _, line := range lines {
		fields := strings.Split(line, "\t")
		if len(fields) != 4 || fields[0] == "" {
			continue
		}
		activity, _ := strconv.ParseInt(fields[3], 10, 64)
		clients = append(clients, model.TmuxClient{
			Name:      fields[0],
			TTY:       fields[1],
			SessionID: fields[2],
			Activity:  activity,
		})
	}
	return clients, nil
}

// ResolveClient reports which tmux client sesh should act on.
//
// A tmux popup has no pane identity — it gets no $TMUX_PANE — so a bare
// `switch-client -t` leaves tmux to work out who is asking. It falls back to
// the most recently active client on the whole server, which may well be
// attached to an unrelated session, and sesh then moves a terminal the user
// isn't looking at. $TMUX still carries the id of the session the popup was
// launched from, which is enough to name the right client explicitly.
//
// Returns "" when no client can be determined; callers then let tmux guess,
// which is no worse than not resolving at all.
func (t *RealTmux) ResolveClient() string {
	if client := t.os.Getenv(seshClientEnv); client != "" {
		return client
	}
	clients, err := t.ListClients()
	if err != nil || len(clients) == 0 {
		return ""
	}
	if sessionID := currentSessionID(t.os.Getenv("TMUX")); sessionID != "" {
		if client := mostRecent(clients, sessionID); client != "" {
			return client
		}
	}
	return mostRecent(clients, "")
}

// switchClient moves the resolved client, falling back to letting tmux pick
// one when sesh cannot name it.
func (t *RealTmux) switchClient(targetSession string) (string, error) {
	if client := t.ResolveClient(); client != "" {
		return t.SwitchClientTarget(client, targetSession)
	}
	return t.SwitchClient(targetSession)
}

// currentSessionID pulls the session id out of $TMUX, which tmux sets to
// "<socket path>,<server pid>,<session number>". The number is the session id
// without its "$" sigil.
func currentSessionID(tmuxEnv string) string {
	fields := strings.Split(tmuxEnv, ",")
	if len(fields) < 3 || fields[2] == "" {
		return ""
	}
	return "$" + fields[2]
}

// mostRecent returns the name of the most recently active client, limited to
// those attached to sessionID when it is non-empty. When several clients share
// a session, the one touched last is the one the user is most likely watching.
func mostRecent(clients []model.TmuxClient, sessionID string) string {
	name := ""
	var activity int64
	for _, client := range clients {
		if sessionID != "" && client.SessionID != sessionID {
			continue
		}
		if name == "" || client.Activity > activity {
			name = client.Name
			activity = client.Activity
		}
	}
	return name
}
