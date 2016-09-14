package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time" // or "runtime"
)

type signalFunc func(s os.Signal, arg interface{})

type Hook struct {
    m map[os.Signal]signalFunc
}

func NewHook() *Hook {
    ss := new(Hook)
    ss.m = make(map[os.Signal]signalFunc)
    return ss
}

func (set *Hook) register(s os.Signal, handler signalFunc) {
    if _, found := set.m[s]; !found {
        set.m[s] = handler
    }
}

func (set *Hook) handle(sig os.Signal, arg interface{}) (err error) {
    if _, found := set.m[sig]; found {
        set.m[sig](sig, arg)
        return nil
    } else {
        return fmt.Errorf("No handler available for signal %v", sig)
    }

    panic("won't reach here")
}

func main() {
    go sysSignalHandleDemo()
    time.Sleep(time.Hour) // make the main goroutine wait!
}

func sysSignalHandleDemo() {
    ss := NewHook()
    handler := func(s os.Signal, arg interface{}) {
        fmt.Printf("handle signal: %v\n", s)
    }

    ss.register(os.Interrupt, handler)
    ss.register(os.Kill, handler)
    ss.register(syscall.SIGUSR1, handler)
    ss.register(syscall.SIGUSR2, handler)

    for {
        c := make(chan os.Signal)
        //		var sigs []os.Signal
        //		for sig := range ss.m {
        //			sigs = append(sigs, sig)
        //		}
        signal.Notify(c)
        sig := <-c

        err := ss.handle(sig, nil)
        if err != nil {
            fmt.Printf("unknown signal received: %v\n", sig)
            //			os.Exit(1)
        }
    }
}
