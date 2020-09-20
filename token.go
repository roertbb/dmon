package dmon

import (
	"encoding/json"
	"errors"
	"fmt"
)

type token struct {
	LRN         map[string]int    `json:"lrn"`
	Q           []string          `json:"q"`
	Data        map[string][]byte `json:"data"`
	CondWaiting map[string][]byte `json:"condWaiting"`
	monitor     *Monitor
}

func newToken(mon *Monitor) *token {
	token := token{
		LRN:         map[string]int{},
		Q:           []string{},
		Data:        map[string][]byte{},
		CondWaiting: map[string][]byte{},
		monitor:     mon,
	}

	token.LRN[mon.env.address] = 0
	for address := range mon.RN {
		token.LRN[address] = 0
	}

	return &token
}

func (t *token) serializeData(data *map[string]interface{}) {
	for k, v := range *data {
		marshData, err := json.Marshal(v)
		if err != nil {
			fmt.Println("failed to serialize data to token")
		}
		t.Data[k] = marshData
	}
}

func (t *token) deserializeData(data *map[string]interface{}, mon *Monitor) error {
	for key := range *data {
		val, _ := t.Data[key]
		err := json.Unmarshal(val, (*data)[key])
		if err != nil {
			return errors.New("failed to deserialize data form Token")
		}
	}
	t.monitor = mon
	return nil
}

func (t *token) serializeCondWaiting(conds *map[string]*Conditional) {
	for k, v := range *conds {
		marshData, err := json.Marshal(v.waiting)
		if err != nil {
			fmt.Println("failed to serialize data to token")
		}
		t.CondWaiting[k] = marshData
	}
}

func (t *token) deserializeCondWaiting(conds *map[string]*Conditional) error {
	for key := range *conds {
		condWaiting, _ := t.CondWaiting[key]
		err := json.Unmarshal(condWaiting, &(*conds)[key].waiting)
		if err != nil {
			return errors.New("failed to deserialize data form Token")
		}
	}
	return nil
}

func (t *token) updateQ(mon *Monitor) {
	for address := range t.LRN {
		if stringIndex(t.Q, address) == -1 && mon.RN[address] == t.LRN[address]+1 {
			t.Q = append(t.Q, address)
		}
	}
}

func (t *token) lastSignaledOrPop() (string, error) {
	if len(t.Q) > 0 {
		var signaledAndInToken []string
		for _, addr := range t.monitor.lastSignaled {
			if stringIndex(t.Q, addr) != -1 {
				signaledAndInToken = append(signaledAndInToken, addr)
			}
		}

		var address string
		if len(signaledAndInToken) > 0 {
			address = signaledAndInToken[0]
			t.Q = removeStringFromSlice(t.Q, address)
		} else {
			address = t.Q[0]
			t.Q = t.Q[1:]
		}

		t.monitor.lastSignaled = []string{}
		return address, nil
	}

	return "", errors.New("empty token queue")
}
