package main

// usage: go run prodcons.go <prod | cons> <myAddress> <addresses of other nodes>

// example usage of 2 producers and 1 consumer
// go run prodcons.go prod 127.0.0.1:3001 127.0.0.1:3002 127.0.0.1:3003
// go run prodcons.go prod 127.0.0.1:3002 127.0.0.1:3001 127.0.0.1:3003
// go run prodcons.go cons 127.0.0.1:3003 127.0.0.1:3001 127.0.0.1:3002

import (
	"fmt"
	"math/rand"
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

		fmt.Println("produce", *buf, "<-", value)
		*buf = append(*buf, value)

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
	fmt.Println("consume", v, "<-", *buf)

	full.Notify()

	mon.Exit()
}

func main() {
	procType := os.Args[1]
	myAddress := os.Args[2]
	addresses := os.Args[3:]

	buf := []int{}

	env, _ := dmon.NewEnv(myAddress, addresses...)
	mon := env.NewMonitor()

	mon.RegisterSharedData(&buf)

	empty := mon.NewConditional()
	full := mon.NewConditional()

	if procType == "prod" {
		for i := 0; true; i++ {
			time.Sleep(time.Second * time.Duration(rand.Intn(2)+1))
			produce(&buf, mon, empty, full, i)
		}
	} else if procType == "cons" {
		for i := 0; true; i++ {
			consume(&buf, mon, empty, full)
			time.Sleep(time.Second * time.Duration(rand.Intn(2)+1))
		}
	}
}
