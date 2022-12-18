package core

import (
	"bytes"
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

func TestEncodeDecode(test *testing.T) {
	transaction := randomTransactionWithSignature(test)
	buff := &bytes.Buffer{}
	assert.Nil(test, transaction.Encode(NewGobTxEncoder(buff)))

	transactionDecoded := new(Transaction)
	assert.Nil(test, transactionDecoded.Decode(NewGobTxDecoder(buff)))
	assert.Equal(test, &transaction, transactionDecoded)
}

func randomTransactionWithSignature(test *testing.T) Transaction {
	privKey := crypto.GeneratePrivateKey()
	transaction := Transaction{
		Data: []byte("meow"),
	}
	assert.Nil(test, transaction.Sign(privKey))

	return transaction
}
