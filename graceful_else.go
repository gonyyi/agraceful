// +build !windows

package graceful

import "syscall"

// If a user to terminal a program, he/she can also use function (RequestShutdown()) anywhere.
// I think only works with non-Windows.

func RequestStop() error {
	return syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}
