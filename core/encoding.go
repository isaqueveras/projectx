package core

import (
	"crypto/elliptic"
	"encoding/gob"
	"io"
)

type (
	Encoder[T any] interface{ Encode(T) error }
	Decoder[T any] interface{ Decode(T) error }
)

type (
	GobTxEncoder struct{ w io.Writer }
	GobTxDecoder struct{ r io.Reader }
)

func NewGotTxEncoder(w io.Writer) *GobTxEncoder {
	gob.Register(elliptic.P256())
	return &GobTxEncoder{w: w}
}

func (e *GobTxEncoder) Encode(tx *Transaction) error {
	return gob.NewEncoder(e.w).Encode(tx)
}

func NewGobTxDecoder(r io.Reader) *GobTxDecoder {
	gob.Register(elliptic.P256())
	return &GobTxDecoder{r: r}
}

func (e *GobTxDecoder) Decode(tx *Transaction) error {
	return gob.NewDecoder(e.r).Decode(tx)
}
