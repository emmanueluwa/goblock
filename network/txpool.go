package network

import (
	"github.com/emmanueluwa/goblock/core"
	"github.com/emmanueluwa/goblock/types"
)

type TxPool struct {
	//hash of transaction and corresponding transaction
	transactions map[types.Hash]*core.Transaction
}

func NewTxPool() *TxPool {
	return &TxPool{
		transactions: make(map[types.Hash]*core.Transaction),
	}
}

// Add will add transaction to pool, caller is responsible for if transaction
// already exists in mempool
func (pool *TxPool) Add(transaction *core.Transaction) error {
	hash := transaction.Hash(core.TxHasher{})
	pool.transactions[hash] = transaction

	return nil
}

// check if transaction is already inside pool
func (pool *TxPool) Has(hash types.Hash) bool {
	_, ok := pool.transactions[hash]
	return ok
}

func (pool *TxPool) Len() int {
	return len(pool.transactions)
}

func (pool *TxPool) Flush() {
	pool.transactions = make(map[types.Hash]*core.Transaction)
}
