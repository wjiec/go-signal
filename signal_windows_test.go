//go:build windows
// +build windows

package signal

import (
	"errors"
	"syscall"
)

const (
	SigUsr1 = syscall.SIGHUP
	SigUsr2 = syscall.SIGINT
)

func SendSignalUser1(_ int) error {
	return errors.New("testing on windows is not supported for now")
}

func SendSignalUser2(_ int) error {
	return errors.New("testing on windows is not supported for now")
}
