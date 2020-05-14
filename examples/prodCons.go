package main

import (
	"fmt"
	"os"
	"time"

	"github.com/roertbb/dmon"
)

const maxBufSize = 5

func produce(buf *[]int, mon *dmon.Monitor, empty *dmon.Conditional, full *dmon.Conditional, value int) {
	dmon.Synchronized(mon)(func() {
		for len(*buf) == maxBufSize {
			full.Wait()
		}

		*buf = append(*buf, value)
		fmt.Println("produce", value, *buf)

		empty.Notify()
	})
}

func consume(buf *[]int, mon *dmon.Monitor, empty *dmon.Conditional, full *dmon.Conditional) {
	mon.Enter()

	for len(*buf) == 0 {
		empty.Wait()
	}

	v := (*buf)[0]
	*buf = (*buf)[1:]
	fmt.Println("consume", *buf, v)

	full.Notify()

	mon.Exit()
}

func main() {
	procType := os.Args[1]
	myAddress := os.Args[2]
	addresses := os.Args[3:]

	buf := []int{}

	env, _ := dmon.NewEnv(myAddress, addresses...)
	mon, _ := env.NewMonitor()

	mon.RegisterSharedData(&buf)

	empty := mon.NewConditional()
	full := mon.NewConditional()

	if procType == "prod" {
		for i := 0; true; i++ {
			produce(&buf, mon, empty, full, i)
			time.Sleep(time.Second)
		}
	} else if procType == "cons" {
		for i := 0; true; i++ {
			consume(&buf, mon, empty, full)
			time.Sleep(time.Second)
		}
	}
}

// go run prodcons.go prod 127.0.0.1:5001 127.0.0.1:5002 127.0.0.1:5003
// go run prodcons.go prod 127.0.0.1:5002 127.0.0.1:5001 127.0.0.1:5003
// go run prodcons.go cons 127.0.0.1:5003 127.0.0.1:5001 127.0.0.1:5002
