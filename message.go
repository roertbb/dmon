package dmon

import (
	"encoding/json"
	"errors"
	"fmt"
)

const requestCSMessageType = "requestCS"
const tokenMessageType = "token"
const conditionalWaitMessageType = "conditionalWait"
const conditionalSignalMessageType = "conditionalSignal"

type message struct {
	Type string `json:"type"`
	Mid  string `json:"mid"`
	Data []byte `json:"data"`
}

func deserializeMessage(data []byte) (*message, error) {
	var msg message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("failed to deserialize message")
	}
	return &msg, nil
}

type requestCSMessage struct {
	From string `json:"from"`
	SN   int    `json:"sn"`
}

func serializeRequestCSMessage(from string, mid string, SN int) ([]byte, error) {
	request := requestCSMessage{
		From: from,
		SN:   SN,
	}
	marshRequest, err := json.Marshal(request)
	if err != nil {
		return nil, errors.New("failed to serialize RequestCSMessage")
	}

	message := message{
		Type: requestCSMessageType,
		Mid:  mid,
		Data: marshRequest,
	}
	marshMessage, err := json.Marshal(message)
	if err != nil {
		return nil, errors.New("failed to serialize Message for RequestCSMessage")
	}

	return marshMessage, nil
}

func deserializeRequestCSMessage(data []byte) (*requestCSMessage, error) {
	var req requestCSMessage
	err := json.Unmarshal(data, &req)
	if err != nil {
		return nil, errors.New("failed to deserialize RequestCSMessage")
	}
	return &req, nil
}

func serializeTokenMessage(mid string, t *token) ([]byte, error) {
	marshToken, err := json.Marshal(t)
	if err != nil {
		return nil, errors.New("failed to serialize token")
	}

	message := message{
		Type: tokenMessageType,
		Mid:  mid,
		Data: marshToken,
	}
	marshMessage, err := json.Marshal(message)
	if err != nil {
		return nil, errors.New("failed to serialize Message for Token")
	}

	return marshMessage, nil
}

func deserializeTokenMessage(data []byte) (*token, error) {
	var token token
	err := json.Unmarshal(data, &token)
	if err != nil {
		return nil, errors.New("failed to deserialize Token")
	}
	return &token, nil
}

type conditionalWaitMessage struct {
	From string `json:"from"`
	Cid  string `json:"cond"`
}

func serializeConditionalWaitMessage(from string, mid string, cid string) ([]byte, error) {
	condWait := conditionalWaitMessage{
		From: from,
		Cid:  cid,
	}
	marshCondWait, err := json.Marshal(condWait)
	if err != nil {
		return nil, errors.New("failed to serialize ConditionalWaitMessage")
	}
	msg := message{
		Type: conditionalWaitMessageType,
		Mid:  mid,
		Data: marshCondWait,
	}
	marshMsg, err := json.Marshal(msg)
	if err != nil {
		return nil, errors.New("failed to serialize Message for ConditionalWaitMessage")
	}
	return marshMsg, nil
}

func deserializeConditionalWaitMessage(data []byte) (*conditionalWaitMessage, error) {
	var condWait conditionalWaitMessage
	err := json.Unmarshal(data, &condWait)
	if err != nil {
		return nil, errors.New("failed to deserialize RequestCSMessage")
	}
	return &condWait, nil
}

type conditionalSignalMessage struct {
	Cid string `json:"cond"`
}

func serializeConditionalSignalMessage(mid string, cid string) ([]byte, error) {
	condSignal := conditionalSignalMessage{
		Cid: cid,
	}
	marshCondSignal, err := json.Marshal(condSignal)
	if err != nil {
		return nil, errors.New("failed to serialize ConditionalWaitMessage")
	}
	msg := message{
		Type: conditionalSignalMessageType,
		Mid:  mid,
		Data: marshCondSignal,
	}
	marshMsg, err := json.Marshal(msg)
	if err != nil {
		return nil, errors.New("failed to serialize Message for ConditionalWaitMessage")
	}
	return marshMsg, nil
}

func deserializeConditionalSignalMessage(data []byte) (*conditionalSignalMessage, error) {
	var condSignal conditionalSignalMessage
	err := json.Unmarshal(data, &condSignal)
	if err != nil {
		return nil, errors.New("failed to deserialize ConditionalSignalMessage")
	}
	return &condSignal, nil
}
