package tmux

import (
	"fmt"

	"github.com/joshmedeski/sesh/model"
)

func (t *RealTmux) SwitchOrAttach(name string, opts model.ConnectOpts) (string, error) {
	if opts.Switch || t.IsAttached() {
		if _, err := t.SwitchClient(name); err != nil {
			return "", fmt.Errorf("failed to switch to tmux session: %w", err)
		} else {
			return fmt.Sprintf("switching to existing tmux session: %s", name), nil
		}
	} else {
		if _, err := t.AttachSession(name); err != nil {
			return "", fmt.Errorf("failed to attach to tmux session: %w", err)
		} else {
			return fmt.Sprintf("attaching to existing tmux session: %s", name), nil
		}
	}
}
