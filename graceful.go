package graceful

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Graceful struct {
	signal   os.Signal
	graceful chan os.Signal
	//logger        *log.Logger
	log           *log.Logger
	tmpFinalFuncs []func(*Graceful)
	finalFuncs    []func(*Graceful)
	waitSec       int
}

func (m *Graceful) GetSignal() os.Signal {
	return m.signal
}

func (m *Graceful) SetLogger(l *log.Logger) *Graceful {
	m.log = l
	return m
}
func (m *Graceful) SetWait(seconds int) *Graceful {
	m.waitSec = seconds
	return m
}
func (m *Graceful) SetFuncs(fns ...func(*Graceful)) *Graceful {
	m.finalFuncs = fns
	return m
}

// Temporarily replace final function to this one
func (m *Graceful) SetTempFuncs(fns ...func(*Graceful)) *Graceful {
	m.tmpFinalFuncs = m.finalFuncs
	m.finalFuncs = fns
	return m
}

// This reverts back to original final function (see SetTempFuncs())
func (m *Graceful) RevertFuncs() *Graceful {
	m.finalFuncs = m.tmpFinalFuncs
	return m
}

func (m *Graceful) Stop() *Graceful {
	m.graceful <- os.Kill
	return m
}
func (m *Graceful) Recover(fn ...func(*Graceful)) *Graceful {
	if r := recover(); r != nil {
		m.log.Println(L{
			Message: MSG_GRACEFUL_FAIL_RECOVER,
		})
		for _, v := range fn {
			v(m)
		}
	}
	return m
}

func Recover(fn ...func()) {
	if r := recover(); r != nil {
		for _, v := range fn {
			v()
		}
	}
}

func NewGraceful(fn ...func(*Graceful)) *Graceful {
	m := Graceful{
		graceful: make(chan os.Signal, 2),
		//logger:     log.New(ioutil.Discard, "modme.systm.Graceful:", log.Ltime),
		log:        log.New(ioutil.Discard, "", 0),
		waitSec:    1,
		finalFuncs: fn}

	signal.Notify(m.graceful, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTRAP)

	for k, v := range m.finalFuncs {
		if v == nil {
			m.finalFuncs[k] = func(*Graceful) {}
		}
	}

	go func() {
		sig := <-m.graceful
		m.signal = sig
		{
			m.log.Println(fmt.Sprintf("Received a termination signal: %s / will be terminating in %d seconds", sig.String(), m.waitSec))
		}

		for _, v := range m.finalFuncs {
			v(&m)
		}
		time.Sleep(time.Duration(m.waitSec) * time.Second)
		os.Exit(0)
	}()

	return &m
}
