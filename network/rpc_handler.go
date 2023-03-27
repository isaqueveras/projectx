package network

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/isaqueveras/projectx/core"
	"github.com/sirupsen/logrus"
)

type DecodedMessage struct {
	From NetAddr
	Data any
}

type RPCDecodeFunc func(RPC) (*DecodedMessage, error)

type RPCProcessor interface {
	ProcessMessage(*DecodedMessage) error
}

func DefaultRPCDecodeFunc(rpc RPC) (_ *DecodedMessage, err error) {
	msg := Message{}
	if err = gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
		return nil, fmt.Errorf("failed to decode message from %s: %s", rpc.From, err)
	}

	logrus.WithFields(logrus.Fields{
		"from": rpc.From,
		"type": msg.Header,
	}).Debug("new incoming message")

	switch msg.Header {
	case MessageTypeTx:
		tx := new(core.Transaction)
		if err = tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: tx,
		}, nil
	default:
		return nil, fmt.Errorf("invalid message header %x", msg.Header)
	}
}
