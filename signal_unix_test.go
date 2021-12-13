//go:build linux || unix || openbsd || darwin
// +build linux unix openbsd darwin

package signal

import (
	"syscall"
)

const (
	SigUsr1 = syscall.SIGUSR1
	SigUsr2 = syscall.SIGUSR2
)

func SendSignalUser1(pid int) error {
	return syscall.Kill(pid, SigUsr1)
}

func SendSignalUser2(pid int) error {
	return syscall.Kill(pid, SigUsr2)
}
