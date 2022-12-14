package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// helper function
func newBlockchainWithGenesis(test *testing.T) *Blockchain {
	blockchain, err := NewBlockchain(randomBlock(0))
	assert.Nil(test, err)

	return blockchain
}

func TestAddBlock(test *testing.T) {
	blockchain := newBlockchainWithGenesis(test)

	lengthBlocks := 1000
	for i := 0; i < lengthBlocks; i++ {
		block := randomBlockWithSignature(test, uint32(i+1))
		assert.Nil(test, blockchain.AddBlock(block))
	}

	assert.Equal(test, blockchain.Height(), uint32(lengthBlocks))
	assert.Equal(test, len(blockchain.headers), lengthBlocks+1)
	assert.NotNil(test, blockchain.AddBlock(randomBlock(101)))
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
}
