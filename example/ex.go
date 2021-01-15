package main

import (
	"github.com/gonyyi/agraceful"
	"os"
	"strings"
	"time"
)

func main() {
	// Although IfTerm can be used in goroutine, it shouldn't.
	agraceful.IfTerm(func(){
		println("1: being forced to shut down...")
		if s := agraceful.GetSignal(); s != nil {
			println("2: signal received: ", s.String())
		}
		os.Exit(0)
	})

	go test()

	time.Sleep(time.Second*10)
	println("3: all finished")
}

func test() {
	// IfPanic always should be used with "defer"
	defer agraceful.IfPanic(func(){
		// GetStack() will pull stack trace info.
		println("4: test: panicked..",
			strings.ReplaceAll(
				string(agraceful.GetStack()[:20]), "\n", ""))
	})

	// intentionally cause panic... this can be replaced with panic(nil)
	var a []bool
	println(a[2])

	println("5: OK test ")
}
