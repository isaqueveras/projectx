package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"io"

	"github.com/isaqueveras/projectx/types"
)

type Header struct {
	Version   uint32
	PrevBlock types.Hash
	Timestamp int64
	Height    uint32
	Nonce     uint64
}

func (h *Header) EncodeBinary(w io.Writer) (err error) {
	if err = binary.Write(w, binary.LittleEndian, &h.Version); err != nil {
		return
	}

	if err = binary.Write(w, binary.LittleEndian, &h.PrevBlock); err != nil {
		return
	}

	if err = binary.Write(w, binary.LittleEndian, &h.Timestamp); err != nil {
		return
	}

	if err = binary.Write(w, binary.LittleEndian, &h.Height); err != nil {
		return
	}

	if err = binary.Write(w, binary.LittleEndian, &h.Nonce); err != nil {
		return
	}

	return
}

func (h *Header) DecodeBinary(r io.Reader) (err error) {
	if err = binary.Read(r, binary.LittleEndian, &h.Version); err != nil {
		return
	}

	if err = binary.Read(r, binary.LittleEndian, &h.PrevBlock); err != nil {
		return
	}

	if err = binary.Read(r, binary.LittleEndian, &h.Timestamp); err != nil {
		return
	}

	if err = binary.Read(r, binary.LittleEndian, &h.Height); err != nil {
		return
	}

	if err = binary.Read(r, binary.LittleEndian, &h.Nonce); err != nil {
		return
	}

	return
}

type Block struct {
	Header
	Transaction []Transaction

	// Cached version of the header hash
	hash types.Hash
}

func (b *Block) Hash() types.Hash {
	buf := &bytes.Buffer{}
	b.Header.EncodeBinary(buf)

	if b.hash.IsZero() {
		b.hash = types.Hash(sha256.Sum256(buf.Bytes()))
	}

	return b.hash
}

func (b *Block) EncodeBinary(w io.Writer) (err error) {
	if err = b.Header.EncodeBinary(w); err != nil {
		return
	}

	for idx := range b.Transaction {
		if err = b.Transaction[idx].EncodeBinary(w); err != nil {
			return err
		}
	}

	return
}

func (b *Block) DecodeBinary(r io.Reader) (err error) {
	if err = b.Header.DecodeBinary(r); err != nil {
		return
	}

	for idx := range b.Transaction {
		if err = b.Transaction[idx].DecodeBinary(r); err != nil {
			return err
		}
	}

	return
}
