package main

// example execution for 3 processes
// go run counter.go 127.0.0.1:3001 127.0.0.1:3002 127.0.0.1:3003
// go run counter.go 127.0.0.1:3002 127.0.0.1:3001 127.0.0.1:3003
// go run counter.go 127.0.0.1:3003 127.0.0.1:3001 127.0.0.1:3002

import (
	"fmt"
	"os"
	"time"

	"github.com/roertbb/dmon"
)

func main() {
	myAddress := os.Args[1]
	otherNodes := os.Args[2:]

	buf := []int{0}

	env, _ := dmon.NewEnv(myAddress, otherNodes...)
	mon := env.NewMonitor()

	mon.RegisterSharedData(&buf)

	for i := 0; i < 5; i++ {
		mon.Enter()
		fmt.Println("Entered critical section")

		if len(buf) == 5 {
			buf = buf[1:]
		}
		buf = append(buf, buf[len(buf)-1]+1)

		time.Sleep(time.Second)
		fmt.Println("buff: ", buf)

		mon.Exit()
		fmt.Println("Left critical section")
		time.Sleep(time.Second)
	}
}
