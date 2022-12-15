package core

import (
	"testing"
	"time"

	"github.com/emmanueluwa/goblock/crypto"
	"github.com/emmanueluwa/goblock/types"
	"github.com/stretchr/testify/assert"
)

func TestSignBlock(test *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	block := randomBlock(0, types.Hash{})

	assert.Nil(test, block.Sign(privKey))
	assert.NotNil(test, block.Signature)
}

// verify validated pubKey signed block in questions
func TestVerifyBlock(test *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	block := randomBlock(0, types.Hash{})

	assert.Nil(test, block.Sign(privKey))
	assert.Nil(test, block.Verify())

	//testing against random pubKey
	privKeyHack := crypto.GeneratePrivateKey()
	randomPubKey := privKeyHack.PublicKey()
	block.Validator = randomPubKey
	assert.NotNil(test, block.Verify())

	//testing change in block data
	block.Height = 100
	assert.NotNil(test, block.Verify())
}

// helper functions
func randomBlock(height uint32, prevBlockHash types.Hash) *Block {
	header := &Header{
		Version:           1,
		PreviousBlockHash: prevBlockHash,
		Height:            height,
		TimeStamp:         time.Now().UnixNano(),
	}

	return NewBlock(header, []Transaction{})
}

func randomBlockWithSignature(test *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	privKey := crypto.GeneratePrivateKey()
	block := randomBlock(height, prevBlockHash)
	transaction := randomTransactionWithSignature(test)
	block.AddTransaction(transaction)
	assert.Nil(test, block.Sign(privKey))

	return block
}
