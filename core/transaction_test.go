package core

import (
	"testing"

	"github.com/emmanueluwa/goblock/crypto"
	"github.com/stretchr/testify/assert"
)

func TestSignTransaction(test *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	data := []byte("meow")
	transaction := &Transaction{
		Data: data,
	}

	assert.Nil(test, transaction.Sign(privKey))
	assert.NotNil(test, transaction.Signature)
}

func TestVerifyTransaction(test *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	data := []byte("meow")
	transaction := &Transaction{
		Data: data,
	}

	assert.Nil(test, transaction.Sign(privKey))
	assert.Nil(test, transaction.Verify())

	randomPrivKey := crypto.GeneratePrivateKey()
	transaction.From = randomPrivKey.PublicKey()

	assert.NotNil(test, transaction.Verify())
}

func randomTransactionWithSignature(test *testing.T) *Transaction {
	privKey := crypto.GeneratePrivateKey()
	transaction := &Transaction{
		Data: []byte("meow"),
	}
	assert.Nil(test, transaction.Sign(privKey))

	return transaction
}
