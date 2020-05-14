package dmon

import (
	"encoding/json"
	"errors"
	"fmt"
)

type token struct {
	LRN  map[string]int    `json:"lrn"`
	Q    []string          `json:"q"`
	Data map[string][]byte `json:"data"`
}

func newToken(mon *Monitor) *token {
	token := token{
		LRN:  map[string]int{},
		Q:    []string{},
		Data: map[string][]byte{},
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

func (t *token) deserializeData(data *map[string]interface{}) error {
	for key := range *data {
		val, _ := t.Data[key]
		err := json.Unmarshal(val, (*data)[key])
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

func (t *token) pop() (string, error) {
	if len(t.Q) > 0 {
		address := t.Q[0]
		t.Q = t.Q[1:]
		return address, nil
	}

	return "", errors.New("empty token queue")
}
