package core

import (
	"fmt"

	"github.com/emmanueluwa/goblock/crypto"
)

type Transaction struct {
	Data []byte

	PublicKey crypto.PublicKey
	Signature *crypto.Signature
}

// signing transaction
func (transaction *Transaction) Sign(privKey crypto.PrivateKey) error {
	signature, err := privKey.Sign(transaction.Data)
	if err != nil {
		return err
	}

	transaction.PublicKey = privKey.PublicKey()
	transaction.Signature = signature

	return nil
}

// checking signature with transaction is valid
func (transaction *Transaction) Verify() error {
	if transaction.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}

	if !transaction.Signature.Verify(transaction.PublicKey, transaction.Data) {
		return fmt.Errorf("invalid transaction signature")
	}

	return nil
}
