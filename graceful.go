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
	DoFinal  func() int
}

func (m *Graceful) SetDoFinal(fn func() int) *Graceful {
	m.DoFinal = fn
	return m
}

func New(fn func() int) *Graceful {
	m := Graceful{
		graceful: make(chan os.Signal),
		DoFinal:  fn,
	}

	// If found os.Interrupt, os.Kill, etc.. send a signal to m.graceful.
	signal.Notify(m.graceful, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTRAP, syscall.SIGHUP, syscall.SIGQUIT)

	go func() {
		SignalReceived = <-m.graceful

		// If DoFinal function is not defined, use exit with code 1.
		if m.DoFinal != nil {
			os.Exit(m.DoFinal())
			return
		}
		os.Exit(1)
		return
	}()

	return &m
}
