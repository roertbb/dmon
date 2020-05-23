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

// Wait waits on Conditional variables
func (cond *Conditional) Wait() {
	cond.monitor.local.Lock()
	if stringIndex(cond.waiting, cond.monitor.env.address) == -1 {
		cond.waiting = append(cond.waiting, cond.monitor.env.address)
	}
	cond.monitor.local.Unlock()

	cond.monitor.Exit()
	<-cond.signalChan
	cond.monitor.Enter()

	cond.monitor.local.Lock()
	cond.waiting = removeStringFromSlice(cond.waiting, cond.monitor.env.address)
	cond.monitor.local.Unlock()
}

// Notify sends signal message to one of the processes waiting on Conditional variable
func (cond *Conditional) Notify() {
	cond.monitor.local.Lock()

	if len(cond.waiting) > 0 {
		waitingAddress := cond.waiting[0]
		if waitingAddress == cond.monitor.env.address {
			cond.signalChan <- true
		} else {
			signalMsg, _ := serializeConditionalSignalMessage(cond.monitor.mid, cond.cid)
			cond.monitor.env.send(waitingAddress, signalMsg)
		}
	}

	cond.monitor.local.Unlock()
}

// NotifyAll sends signal message to all of the processes waiting on Conditional variable
func (cond *Conditional) NotifyAll() {
	cond.monitor.local.Lock()

	for _, addr := range cond.waiting {
		if addr == cond.monitor.env.address {
			cond.signalChan <- true
		} else {
			signalMsg, _ := serializeConditionalSignalMessage(cond.monitor.mid, cond.cid)
			cond.monitor.env.send(addr, signalMsg)
		}
	}

	cond.monitor.local.Unlock()
}

func (cond *Conditional) receiveSignal() {
	if stringIndex(cond.waiting, cond.monitor.env.address) != -1 {
		cond.signalChan <- true
	}
}
