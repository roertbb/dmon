package dmon

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pebbe/zmq4"
)

// Env is responsible for communication within distributed nodes
type Env struct {
	address  string
	receiver *zmq4.Socket
	sockets  map[string]*zmq4.Socket
	monitors map[string]*Monitor
	running  bool
}

// NewEnv creates distributed Env for group of nodes, that allows exchanging messages between them and creating Monitors
func NewEnv(myAddress string, addresses ...string) (*Env, error) {
	time.Sleep(time.Second)
	fmt.Println("starting env on", myAddress)

	env := Env{
		address:  myAddress,
		sockets:  map[string]*zmq4.Socket{},
		monitors: map[string]*Monitor{},
		running:  true,
	}

	env.receiver, _ = zmq4.NewSocket(zmq4.PULL)
	env.receiver.Bind(fmt.Sprintf("tcp://%s", myAddress))

	for _, addr := range addresses {
		socket, err := zmq4.NewSocket(zmq4.PUSH)
		if err != nil {
			fmt.Println("failed to create socket", err)
		}
		env.sockets[addr] = socket
		env.sockets[addr].Connect(fmt.Sprintf("tcp://%s", addr))
	}

	go env.listener()

	return &env, nil
}

// NewMonitor creates new Monitor for defined distributed Env
func (env *Env) NewMonitor() *Monitor {
	nextID := strconv.Itoa(len(env.monitors))
	monitor := newMonitor(nextID, env)
	env.monitors[nextID] = monitor

	return monitor
}

// NewMonitor creates new distributed monitor with assigned name for distributed Env
// TODO: func (env *Env) NewMonitorWithName()

func (env *Env) getLowestAddress() string {
	lowestAddress := env.address
	for address := range env.sockets {
		if address < lowestAddress {
			lowestAddress = address
		}
	}
	return lowestAddress
}

func (env *Env) broadcast(message []byte) {
	for _, socket := range env.sockets {
		_, err := socket.SendBytes(message, zmq4.DONTWAIT)
		if err != nil {
			fmt.Println("failed to broadcast message")
		}
	}
}

func (env *Env) send(address string, message []byte) {
	_, err := env.sockets[address].SendBytes(message, zmq4.DONTWAIT)
	if err != nil {
		fmt.Println("failed to broadcast message")
	}
}

func (env *Env) listener() {
	for env.running {
		data, _ := env.receiver.RecvBytes(0)
		msg, _ := deserializeMessage(data)

		if msg.Type == requestCSMessageType {
			env.monitors[msg.Mid].handleRequestCSMessage(msg.Data)
		} else if msg.Type == tokenMessageType {
			env.monitors[msg.Mid].handleTokenMessage(msg.Data)
		} else if msg.Type == conditionalWaitMessageType {
			env.monitors[msg.Mid].handleConditionalWaitMessage(msg.Data)
		} else if msg.Type == conditionalSignalMessageType {
			env.monitors[msg.Mid].handleConditionalSignalMessage(msg.Data)
		}
	}
}
