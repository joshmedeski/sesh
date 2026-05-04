package live

import (
	"os"
	"syscall"
)

// realProcessChecker 用 kill -0 (signal 0) 探活，跨 Unix 平台通用。
type realProcessChecker struct{}

func NewProcessChecker() ProcessChecker { return realProcessChecker{} }

func (realProcessChecker) IsAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	if err := p.Signal(syscall.Signal(0)); err != nil {
		// ESRCH = no such process; EPERM = exists but not ours, 也算活
		if err == syscall.EPERM {
			return true
		}
		return false
	}
	return true
}
