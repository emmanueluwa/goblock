package core

import "io"

type Transaction struct {
	Data []byte
}

func (tx *Transaction) DecodeBinary(reader io.Reader) error { return nil }

func (tx *Transaction) EncodeBinary(writer io.Writer) error { return nil }
