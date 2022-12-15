package network

import (
	"sort"

	"github.com/emmanueluwa/goblock/core"
	"github.com/emmanueluwa/goblock/types"
)

type TxMapSorter struct {
	transactions []*core.Transaction
}

func NewTxMapSorter(txMap map[types.Hash]*core.Transaction) *TxMapSorter {
	transaction := make([]*core.Transaction, len(txMap))

	i := 0
	for _, val := range txMap {
		transaction[i] = val
		i++
	}

	s := &TxMapSorter{transaction}

	sort.Sort(s)

	return s
}

func (sorter *TxMapSorter) Len() int { return len(sorter.transactions) }

func (sorter *TxMapSorter) Swap(i, j int) {
	sorter.transactions[i], sorter.transactions[j] = sorter.transactions[j], sorter.transactions[i]
}

func (sorter *TxMapSorter) Less(i, j int) bool {
	return sorter.transactions[i].FirstSeen() < sorter.transactions[j].FirstSeen()
}

type TxPool struct {
	//hash of transaction and corresponding transaction
	transactions map[types.Hash]*core.Transaction
}

func NewTxPool() *TxPool {
	return &TxPool{
		transactions: make(map[types.Hash]*core.Transaction),
	}
}

func (pool *TxPool) Transactions() []*core.Transaction {
	sorter := NewTxMapSorter(pool.transactions)
	return sorter.transactions
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
