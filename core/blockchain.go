package core

import (
	"fmt"
	"sync"

	"github.com/go-kit/log"
)

type Blockchain struct {
	logger log.Logger
	//for json rpc(recieving block) and protocol, other nodes can ask blocks to sync
	store Storage
	//for distirbuted env, blockchain data structure needs to be threadsafe
	lock sync.RWMutex
	//maintain in memory, RAM intensive instead of Disk intensive
	headers   []*Header
	validator Validator
	// TODO: MAKE THIS AN INTERFACE
	contractState *State
}

func NewBlockchain(log log.Logger, genesis *Block) (*Blockchain, error) {
	blockchain := &Blockchain{
		contractState: NewState(),
		headers:       []*Header{},
		store:         NewMemoryStore(),
		logger:        log,
	}
	blockchain.validator = NewBlockValidator(blockchain)
	err := blockchain.addBlockWithoutValidation(genesis)

	return blockchain, err
}

func (blockchain *Blockchain) SetValidator(validator Validator) {
	blockchain.validator = validator
}

func (blockchain *Blockchain) AddBlock(block *Block) error {
	if err := blockchain.validator.ValidateBlock(block); err != nil {
		return err
	}

	for _, transaction := range block.Transactions {
		blockchain.logger.Log("message", "executing code", "len", len(transaction.Data), "hash", transaction.Hash(&TxHasher{}))

		vm := NewVM(transaction.Data, blockchain.contractState)
		if err := vm.Run(); err != nil {
			return err
		}

		fmt.Printf("STATE: %+v\n", vm.contractState)

		result := vm.stack.Pop()
		fmt.Printf("VM Result: %+v\n", result)
	}

	//validation already done if no error
	return blockchain.addBlockWithoutValidation(block)
}

func (blockchain *Blockchain) GetHeader(height uint32) (*Header, error) {
	if height > blockchain.Height() {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}

	blockchain.lock.Lock()
	defer blockchain.lock.Unlock()

	return blockchain.headers[height], nil
}

// uninitalised uint32 --> 4294967295
func (blockchain *Blockchain) HasBlock(height uint32) bool {
	return height <= blockchain.Height()
}

func (blockchain *Blockchain) Height() uint32 {
	blockchain.lock.RLock()
	defer blockchain.lock.RUnlock()

	//height of chain minus genesis block
	return uint32(len(blockchain.headers)) - 1
}

// internal addblock
func (blockchain *Blockchain) addBlockWithoutValidation(block *Block) error {
	blockchain.lock.Lock()
	blockchain.headers = append(blockchain.headers, block.Header)
	blockchain.lock.Unlock()

	blockchain.logger.Log(
		"message", "new block",
		"hash", block.Hash(BlockHasher{}),
		"height", block.Height,
		"transactions", len(block.Transactions),
	)

	return blockchain.store.Put(block)
}
