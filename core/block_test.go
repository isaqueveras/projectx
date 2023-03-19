package core

import (
	"testing"
	"time"

	"github.com/isaqueveras/projectx/crypto"
	"github.com/isaqueveras/projectx/types"
	"github.com/stretchr/testify/assert"
)

func randomBlock(height uint32) *Block {
	var (
		header = &Header{
			Version:       1,
			PrevBlockHash: types.RandomHash(),
			Height:        height,
			Timestamp:     time.Now().UnixNano(),
		}
		tx = Transaction{
			Data: []byte("foo"),
		}
	)

	return NewBlock(header, []Transaction{tx})
}

func TestHashBlock(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(0)
	assert.Nil(t, b.Sign(privKey))
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(0)
	assert.Nil(t, b.Sign(privKey))
	assert.NotNil(t, b.Signature)

	otherPrivKey := crypto.GeneratePrivateKey()
	b.Validator = otherPrivKey.PublicKey()
	assert.NotNil(t, b.Verify())

	b.Height = 100
	assert.Nil(t, b.Verify())
}
