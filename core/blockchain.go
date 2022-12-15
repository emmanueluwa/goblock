package core

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

type Blockchain struct {
	//for json rpc(recieving block) and protocol, other nodes can ask blocks to sync
	store Storage
	//for distirbuted env, blockchain data structure needs to be threadsafe
	lock sync.RWMutex
	//maintain in memory, RAM intensive instead of Disk intensive
	headers   []*Header
	validator Validator
}

func NewBlockchain(genesis *Block) (*Blockchain, error) {
	blockchain := &Blockchain{
		headers: []*Header{},
		store:   NewMemoryStore(),
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

	logrus.WithFields(logrus.Fields{
		"height": block.Height,
		"hash":   block.Hash(BlockHasher{}),
	}).Info("adding new block")

	return blockchain.store.Put(block)
}
