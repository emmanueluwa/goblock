package core

import (
	"fmt"
	"testing"

	"github.com/emmanueluwa/goblock/types"
	"github.com/stretchr/testify/assert"
)

func TestAddBlock(test *testing.T) {
	blockchain := newBlockchainWithGenesis(test)

	lengthBlocks := 1000
	for i := 0; i < lengthBlocks; i++ {
		block := randomBlockWithSignature(test, uint32(i+1), getPrevBlockHash(test, blockchain, uint32(i+1)))
		assert.Nil(test, blockchain.AddBlock(block))
	}

	assert.Equal(test, blockchain.Height(), uint32(lengthBlocks))
	assert.Equal(test, len(blockchain.headers), lengthBlocks+1)
	assert.NotNil(test, blockchain.AddBlock(randomBlock(101, types.Hash{})))
}

func TestNewBlockchain(test *testing.T) {
	blockchain := newBlockchainWithGenesis(test)
	//not yet validated
	assert.NotNil(test, blockchain.validator)
	assert.Equal(test, blockchain.Height(), uint32(0))

	fmt.Print(blockchain.Height())
}

func TestHasBlock(test *testing.T) {
	blockchain := newBlockchainWithGenesis(test)
	assert.True(test, blockchain.HasBlock(0))
	assert.False(test, blockchain.HasBlock(1))
	assert.False(test, blockchain.HasBlock(101))
}

func TestGetHeader(test *testing.T) {
	blockchain := newBlockchainWithGenesis(test)

	lengthBlocks := 1000
	for i := 0; i < lengthBlocks; i++ {
		block := randomBlockWithSignature(test, uint32(i+1), getPrevBlockHash(test, blockchain, uint32(i+1)))
		assert.Nil(test, blockchain.AddBlock(block))
		header, err := blockchain.GetHeader(block.Height)
		assert.Nil(test, err)
		assert.Equal(test, header, block.Header)
	}

}

func TestAddingBlockTooHigh(test *testing.T) {
	blockchain := newBlockchainWithGenesis(test)

	assert.Nil(test, blockchain.AddBlock(randomBlockWithSignature(test, 1, getPrevBlockHash(test, blockchain, uint32(1)))))
	assert.NotNil(test, blockchain.AddBlock(randomBlockWithSignature(test, 101, types.Hash{})))
}

// helper function
func newBlockchainWithGenesis(test *testing.T) *Blockchain {
	blockchain, err := NewBlockchain(randomBlock(0, types.Hash{}))
	assert.Nil(test, err)

	return blockchain
}

func getPrevBlockHash(test *testing.T, blockchain *Blockchain, height uint32) types.Hash {
	prevHeader, err := blockchain.GetHeader(height - 1)
	assert.Nil(test, err)

	return BlockHasher{}.Hash(prevHeader)
}
