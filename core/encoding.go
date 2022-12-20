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
	return &GobTxEncoder{
		writer: writer,
	}
}

func (encoder *GobTxEncoder) Encode(transaction *Transaction) error {
	return gob.NewEncoder(encoder.writer).Encode(transaction)
}

func NewGobTxDecoder(reader io.Reader) *GobTxDecoder {
	return &GobTxDecoder{
		reader: reader,
	}
}

func (decoder *GobTxDecoder) Decode(transaction *Transaction) error {
	return gob.NewDecoder(decoder.reader).Decode(transaction)
}

type GobBlockEncoder struct {
	writer io.Writer
}

func NewGobBlockEncoder(writer io.Writer) *GobBlockEncoder {
	return &GobBlockEncoder{
		writer: writer,
	}
}

func (encoder *GobBlockEncoder) Encode(block *Block) error {
	return gob.NewEncoder(encoder.writer).Encode(block)
}

type GobBlockDecoder struct {
	reader io.Reader
}

func NewGobBlockDecoder(reader io.Reader) *GobBlockDecoder {
	return &GobBlockDecoder{
		reader: reader,
	}
}

func (decoder *GobBlockDecoder) Decode(block *Block) error {
	return gob.NewDecoder(decoder.reader).Decode(block)
}

func init() {
	gob.Register(elliptic.P256())
}
