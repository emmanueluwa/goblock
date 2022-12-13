package core

import (
	"io"

	"github.com/emmanueluwa/goblock/crypto"
	"github.com/emmanueluwa/goblock/types"
)

type Header struct {
	Version           uint32
	DataHash          types.Hash
	PreviousBlockHash types.Hash
	TimeStamp         int64
	//eg 3 blocks, 1 genesis height=2
	Height uint32
}

type Block struct {
	*Header
	Transactions []Transaction

	//validator to enable a block to propose blocks into the network
	Validator crypto.PublicKey
	Signature *crypto.Signature

	//chached version, costly to repeat it each time its needed
	hash types.Hash
}

func NewBlock(header *Header, transaction []Transaction) *Block {
	return &Block{
		Header:       header,
		Transactions: transaction,
	}
}

func (block *Block) Decode(reader io.Reader, decoder Decoder[*Block]) error {
	return decoder.Decode(reader, block)
}

func (block *Block) Encode(writer io.Writer, encoder Encoder[*Block]) error {
	return encoder.Encode(writer, block)
}

func (block *Block) Hash(hasher Hasher[*Block]) types.Hash {
	if block.hash.IsZero() {
		block.hash = hasher.Hash(block)
	}
	return block.hash
}
