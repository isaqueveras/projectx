package network

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/isaqueveras/projectx/core"
)

type DefaultRPCHandler struct {
	p RPCProcessor
}

type (
	RPCHandler interface {
		HandleRPC(rpc RPC) error
	}

	RPCProcessor interface {
		ProcessTransaction(NetAddr, *core.Transaction) error
	}
)

func NewDefaultRPCHandler(p RPCProcessor) *DefaultRPCHandler {
	return &DefaultRPCHandler{p: p}
}

func (h *DefaultRPCHandler) HandleRPC(rpc RPC) (err error) {
	msg := Message{}
	if err = gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
		return fmt.Errorf("failed to decode message from %s: %s", rpc.From, err)
	}

	switch msg.Header {
	case MessageTypeTx:
		tx := new(core.Transaction)
		if err = tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
			return err
		}
		return h.p.ProcessTransaction(rpc.From, tx)
	default:
		return fmt.Errorf("invalid message header %x", msg.Header)
	}
}
