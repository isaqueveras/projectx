package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	trA := NewLocalTransport("A")
	trB := NewLocalTransport("B")

	trA.Connect(trB)
	trB.Connect(trA)

	assert.Equal(t, trA.peers[trB.Addr()], trB)
	assert.Equal(t, trB.peers[trA.Addr()], trA)
}

func TestSendMessage(t *testing.T) {
	trA := NewLocalTransport("A")
	trB := NewLocalTransport("B")

	trA.Connect(trB)
	trB.Connect(trA)

	msg := []byte("Ayrton Senna")
	assert.Nil(t, trA.SendMessage(trB.addr, msg))

	rpc := <-trB.Consume()
	buf := make([]byte, len(msg))
	n, err := rpc.Payload.Read(buf)
	assert.Nil(t, err)
	assert.Equal(t, n, len(msg))

	assert.Equal(t, buf, msg)
	assert.Equal(t, rpc.From, trA.addr)
}
