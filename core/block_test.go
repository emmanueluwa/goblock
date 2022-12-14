package core

import (
	"testing"
	"time"

	"github.com/emmanueluwa/goblock/crypto"
	"github.com/emmanueluwa/goblock/types"
	"github.com/stretchr/testify/assert"
)

func randomBlock(height uint32) *Block {
	header := &Header{
		Version:           1,
		PreviousBlockHash: types.RandomHash(),
		Height:            height,
		TimeStamp:         time.Now().UnixNano(),
	}
	transaction := Transaction{
		Data: []byte("block"),
	}
	return NewBlock(header, []Transaction{transaction})
}

func randomBlockWithSignature(test *testing.T, height uint32) *Block {
	privKey := crypto.GeneratePrivateKey()
	block := randomBlock(height)
	assert.Nil(test, block.Sign(privKey))

	return block
}

func TestSignBlock(test *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	block := randomBlock(0)

	assert.Nil(test, block.Sign(privKey))
	assert.NotNil(test, block.Signature)
}

// verify validated pubKey signed block in questions
func TestVerifyBlock(test *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	block := randomBlock(0)

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
