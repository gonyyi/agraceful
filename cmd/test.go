// (c) 2020 Gon Y Yi. <https://gonyyi.com/copyright.txt>

package main

import (
	"github.com/gonyyi/graceful"
	"syscall"
	"time"
)

func main() {
	graceful.New(func() int {
		switch graceful.SignalReceived {
		case syscall.SIGKILL:
			println("Received SIGKILL")
			return 0
		case syscall.SIGINT:
			println("Received SIGINT")
			return 0
		default:
			println("Received Others")
			return 1
		}
		return 0
	})

	go syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	time.Sleep(time.Second * 2)
	println("Normal termination of a process")

}
