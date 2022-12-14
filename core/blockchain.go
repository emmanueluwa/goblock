package core

type Blockchain struct {
	store Storage
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

// uninitalised uint32 --> 4294967295
func (blockchain *Blockchain) HasBlock(height uint32) bool {
	return height <= blockchain.Height()
}

func (blockchain *Blockchain) Height() uint32 {
	//height of chain minus genesis block
	return uint32(len(blockchain.headers)) - 1
}

// internal addblock
func (blockchain *Blockchain) addBlockWithoutValidation(block *Block) error {
	blockchain.headers = append(blockchain.headers, block.Header)

	return blockchain.store.Put(block)
}
