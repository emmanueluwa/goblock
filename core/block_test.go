package core

import (
	"bytes"
	"testing"
	"time"

	"github.com/emmanueluwa/goblock/crypto"
	"github.com/emmanueluwa/goblock/types"
	"github.com/stretchr/testify/assert"
)

func TestSignBlock(test *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	block := randomBlock(test, 0, types.Hash{})

	assert.Nil(test, block.Sign(privKey))
	assert.NotNil(test, block.Signature)
}

// verify validated pubKey signed block in questions
func TestVerifyBlock(test *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	block := randomBlock(test, 0, types.Hash{})

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

func TestDecodeEncodeBlock(test *testing.T) {
	block := randomBlock(test, 1, types.Hash{})
	buffer := &bytes.Buffer{}
	assert.Nil(test, block.Encode(NewGobBlockEncoder(buffer)))

	blockDecode := new(Block)
	assert.Nil(test, blockDecode.Decode(NewGobBlockDecoder(buffer)))
	assert.Equal(test, block, blockDecode)
}

// helper functions
func randomBlock(test *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	privKey := crypto.GeneratePrivateKey()
	transaction := randomTransactionWithSignature(test)

	header := &Header{
		Version:           1,
		PreviousBlockHash: prevBlockHash,
		Height:            height,
		TimeStamp:         time.Now().UnixNano(),
	}

	block, err := NewBlock(header, []*Transaction{transaction})
	assert.Nil(test, err)
	dataHash, err := CalculateDataHash(block.Transactions)
	assert.Nil(test, err)
	block.Header.DataHash = dataHash
	assert.Nil(test, block.Sign(privKey))

	return block
}
