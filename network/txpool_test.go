package network

import (
	"math/rand"
	"strconv"
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

func TestSortTransactions(test *testing.T) {
	pool := NewTxPool()
	txLength := 1000

	for i := 0; i < txLength; i++ {
		transaction := core.NewTransaction([]byte(strconv.FormatInt(int64(i), 10)))
		transaction.SetFirstSeen(int64(i * rand.Intn(10000)))
		assert.Nil(test, pool.Add(transaction))
	}

	assert.Equal(test, pool.Len(), txLength)

	transactions := pool.Transactions()
	for i := 0; i < len(transactions)-1; i++ {
		assert.True(test, transactions[i].FirstSeen() < transactions[i+1].FirstSeen())
	}

}
