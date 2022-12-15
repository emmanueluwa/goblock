package core

import "fmt"

/***
Reason for interface:
- so it can be mocked and tested
- use it to make other stuff
- can be used as default or another can be made by others
***/
type Validator interface {
	ValidateBlock(*Block) error
}

type BlockValidator struct {
	blockchain *Blockchain
}

func NewBlockValidator(blockchain *Blockchain) *BlockValidator {
	return &BlockValidator{
		blockchain: blockchain,
	}
}

func (validator *BlockValidator) ValidateBlock(block *Block) error {
	if validator.blockchain.HasBlock(block.Height) {
		return fmt.Errorf("chain already contains block (%d) with hash (%s)", block.Height, block.Hash(BlockHasher{}))
	}

	//ensuring block comes exactly after current height(no blocks skipped)
	if block.Height != validator.blockchain.Height()+1 {
		return fmt.Errorf("Block (%s) too hight", block.Hash(BlockHasher{}))
	}

	//checking hash of previous block
	prevHeader, err := validator.blockchain.GetHeader(block.Height - 1)
	if err != nil {
		return err
	}

	hash := BlockHasher{}.Hash(prevHeader)
	if hash != block.PreviousBlockHash {
		return fmt.Errorf("The hash of the previous block (%s) is invalid", block.PreviousBlockHash)
	}

	if err := block.Verify(); err != nil {
		return err
	}

	return nil
}
