package core

import (
	"crypto/sha256"

	"github.com/emmanueluwa/goblock/types"
)

// hash the type and return hash, allows the use of different hash methods
type Hasher[T any] interface {
	Hash(T) types.Hash
}

type BlockHasher struct {
}

func (BlockHasher) Hash(block *Block) types.Hash {
	header := sha256.Sum256(block.HeaderData())
	return types.Hash(header)
}
