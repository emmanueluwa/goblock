package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"

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

// allowing us to check hash of previous block/
func (header *Header) Bytes() []byte {
	buffer := &bytes.Buffer{}
	encode := gob.NewEncoder(buffer)
	//returning bytes of the header
	encode.Encode(header)

	return buffer.Bytes()
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

func NewBlock(header *Header, transactions []Transaction) (*Block, error) {
	return &Block{
		Header:       header,
		Transactions: transactions,
	}, nil
}

func NewBlockFromPrevHeader(prevHeader *Header, transactions []Transaction) (*Block, error) {
	dataHash, err := CalculateDataHash(transactions)
	if err != nil {
		return nil, err
	}

	header := &Header{
		Version:           1,
		DataHash:          dataHash,
		PreviousBlockHash: BlockHasher{}.Hash(prevHeader),
		TimeStamp:         time.Now().UnixNano(),
		Height:            prevHeader.Height + 1,
	}

	return NewBlock(header, transactions)
}

func (block *Block) AddTransaction(transaction *Transaction) {
	block.Transactions = append(block.Transactions, *transaction)
}

// signature is embedded in block
func (block *Block) Sign(privKey crypto.PrivateKey) error {
	signature, err := privKey.Sign(block.Header.Bytes())
	if err != nil {
		return err
	}

	block.Validator = privKey.PublicKey()
	block.Signature = signature

	return nil
}

func (block *Block) Verify() error {
	//verifying block
	if block.Signature == nil {
		return fmt.Errorf("block has no signature")
	}

	if !block.Signature.Verify(block.Validator, block.Header.Bytes()) {
		return fmt.Errorf("block has invalid signature")
	}

	//verifying transaction
	for _, transaction := range block.Transactions {
		if err := transaction.Verify(); err != nil {
			return err
		}
	}

	//verifying datahash
	dataHash, err := CalculateDataHash(block.Transactions)
	if err != nil {
		return err
	}
	if dataHash != block.DataHash {
		return fmt.Errorf("block (%s) has an invalid data hash", block.Hash(BlockHasher{}))
	}

	return nil
}

func (block *Block) Decode(decoder Decoder[*Block]) error {
	return decoder.Decode(block)
}

func (block *Block) Encode(encoder Encoder[*Block]) error {
	return encoder.Encode(block)
}

func (block *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if block.hash.IsZero() {
		block.hash = hasher.Hash(block.Header)
	}
	return block.hash
}

func CalculateDataHash(transactions []Transaction) (hash types.Hash, err error) {
	buffer := &bytes.Buffer{}

	for _, transaction := range transactions {
		if err := transaction.Encode(NewGobTxEncoder(buffer)); err != nil {
			return types.Hash{}, err
		}
	}

	hash = sha256.Sum256(buffer.Bytes())

	return
}
