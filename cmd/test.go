// (c) 2020 Gon Y Yi. <https://gonyyi.com/copyright.txt>

package main

import (
	"github.com/gonyyi/graceful"
	"syscall"
	"time"
)

func main() {
	graceful.New(func() {
		switch graceful.SignalReceived {
		case syscall.SIGKILL:
			println("Received SIGKILL")
		case syscall.SIGINT:
			println("Received SIGINT")
		default:
			println("Received Others")
		}
		return 0
	})

	go syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	time.Sleep(time.Second * 2)
	println("Normal termination of a process")

}
