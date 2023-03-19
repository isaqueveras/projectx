package core

import "io"

type (
	Encoder[T any] interface {
		Encode(io.Writer, T) error
	}

	Decoder[T any] interface {
		Decode(io.Reader, T) error
	}
)
