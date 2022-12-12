package core

import (
	"encoding/binary"
	"io"

	"github.com/emmanueluwa/goblock/types"
)

type Header struct {
	Version       uint32
	PreviousBlock types.Hash
	TimeStamp     int64
	//eg 3 blocks, 1 genesis height=2
	Height uint32
	Nonce  uint64
}

// encode header to byte slice
func (header *Header) EncodeBinary(writer io.Writer) error {
	if err := binary.Write(writer, binary.LittleEndian, &header.Version); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, &header.PreviousBlock); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, &header.TimeStamp); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, &header.Height); err != nil {
		return err
	}
	return binary.Write(writer, binary.LittleEndian, &header.Nonce)
}

// received block byte slice, decoded to header
func (header *Header) DecodeBinary(reader io.Reader) error {
	if err := binary.Read(reader, binary.LittleEndian, &header.Version); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &header.PreviousBlock); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &header.TimeStamp); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &header.Height); err != nil {
		return err
	}
	return binary.Read(reader, binary.LittleEndian, &header.Nonce)
}

type Block struct {
	Header
	Transactions []Transaction
}

func (block *Block) EncodeBinary(writer io.Writer) error {
	if err := block.Header.EncodeBinary(writer); err != nil {
		return err
	}

	for _, tx := range block.Transactions {
		if err := tx.EncodeBinary(writer); err != nil {
			return err
		}
	}

	return nil
}

func (block *Block) DecodeBinary(reader io.Reader) error {
	if err := block.Header.DecodeBinary(reader); err != nil {
		return err
	}

	for _, tx := range block.Transactions {
		if err := tx.DecodeBinary(reader); err != nil {
			return err
		}
	}

	return nil
}
