package network

import (
	"sync"

	"github.com/emmanueluwa/goblock/core"
	"github.com/emmanueluwa/goblock/types"
)

type TransactionPool struct {
	//all transactions seen
	all *TxSortedMap
	//transactions waiting to be placed in a block
	pending *TxSortedMap
	//when the pool is full the oldest transaction is pruned
	maxLength int
}

func NewTransactionPool(maxLength int) *TransactionPool {
	return &TransactionPool{
		all:       NewTxSortedMap(),
		pending:   NewTxSortedMap(),
		maxLength: maxLength,
	}
}

// Add will add transaction to pool, caller is responsible for if transaction
// already exists in mempool
func (pool *TransactionPool) Add(transaction *core.Transaction) {
	//prune oldest transaction in pool
	if pool.all.Count() == pool.maxLength {
		oldest := pool.all.First()
		pool.all.Remove(oldest.Hash(core.TxHasher{}))
	}

	if !pool.all.Contains(transaction.Hash(core.TxHasher{})) {
		pool.all.Add(transaction)
		pool.pending.Add(transaction)
	}
}

// check if transaction is already inside pool
func (pool *TransactionPool) Contains(hash types.Hash) bool {
	return pool.all.Contains(hash)
}

// returns slice of transactions in pending pool
func (pool *TransactionPool) Pending() []*core.Transaction {
	return pool.pending.transactions.Data
}

func (pool *TransactionPool) ClearPending() {
	pool.pending.Clear()
}

func (pool *TransactionPool) PendingCount() int {
	return pool.pending.Count()
}

type TxSortedMap struct {
	lock         sync.RWMutex
	lookup       map[types.Hash]*core.Transaction
	transactions *types.List[*core.Transaction]
}

func NewTxSortedMap() *TxSortedMap {
	return &TxSortedMap{
		lookup:       make(map[types.Hash]*core.Transaction),
		transactions: types.NewList[*core.Transaction](),
	}
}

func (transactions *TxSortedMap) First() *core.Transaction {
	transactions.lock.RLock()
	defer transactions.lock.RUnlock()

	first := transactions.transactions.Get(0)
	return transactions.lookup[first.Hash(core.TxHasher{})]
}

func (transactions *TxSortedMap) Get(hash types.Hash) *core.Transaction {
	transactions.lock.RLock()
	defer transactions.lock.RUnlock()

	return transactions.lookup[hash]
}

func (transactions *TxSortedMap) Add(transaction *core.Transaction) {
	hash := transaction.Hash(core.TxHasher{})

	transactions.lock.Lock()
	defer transactions.lock.Unlock()

	if _, ok := transactions.lookup[hash]; !ok {
		transactions.lookup[hash] = transaction
		transactions.transactions.Insert(transaction)
	}
}

func (transactions *TxSortedMap) Remove(hash types.Hash) {
	transactions.lock.Lock()
	defer transactions.lock.Unlock()

	transactions.transactions.Remove(transactions.lookup[hash])
	delete(transactions.lookup, hash)
}

func (transactions *TxSortedMap) Count() int {
	transactions.lock.RLock()
	defer transactions.lock.RUnlock()

	return len(transactions.lookup)
}

func (transactions *TxSortedMap) Contains(hash types.Hash) bool {
	transactions.lock.RLock()
	defer transactions.lock.RUnlock()

	_, ok := transactions.lookup[hash]
	return ok
}

func (transactions *TxSortedMap) Clear() {
	transactions.lock.Lock()
	defer transactions.lock.Unlock()

	transactions.lookup = make(map[types.Hash]*core.Transaction)
	transactions.transactions.Clear()
}
