package dmon

// Conditional ...
type Conditional struct {
	cid        string
	monitor    *Monitor
	signalChan chan bool
	waiting    []string
}

func newConditional(mon *Monitor, cid string) *Conditional {
	cond := Conditional{
		cid:        cid,
		monitor:    mon,
		signalChan: make(chan bool),
		waiting:    []string{},
	}

	return &cond
}

// Wait ...
func (cond *Conditional) Wait() {
	cond.waiting = append(cond.waiting, cond.monitor.env.address)
	waitMsg, _ := serializeConditionalWaitMessage(cond.monitor.env.address, cond.monitor.mid, cond.cid)
	cond.monitor.env.broadcast(waitMsg)

	cond.monitor.Exit()
	<-cond.signalChan
	cond.monitor.Enter()
}

// Notify ...
func (cond *Conditional) Notify() {
	if len(cond.waiting) > 0 {
		waitingAddress := cond.waiting[0]
		cond.waiting = cond.waiting[1:]
		if waitingAddress == cond.monitor.env.address {
			cond.signalChan <- true
		} else {
			signalMsg, _ := serializeConditionalSignalMessage(cond.monitor.mid, cond.cid)
			cond.waiting = removeStringFromSlice(cond.waiting, waitingAddress)
			cond.monitor.env.send(waitingAddress, signalMsg)
		}
	}
}

// NotifyAll ...
func (cond *Conditional) NotifyAll() {
	signalMsg, _ := serializeConditionalSignalMessage(cond.monitor.mid, cond.cid)
	for _, addr := range cond.waiting {
		if addr == cond.monitor.env.address {
			cond.signalChan <- true
		} else {
			cond.monitor.env.send(addr, signalMsg)
		}
	}
	cond.waiting = []string{}
}

func (cond *Conditional) receiveSignal() {
	if stringIndex(cond.waiting, cond.monitor.env.address) != -1 {
		cond.waiting = removeStringFromSlice(cond.waiting, cond.monitor.env.address)
		cond.signalChan <- true
	}
	// if I'm not waiting, but received signal, should have to send it to someone else?
}

func (cond *Conditional) addToWaiting(address string) {
	cond.waiting = append(cond.waiting, address)
}
