package core

import (
	"bytes"
	"testing"
	"time"

	"github.com/emmanueluwa/goblock/types"
	"github.com/stretchr/testify/assert"
)

func TestHeader_Encode_Decode(test *testing.T) {
	header := &Header{
		Version:       1,
		PreviousBlock: types.RandomHash(),
		TimeStamp:     time.Now().UnixNano(),
		Height:        12,
		Nonce:         993894,
	}

	buffer := &bytes.Buffer{}
	assert.Nil(test, header.EncodeBinary(buffer))

	headerDecode := &Header{}
	assert.Nil(test, headerDecode.DecodeBinary(buffer))
	assert.Equal(test, header, headerDecode)
}

func TestBlock_Encode_Decode(test *testing.T) {
	block := &Block{
		Header: Header{
			Version:       1,
			PreviousBlock: types.RandomHash(),
			TimeStamp:     time.Now().UnixNano(),
			Height:        12,
			Nonce:         993894,
		},
		Transactions: nil,
	}

	buffer := &bytes.Buffer{}
	assert.Nil(test, block.EncodeBinary(buffer))

	blockDecode := &Block{}
	assert.Nil(test, blockDecode.DecodeBinary(buffer))
	assert.Equal(test, block, blockDecode)
}
