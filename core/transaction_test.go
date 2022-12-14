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
	transaction.PublicKey = randomPrivKey.PublicKey()

	assert.NotNil(test, transaction.Verify())
}
