package core

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/isaqueveras/projectx/crypto"
	"github.com/isaqueveras/projectx/types"
)

type Header struct {
	Version       uint32
	DataHash      types.Hash
	PrevBlockHash types.Hash
	Timestamp     int64
	Height        uint32
	Nonce         uint64
}

func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(h)
	return buf.Bytes()
}

type Block struct {
	*Header

	Transaction []Transaction
	Validator   crypto.PublicKey
	Signature   *crypto.Signature

	// Cached version of the header hash
	hash types.Hash
}

func NewBlock(h *Header, tx []Transaction) *Block {
	return &Block{Header: h, Transaction: tx}
}

func (b *Block) AddTransaction(tx *Transaction) {
	b.Transaction = append(b.Transaction, *tx)
}

func (b *Block) Sign(privKey crypto.PrivateKey) error {
	sig, err := privKey.Sign(b.Header.Bytes())
	if err != nil {
		return err
	}

	b.Validator = privKey.PublicKey()
	b.Signature = sig
	return nil
}

func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("block has no signature")
	}

	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return fmt.Errorf("block has invalid signature")
	}

	for idx := range b.Transaction {
		if err := b.Transaction[idx].Verify(); err != nil {
			return err
		}
	}

	return nil
}

func (b *Block) Decode(dec Decoder[*Block]) error {
	return dec.Decode(b)
}

func (b *Block) Encode(enc Encoder[*Block]) error {
	return enc.Encode(b)
}

func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b.Header)
	}
	return b.hash
}
