// +build windows

package graceful

import "fmt"

func RequestStop() error {
	//return syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	return fmt.Errorf("No syscall.Kill support for Windows")
}
