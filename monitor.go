package dmon

import (
	"strconv"
	"sync"
)

// Monitor is distributed monitor structure handling data, token exchange and it's Conditional variables
type Monitor struct {
	mid    string
	env    *Env
	local  *sync.Mutex
	isInCS bool
	// csMutex *sync.Mutex
	// conditionals map[string]*Conditional
	data      *map[string]interface{}
	RN        map[string]int
	token     *Token
	tokenChan chan bool
}

func newMonitor(mid string, env *Env) (*Monitor, error) {
	monitor := Monitor{
		mid:       mid,
		env:       env,
		local:     &sync.Mutex{},
		data:      &map[string]interface{}{},
		RN:        map[string]int{},
		isInCS:    false,
		token:     nil,
		tokenChan: make(chan bool),
	}

	monitor.RN[env.address] = 0
	for address := range env.sockets {
		monitor.RN[address] = 0
	}

	lowestAddress := env.getLowestAddress()
	if lowestAddress == env.address {
		monitor.token = newToken(&monitor)
	}

	return &monitor, nil
}

// RegisterSharedData ...
func (mon *Monitor) RegisterSharedData(data ...interface{}) {
	for id, value := range data {
		(*mon.data)[strconv.Itoa(id)] = value
	}
}

// Enter ...
func (mon *Monitor) Enter() {
	// mon.csMutex.Lock()
	mon.local.Lock()

	if mon.token == nil {
		mon.RN[mon.env.address]++
		requestMsg, _ := serializeRequestCSMessage(mon.env.address, mon.mid, mon.RN[mon.env.address])
		mon.env.broadcast(requestMsg)
		mon.local.Unlock()

		<-mon.tokenChan

		mon.local.Lock()
		mon.token.deserializeData(mon.data)
	}

	mon.isInCS = true
	mon.local.Unlock()
}

// Exit ...
func (mon *Monitor) Exit() {
	mon.local.Lock()

	mon.token.LRN[mon.env.address] = mon.RN[mon.env.address]
	mon.token.updateQ(mon)

	address, _ := mon.token.pop()
	if address != "" {
		mon.sendToken(address)
	}

	mon.isInCS = false

	mon.local.Unlock()
	// mon.csMutex.Unlock()
}

func (mon *Monitor) sendToken(address string) {
	mon.token.serializeData(mon.data)
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
	if mon.token != nil && !mon.isInCS && mon.RN[requestCS.From] == mon.token.LRN[requestCS.From]+1 {
		mon.sendToken(requestCS.From)
	}
	mon.local.Unlock()
}

func (mon *Monitor) handleTokenMessage(data []byte) {
	token, _ := deserializeTokenMessage(data)

	mon.local.Lock()
	mon.token = token
	mon.token.deserializeData(mon.data)
	mon.local.Unlock()

	mon.tokenChan <- true
}
