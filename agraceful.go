// MIT License
//
// Copyright (c) 2020-2021 Gon Y. Yi <https://gonyyi.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package agraceful

import (
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
)

// As signal (termination) can be global
var signalReceived os.Signal

// GetSignal will return received os.Signal. This to be used when terminated.
func GetSignal() os.Signal {
	return signalReceived
}

// New takes a function to execute when terminated. This should NOT be used in a goroutine
func IfTerm(f func()) {
	graceful := make(chan os.Signal)
	// If found os.Interrupt, os.Kill, etc.. send a signal to m.graceful.
	signal.Notify(graceful, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTRAP, syscall.SIGHUP, syscall.SIGQUIT)

	go func() {
		signalReceived = <-graceful
		// If DoFinal function is not defined, use exit with code 1.
		if f != nil {
			f()
		}
		return
	}()
	return
}

// IfPanic should be used inside go routine or any function with defer method.
func IfPanic(f func()) {
	if r := recover(); r != nil {
		if f != nil {
			f()
		}
	}
}

// GetStack will returns stack trace in []byte. This can be used with Recover
func GetStack() []byte {
	return debug.Stack()
}

