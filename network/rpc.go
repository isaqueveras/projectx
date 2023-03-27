package network

import (
	"bytes"
	"encoding/gob"
	"io"
)

type MessageType byte

const (
	MessageTypeTx MessageType = 0x1
	MessageTypeBlock
)

type RPC struct {
	From    NetAddr
	Payload io.Reader
}

type Message struct {
	Header MessageType
	Data   []byte
}

func NewMessage(t MessageType, data []byte) *Message {
	return &Message{Header: t, Data: data}
}

func (msg *Message) Bytes() []byte {
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(msg)
	return buf.Bytes()
}
