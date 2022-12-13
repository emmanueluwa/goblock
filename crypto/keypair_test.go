package crypto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivatePublicKey(test *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.PublicKey()

	fmt.Println(privKey, pubKey)
}

func TestKeypairSignatureValid(test *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.PublicKey()
	// address := pubKey.Address()

	message := []byte("Obavan people")
	signature, err := privKey.Sign(message)
	assert.Nil(test, err)

	assert.True(test, signature.Verify(pubKey, message))

}

func TestKeypairSignatureInvalid(test *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.PublicKey()

	message := []byte("Obavan people")
	signature, err := privKey.Sign(message)
	assert.Nil(test, err)

	randomPrivKey := GeneratePrivateKey()
	randomPubKey := randomPrivKey.PublicKey()

	assert.False(test, signature.Verify(randomPubKey, message))
	assert.False(test, signature.Verify(pubKey, []byte("False message")))
}
