# dmon

### instalation

download module using

```
go get "github.com/roertbb/dmon"
```

### usage

NewEnv creates distributed Env for group of nodes, that allows exchanging messages between them and creating Monitors

```go
env, err := dmon.NewEnv(hostAddress, otherHostsAddresses...)
```

NewMonitor creates new Monitor for defined distributed Env

```go
monitor := env.NewMonitor()
```

RegisterSharedData registers shared data for defined Monitor (comma-separated pointers for variables)

```go
monitor.RegisterSpharedData(&data, otherData...)
```

Enter requests to enter distributed critical section for defined Monitor

```go
monitor.Enter()
```

Exit leaves distributed critical section for defined Monitor

```go
monitor.Exit()
```

Synchronized is method imitating synchronized block. It's higher order function that takes Monitor and function without argument, that's executed within distributed critical section defined by Monitor

```go
dmon.Synchronized(monitor)(func() {
    // ...
})
```

NewConditional creates new Conditional variable for defined Monitor

```go
conditional := monitor.NewConditional()
```

Wait waits on Conditional variables

```go
conditional.Wait()
```

Notify sends signal message to one of the processes waiting on Conditional variable

```go
conditional.Notify()
```

NotifyAll sends signal message to all of the processes waiting on Conditional variable

```go
conditional.NotifyAll()
```

### concept

- Suzuki-Kasami algorithm used for mutual exclusion of critical section - read/write operations allowed in critical section
- shared data is passed with the token, that grant access to critical section
- transparent access for shared data (requires registering pointer for shared data)
- multiple conditionals for single monitor
- multiple monitors for single group of nodes

### example

Examples can be found in `examples` directory

```go
package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/roertbb/dmon"
)

const maxBufSize = 5

func produce(buf *[]int, mon *dmon.Monitor, empty *dmon.Conditional, full *dmon.Conditional, value int) {
	// dmon.Synchronized(mon)(...)
	// is equivalent to
	// mon.Enter()
	// ...
	// mon.Exit()

	dmon.Synchronized(mon)(func() {
		// waits on full Conditional, when buf has 5 elements
		for len(*buf) == maxBufSize {
			full.Wait()
		}

		fmt.Println("produce", *buf, "<-", value)
		// append elem to buf - transparent access!
		*buf = append(*buf, value)

		// notifies processes waiting on empty condition
		empty.Notify()
	})
}

func consume(buf *[]int, mon *dmon.Monitor, empty *dmon.Conditional, full *dmon.Conditional) {
	// Enters distributed critical section
	mon.Enter()

	// waits on empty Conditional, when buf is empty
	for len(*buf) == 0 {
		empty.Wait()
	}

	// assigns first element to variable v and updates slice, removing first element from it - transparent access!
	v := (*buf)[0]
	*buf = (*buf)[1:]
	fmt.Println("consume", v, "<-", *buf)

	// notifies processes waiting on full condition
	full.Notify()

	// Leaves distributed critical section
	mon.Exit()
}

func main() {
	procType := os.Args[1]
	myAddress := os.Args[2]
	addresses := os.Args[3:]

	// tranparent slice definition as shared buffer
	buf := []int{}

	// create distributed Env responsible for communication between nodes
	env, _ := dmon.NewEnv(myAddress, addresses...)
	// creates prod-cons monitor
	mon := env.NewMonitor()

	// registers shared data - in that case slice (dynamic array) where data is stored
	mon.RegisterSharedData(&buf)

	// creates conditional variables for empty and full cases
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

```
