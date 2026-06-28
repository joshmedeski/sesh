package refresher

import (
	"log/slog"
	"os"
	"os/exec"
	"syscall"
)

// Refresher launches a detached `sesh status --refresh` to repopulate the
// status cache without blocking the caller.
type Refresher interface {
	// Spawn launches a detached refresh for path. An empty path lets the child
	// resolve the directory itself (attached tmux session, then cwd).
	Spawn(path string) error
}

type RealRefresher struct{}

func NewRefresher() Refresher {
	return &RealRefresher{}
}

func refreshArgs(path string) []string {
	args := []string{"status", "--refresh"}
	if path != "" {
		args = append(args, path)
	}
	return args
}

func (r *RealRefresher) Spawn(path string) error {
	self, err := os.Executable()
	if err != nil {
		slog.Debug("refresher: os.Executable failed", "error", err)
		return err
	}
	cmd := exec.Command(self, refreshArgs(path)...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil
	// Detach into its own session so it outlives this (foreground) process.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := cmd.Start(); err != nil {
		slog.Debug("refresher: start failed", "error", err)
		return err
	}
	// Do not Wait — release the child to run independently.
	return cmd.Process.Release()
}
