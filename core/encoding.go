package core

import (
	"crypto/elliptic"
	"encoding/gob"
	"io"
)

type Encoder[T any] interface {
	Encode(T) error
}

type Decoder[T any] interface {
	Decode(T) error
}

type GobTxEncoder struct {
	writer io.Writer
}

type GobTxDecoder struct {
	reader io.Reader
}

func NewGobTxEncoder(writer io.Writer) *GobTxEncoder {
	gob.Register(elliptic.P256())
	return &GobTxEncoder{
		writer: writer,
	}
}

func (encoder *GobTxEncoder) Encode(transaction *Transaction) error {
	return gob.NewEncoder(encoder.writer).Encode(transaction)
}

func NewGobTxDecoder(reader io.Reader) *GobTxDecoder {
	gob.Register(elliptic.P256())
	return &GobTxDecoder{
		reader: reader,
	}
}

func (decoder *GobTxDecoder) Decode(transaction *Transaction) error {
	return gob.NewDecoder(decoder.reader).Decode(transaction)
}
