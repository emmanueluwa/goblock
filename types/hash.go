package types

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type Hash [32]uint8

func (hash Hash) IsZero() bool {
	for i := 0; i < 32; i++ {
		if hash[i] != 0 {
			return false
		}
	}
	return true
}

func (hash Hash) ToSlice() []byte {
	bytes := make([]byte, 32)
	for i := 0; i < 32; i++ {
		bytes[i] = hash[i]
	}
	return bytes
}

// print string represetation
func (hash Hash) String() string {
	return hex.EncodeToString(hash.ToSlice())
}

func HashFromBytes(bytes []byte) Hash {
	//if not 32 the system cannot continue
	if len(bytes) != 32 {
		message := fmt.Sprintf("bytes with length %d should be 32", len(bytes))
		panic(message)
	}

	var value [32]uint8
	for i := 0; i < 32; i++ {
		value[i] = bytes[i]
	}

	return Hash(value)
}

func RandomBytes(size int) []byte {
	token := make([]byte, size)
	rand.Read(token)
	return token
}

func RandomHash() Hash {
	return HashFromBytes(RandomBytes(32))
}
