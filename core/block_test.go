package core

import (
	"fmt"
	"testing"
	"time"

	"github.com/emmanueluwa/goblock/types"
)

func randomBlock(height uint32) *Block {
	header := &Header{
		Version:           1,
		PreviousBlockHash: types.RandomHash(),
		TimeStamp:         time.Now().UnixNano(),
	}
	transaction := Transaction{
		Data: []byte("block"),
	}
	return NewBlock(header, []Transaction{transaction})
}

func TestHashBlock(test *testing.T) {
	block := randomBlock(0)
	fmt.Println(block.Hash(BlockHasher{}))
}
