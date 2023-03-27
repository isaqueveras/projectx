package network

import (
	"bytes"
	"fmt"
	"sync"
)

type LocalTransport struct {
	addr    NetAddr
	consume chan RPC
	peers   map[NetAddr]*LocalTransport

	sync.RWMutex
}

func NewLocalTransport(addr NetAddr) *LocalTransport {
	return &LocalTransport{
		addr:    addr,
		consume: make(chan RPC, 1024),
		peers:   make(map[NetAddr]*LocalTransport),
	}
}

func (l *LocalTransport) Consume() <-chan RPC {
	return l.consume
}

func (l *LocalTransport) Connect(tr Transport) error {
	l.Lock()
	defer l.Unlock()
	l.peers[tr.Addr()] = tr.(*LocalTransport)
	return nil
}

func (l *LocalTransport) SendMessage(to NetAddr, payload []byte) error {
	l.RLock()
	defer l.RUnlock()

	peer, ok := l.peers[to]
	if !ok {
		return fmt.Errorf("%s: could not send message to %s", l.addr, to)
	}

	peer.consume <- RPC{
		From:    l.addr,
		Payload: bytes.NewReader(payload),
	}

	return nil
}

func (l *LocalTransport) Addr() NetAddr {
	return l.addr
}
