package network

import (
	"testing"

	"github.com/emmanueluwa/goblock/core"
	"github.com/stretchr/testify/assert"
)

func TestTxPool(test *testing.T) {
	pool := NewTxPool()
	assert.Equal(test, pool.Len(), 0)
}

func TestTxPoolAddTx(test *testing.T) {
	pool := NewTxPool()
	transaction := core.NewTransaction([]byte("meow"))
	assert.Nil(test, pool.Add(transaction))
	assert.Equal(test, pool.Len(), 1)

	duplicateTransaction := core.NewTransaction([]byte("meow"))
	assert.Nil(test, pool.Add(duplicateTransaction))
	assert.Equal(test, pool.Len(), 1)

	pool.Flush()
	assert.Equal(test, pool.Len(), 0)
}
