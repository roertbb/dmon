package dmon

import (
	"strconv"
	"sync"
)

// Monitor is distributed monitor structure handling data, token exchange and it's Conditional variables
type Monitor struct {
	mid          string
	env          *Env
	local        *sync.Mutex
	keepToken    bool
	conditionals map[string]*Conditional
	data         *map[string]interface{}
	RN           map[string]int
	token        *token
	tokenChan    chan bool
	lastSignaled []string
}

func newMonitor(mid string, env *Env) *Monitor {
	monitor := Monitor{
		mid:          mid,
		env:          env,
		local:        &sync.Mutex{},
		keepToken:    false,
		conditionals: map[string]*Conditional{},
		data:         &map[string]interface{}{},
		RN:           map[string]int{},
		token:        nil,
		tokenChan:    make(chan bool),
		lastSignaled: []string{},
	}

	monitor.RN[env.address] = 0
	for address := range env.sockets {
		monitor.RN[address] = 0
	}

	lowestAddress := env.getLowestAddress()
	if lowestAddress == env.address {
		monitor.token = newToken(&monitor)
	}

	return &monitor
}

// RegisterSharedData registers shared data for defined Monitor (comma-separated pointers for variables)
func (mon *Monitor) RegisterSharedData(data ...interface{}) {
	for id, value := range data {
		(*mon.data)[strconv.Itoa(id)] = value
	}
}

// Enter requests to enter distributed critical section for defined Monitor
func (mon *Monitor) Enter() {
	mon.local.Lock()

	if mon.token == nil {
		mon.RN[mon.env.address]++
		requestMsg, _ := serializeRequestCSMessage(mon.env.address, mon.mid, mon.RN[mon.env.address])
		mon.env.broadcast(requestMsg)
		mon.local.Unlock()

		<-mon.tokenChan

		mon.local.Lock()
	}

	mon.keepToken = true
	mon.local.Unlock()
}

// Exit leaves distributed critical section for defined Monitor
func (mon *Monitor) Exit() {
	mon.local.Lock()

	mon.token.LRN[mon.env.address] = mon.RN[mon.env.address]
	mon.token.updateQ(mon)

	address, _ := mon.token.lastSignaledOrPop()
	if address != "" {
		mon.sendToken(address)
	}

	mon.keepToken = false

	mon.local.Unlock()
}

// NewConditional creates new Conditional variable for defined Monitor
func (mon *Monitor) NewConditional() *Conditional {
	idx := strconv.Itoa(len(mon.conditionals))
	cond := newConditional(mon, idx)
	mon.conditionals[idx] = cond

	return cond
}

func (mon *Monitor) sendToken(address string) {
	mon.token.serializeData(mon.data)
	mon.token.serializeCondWaiting(&mon.conditionals)
	tokenMsg, _ := serializeTokenMessage(mon.mid, mon.token)
	mon.env.send(address, tokenMsg)
	mon.token = nil
}

// message handlers

func (mon *Monitor) handleRequestCSMessage(data []byte) {
	requestCS, _ := deserializeRequestCSMessage(data)

	mon.local.Lock()
	if requestCS.SN > mon.RN[requestCS.From] {
		mon.RN[requestCS.From] = requestCS.SN
	}
	if mon.token != nil && !mon.keepToken && mon.RN[requestCS.From] == mon.token.LRN[requestCS.From]+1 {
		mon.sendToken(requestCS.From)
	}
	mon.local.Unlock()
}

func (mon *Monitor) handleTokenMessage(data []byte) {
	token, _ := deserializeTokenMessage(data)

	mon.local.Lock()
	mon.token = token
	mon.token.deserializeData(mon.data, mon)
	mon.token.deserializeCondWaiting(&mon.conditionals)
	mon.keepToken = true
	mon.local.Unlock()

	mon.tokenChan <- true
}

func (mon *Monitor) handleConditionalSignalMessage(data []byte) {
	signalMsg, _ := deserializeConditionalSignalMessage(data)

	mon.local.Lock()
	mon.conditionals[signalMsg.Cid].receiveSignal()
	mon.local.Unlock()
}

// Synchronized is method imitating synchronized block. It's higher order function that takes Monitor and function without argument, that's executed within distributed critical section defined by Monitor
func Synchronized(mon *Monitor) func(f func()) {
	return func(run func()) {
		mon.Enter()
		run()
		mon.Exit()
	}
}
