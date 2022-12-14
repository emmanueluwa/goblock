package core

type Storage interface {
	Put(*Block) error
}

type MemoryStore struct {
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (store *MemoryStore) Put(block *Block) error {
	return nil
}
