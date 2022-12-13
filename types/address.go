package types

import (
	"encoding/hex"
	"fmt"
)

type Address [20]uint8

func (address Address) ToSlice() []byte {
	bytes := make([]byte, 20)
	for i := 0; i < 20; i++ {
		bytes[i] = address[i]
	}
	return bytes
}

func (address Address) String() string {
	return hex.EncodeToString(address.ToSlice())
}

func AddressFromBytes(bytes []byte) Address {
	if len(bytes) != 20 {
		message := fmt.Sprintf("bytes with length %d should be 20", len(bytes))
		panic(message)
	}

	var value [20]uint8
	for i := 0; i < 20; i++ {
		value[i] = bytes[i]
	}

	return Address(value)
}
