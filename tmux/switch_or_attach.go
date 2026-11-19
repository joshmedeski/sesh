package tmux

import (
	"fmt"

	"github.com/joshmedeski/sesh/v2/model"
)

func (t *RealTmux) SwitchOrAttach(name string, opts model.ConnectOpts) (string, error) {
	if opts.Switch || t.IsAttached() {
		if _, err := t.switchClient(name); err != nil {
			return "", fmt.Errorf("failed to switch to tmux session: %w", err)
		} else {
			if opts.Message != "" {
				_, _ = t.DisplayMessage(opts.Message)
			}
			return fmt.Sprintf("switching to tmux session: %s", name), nil
		}
	} else {
		if _, err := t.AttachSession(name, opts.Message); err != nil {
			return "", fmt.Errorf("failed to attach to tmux session: %w", err)
		} else {
			return fmt.Sprintf("attaching to tmux session: %s", name), nil
		}
	}
}
