// (c) 2020 Gon Y Yi. <https://gonyyi.com/copyright.txt>

package graceful

import (
	"os"
	"os/signal"
	"syscall"
)

// Global variable -- so functions can differentiate (if needed..)
var SignalReceived os.Signal

type Graceful struct {
	graceful chan os.Signal
	DoFinal  func()
}

// GetStackTrace will returns stack trace []byte
func GetStackTrace() []byte {
	return debug.Stack()
}

// Recover should be used inside go routine or any function with defer method.
func Recover(f func()) {
	if r := recover(); r != nil {
		if f != nil {
			f()
		}
	}
}

// New takes a function to execute when terminated
func New(f func()) *Graceful {
	m := Graceful{
		graceful: make(chan os.Signal),
		DoFinal:  f,
	}

	// If found os.Interrupt, os.Kill, etc.. send a signal to m.graceful.
	signal.Notify(m.graceful, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTRAP, syscall.SIGHUP, syscall.SIGQUIT)

	go func() {
		SignalReceived = <-m.graceful

		// If DoFinal function is not defined, use exit with code 1.
		if m.DoFinal != nil {
			m.DoFinal()
		}
		return
	}()

	return &m
}
