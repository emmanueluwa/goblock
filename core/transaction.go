package core

import (
	"fmt"

	"github.com/emmanueluwa/goblock/crypto"
	"github.com/emmanueluwa/goblock/types"
)

type Transaction struct {
	Data []byte

	From      crypto.PublicKey
	Signature *crypto.Signature

	//cashed version of transaction data hash
	hash types.Hash
	// firstSeen, timestamp for when tx is first seen locally
	firstSeen int64
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data: data,
	}
}

func (transaction *Transaction) Hash(hasher Hasher[*Transaction]) types.Hash {
	if transaction.hash.IsZero() {
		transaction.hash = hasher.Hash(transaction)
	}
	return transaction.hash
}

// signing transaction
func (transaction *Transaction) Sign(privKey crypto.PrivateKey) error {
	signature, err := privKey.Sign(transaction.Data)
	if err != nil {
		return err
	}

	transaction.From = privKey.PublicKey()
	transaction.Signature = signature

	return nil
}

// checking signature with transaction is valid
func (transaction *Transaction) Verify() error {
	if transaction.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}

	if !transaction.Signature.Verify(transaction.From, transaction.Data) {
		return fmt.Errorf("invalid transaction signature")
	}

	return nil
}

func (transaction *Transaction) SetFirstSeen(t int64) {
	transaction.firstSeen = t
}

func (transaction *Transaction) FirstSeen() int64 {
	return transaction.firstSeen
}
