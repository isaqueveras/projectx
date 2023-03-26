package core

import (
	"fmt"

	"github.com/isaqueveras/projectx/crypto"
	"github.com/isaqueveras/projectx/types"
)

type Transaction struct {
	Data      []byte
	From      crypto.PublicKey
	Signature *crypto.Signature

	hash      types.Hash // cached version of the tx data hash
	firstSeen int64      // firstSeen is the timestamp of when this tx is first seen locally
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{Data: data}
}

func (tx *Transaction) Hash(hasher Hasher[*Transaction]) types.Hash {
	if tx.hash.IsZero() {
		tx.hash = hasher.Hash(tx)
	}
	return tx.hash
}

func (tx *Transaction) Sign(privKey crypto.PrivateKey) error {
	sig, err := privKey.Sign(tx.Data)
	if err != nil {
		return err
	}

	tx.From = privKey.PublicKey()
	tx.Signature = sig
	return nil
}

func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}

	if !tx.Signature.Verify(tx.From, tx.Data) {
		return fmt.Errorf("invalid transaction signature")
	}

	return nil
}

func (tx *Transaction) Encode(dec Encoder[*Transaction]) (err error) {
	return dec.Encode(tx)
}

func (tx *Transaction) Decode(dec Decoder[*Transaction]) (err error) {
	return dec.Decode(tx)
}

func (tx *Transaction) SetFirstSeen(t int64) {
	tx.firstSeen = t
}

func (tx *Transaction) FirstSeen() int64 {
	return tx.firstSeen
}
