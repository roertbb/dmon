package dmon

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Token ...
type Token struct {
	LRN  map[string]int    `json:"lrn"`
	Q    []string          `json:"q"`
	Data map[string][]byte `json:"data"`
}

func newToken(mon *Monitor) *Token {
	token := Token{
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

func (t *Token) serializeData(data *map[string]interface{}) {
	for k, v := range *data {
		marshData, err := json.Marshal(v)
		if err != nil {
			fmt.Println("failed to serialize data to token")
		}
		t.Data[k] = marshData
	}
}

func (t *Token) deserializeData(data *map[string]interface{}) error {
	for key := range *data {
		val, _ := t.Data[key]
		err := json.Unmarshal(val, (*data)[key])
		if err != nil {
			return errors.New("failed to deserialize data form Token")
		}
	}
	return nil
}

func (t *Token) updateQ(mon *Monitor) {
	for address := range t.LRN {
		if stringIndex(t.Q, address) == -1 && mon.RN[address] == t.LRN[address]+1 {
			t.Q = append(t.Q, address)
		}
	}
}

func (t *Token) pop() (string, error) {
	if len(t.Q) > 0 {
		address := t.Q[0]
		t.Q = t.Q[1:]
		return address, nil
	}

	return "", errors.New("empty token queue")
}
