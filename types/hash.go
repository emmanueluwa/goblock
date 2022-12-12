package types

import (
	"crypto/rand"
	"fmt"
)

type Hash [32]uint8

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
