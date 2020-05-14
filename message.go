package dmon

import (
	"encoding/json"
	"errors"
	"fmt"
)

const requestCSMessageType = "requestCS"
const tokenMessageType = "token"

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

func serializeTokenMessage(mid string, t *Token) ([]byte, error) {
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

func deserializeTokenMessage(data []byte) (*Token, error) {
	var token Token
	err := json.Unmarshal(data, &token)
	if err != nil {
		return nil, errors.New("failed to deserialize Token")
	}
	return &token, nil
}
