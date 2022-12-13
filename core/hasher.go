package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"

	"github.com/emmanueluwa/goblock/types"
)

// hash the type and return hash, allows the use of different hash methods
type Hasher[T any] interface {
	Hash(T) types.Hash
}

type BlockHasher struct {
}

func (BlockHasher) Hash(block *Block) types.Hash {
	buffer := &bytes.Buffer{}
	encode := gob.NewEncoder(buffer)
	if err := encode.Encode(block.Header); err != nil {
		panic(err)
	}

	header := sha256.Sum256(buffer.Bytes())
	return types.Hash(header)
}
