package tmux

import (
	"fmt"
	"log/slog"
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

// tmuxPaneEnv is tmux's own answer to "who is asking". Any command run from a
// normal pane inherits it, and tmux resolves the client from it exactly.
const tmuxPaneEnv = "TMUX_PANE"

// clientFormat uses the same "::" separator as every other listing rather than
// a tab, because tmux does not deliver a tab intact to a command client that
// isn't attached: it sanitizes the control character to "_", so the whole line
// arrives as a single field. That is exactly the GUI-launcher case — Leader
// Key, Raycast, an osascript binding — where every client silently dropped out
// of the list and sesh concluded nothing was attached to switch.
var clientFormat = strings.Join([]string{
	"#{client_name}",
	"#{client_tty}",
	"#{session_id}",
	"#{client_activity}",
}, separator)

const clientFieldCount = 4

func (t *RealTmux) ListClients() ([]model.TmuxClient, error) {
	lines, err := t.shell.ListCmd(t.bin, "list-clients", "-F", clientFormat)
	if err != nil {
		return nil, err
	}
	clients := make([]model.TmuxClient, 0, len(lines))
	unparsed := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Split(line, separator)
		if len(fields) != clientFieldCount || fields[0] == "" {
			unparsed = append(unparsed, line)
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
	// Clients tmux listed but sesh couldn't read must not look like a server
	// with nothing attached: callers treat that as "nothing to switch" and
	// report success without moving anything.
	if len(clients) == 0 && len(unparsed) > 0 {
		return nil, fmt.Errorf("tmux listed %d client(s) in an unreadable format: %q", len(unparsed), unparsed)
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
		slog.Debug("tmux: client named by env", "client", client, "env", seshClientEnv)
		return client
	}
	tmuxEnv := t.os.Getenv("TMUX")
	// A pane inside tmux carries $TMUX_PANE, so tmux can resolve the calling
	// client from the pane itself — exactly, not by guessing. Naming a client
	// here would replace that certainty with the heuristic below, which is only
	// ever a fallback for callers that have no pane identity. Popups are the
	// case this whole resolver exists for, and they get no $TMUX_PANE.
	//
	// $TMUX has to be set for the pane to mean anything: a GUI launcher can
	// inherit a stale $TMUX_PANE from whatever started it, and tmux has no
	// caller to resolve from outside a session anyway.
	if pane := t.os.Getenv(tmuxPaneEnv); pane != "" && tmuxEnv != "" {
		slog.Debug("tmux: leaving client resolution to tmux", "pane", pane)
		return ""
	}
	clients, err := t.ListClients()
	if err != nil {
		// Almost always the tmux binary being unreachable — a GUI launcher
		// hands sesh a bare PATH that Homebrew isn't on. Nothing downstream can
		// work, and every caller degrades to doing nothing visible, so this
		// cannot be a debug line.
		slog.Error("tmux: could not list clients", "error", err, "tmuxCommand", t.bin)
		return ""
	}
	if len(clients) == 0 {
		slog.Debug("tmux: no clients attached")
		return ""
	}
	sessionID := currentSessionID(tmuxEnv)
	if sessionID != "" {
		if client := mostRecent(clients, sessionID); client != "" {
			slog.Debug("tmux: resolved client from calling session", "client", client, "session", sessionID)
			return client
		}
	}
	client := mostRecent(clients, "")
	slog.Debug("tmux: resolved most recently active client", "client", client, "callingSession", sessionID, "clients", clients)
	return client
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
