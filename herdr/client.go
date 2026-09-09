package herdr

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"path/filepath"
	"time"

	"github.com/joshmedeski/sesh/v2/oswrap"
)

// Client sends one request to the herdr socket API and returns the result
// object. It is the only thing the backend needs from herdr, so tests fake it
// with canned responses instead of a socket.
type Client interface {
	Call(method string, params any) (json.RawMessage, error)
}

// SocketClient speaks herdr's newline-delimited JSON protocol over the local
// socket, one connection per request. sesh is a short-lived CLI, so there is
// nothing to gain from keeping a connection open.
type SocketClient struct {
	path    string
	timeout time.Duration
}

func NewSocketClient(path string) *SocketClient {
	return &SocketClient{path: path, timeout: 5 * time.Second}
}

// SocketPath resolves the socket the way herdr's own CLI does: an explicit
// HERDR_SOCKET_PATH wins, then a named session under HERDR_SESSION, then the
// default session's socket.
func SocketPath(os oswrap.Os) string {
	if path := os.Getenv("HERDR_SOCKET_PATH"); path != "" {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	base := filepath.Join(home, ".config", "herdr")
	if session := os.Getenv("HERDR_SESSION"); session != "" {
		return filepath.Join(base, "sessions", session, "herdr.sock")
	}
	return filepath.Join(base, "herdr.sock")
}

type request struct {
	ID     string `json:"id"`
	Method string `json:"method"`
	Params any    `json:"params"`
}

type response struct {
	ID     string          `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (c *SocketClient) Call(method string, params any) (json.RawMessage, error) {
	if params == nil {
		params = map[string]any{}
	}
	conn, err := net.DialTimeout("unix", c.path, c.timeout)
	if err != nil {
		return nil, fmt.Errorf("herdr: couldn't reach %s: %w", c.path, err)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(c.timeout))

	body, err := json.Marshal(request{ID: "sesh", Method: method, Params: params})
	if err != nil {
		return nil, err
	}
	if _, err := conn.Write(append(body, '\n')); err != nil {
		return nil, fmt.Errorf("herdr %s: %w", method, err)
	}

	reader := bufio.NewReaderSize(conn, 1<<20)
	line, err := reader.ReadBytes('\n')
	if err != nil && len(line) == 0 {
		return nil, fmt.Errorf("herdr %s: %w", method, err)
	}
	return decodeResponse(method, line)
}

// decodeResponse unwraps herdr's envelope, turning a protocol-level error into
// a Go error so callers only ever see a result.
func decodeResponse(method string, line []byte) (json.RawMessage, error) {
	var resp response
	if err := json.Unmarshal(line, &resp); err != nil {
		return nil, fmt.Errorf("herdr %s: bad response: %w", method, err)
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("herdr %s: %s: %s", method, resp.Error.Code, resp.Error.Message)
	}
	return resp.Result, nil
}

// call runs a request and decodes its result into out.
func call(c Client, method string, params any, out any) error {
	raw, err := c.Call(method, params)
	if err != nil {
		return err
	}
	if out == nil {
		return nil
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("herdr %s: decode result: %w", method, err)
	}
	return nil
}
