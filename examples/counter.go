package main

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
	mon, _ := env.NewMonitor()

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
