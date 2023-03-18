package core

import "io"

type Transaction struct {
	Data []byte
}

func (tx *Transaction) EncodeBinary(w io.Writer) (err error) {
	return
}

func (tx *Transaction) DecodeBinary(r io.Reader) (err error) {
	return
}
