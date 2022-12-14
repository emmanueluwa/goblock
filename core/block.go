package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
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

// signature is embedded in block
func (block *Block) Sign(privKey crypto.PrivateKey) error {
	signature, err := privKey.Sign(block.HeaderData())
	if err != nil {
		return err
	}

	block.Validator = privKey.PublicKey()
	block.Signature = signature

	return nil
}

func (block *Block) Verify() error {
	if block.Signature == nil {
		return fmt.Errorf("block has no signature")
	}

	if !block.Signature.Verify(block.Validator, block.HeaderData()) {
		return fmt.Errorf("block has invalid signature")
	}

	return nil
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

// bytes of data that needs to be signed and hashed
func (block *Block) HeaderData() []byte {
	buffer := &bytes.Buffer{}
	encode := gob.NewEncoder(buffer)
	encode.Encode(block.Header)

	return buffer.Bytes()
}
