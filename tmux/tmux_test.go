package tmux

import (
	"os"
	"testing"
)

func TestIsRunningTrue(t *testing.T) {
	os.Setenv("TMUX", "/private/tmp/tmux-501/default,39832,2")
	expected := true
	result := IsRunning()
	if result != expected {
		t.Errorf("IsTmuxRunning() = %v; want %v", result, expected)
	}
	os.Unsetenv("TMUX")
}

func TestIsRunningFalse(t *testing.T) {
	os.Unsetenv("TMUX")
	expected := false
	result := IsRunning()
	if result != expected {
		t.Errorf("IsTmuxRunning() = %v; want %v", result, expected)
	}
}
