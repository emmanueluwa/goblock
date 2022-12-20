package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"

	"github.com/emmanueluwa/goblock/types"
)

type PrivateKey struct {
	Key *ecdsa.PrivateKey
}

func (k PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.Key, data)
	if err != nil {
		return nil, err
	}

	return &Signature{
		R: r,
		S: s,
	}, nil
}

func GeneratePrivateKey() PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	//cannot continue if this fails so we panic
	if err != nil {
		panic(err)
	}

	return PrivateKey{
		Key: key,
	}
}

func (k PrivateKey) PublicKey() PublicKey {
	return PublicKey{
		Key: &k.Key.PublicKey,
	}
}

type PublicKey struct {
	Key *ecdsa.PublicKey
}

// access bytes from public key (curve, (x,y)BigInt)
func (k PublicKey) ToSlice() []byte {
	return elliptic.MarshalCompressed(k.Key, k.Key.X, k.Key.Y)
}

// create address using public key
func (pubKey PublicKey) Address() types.Address {
	hash := sha256.Sum256(pubKey.ToSlice())
	//using the last 20 bytes
	return types.AddressFromBytes(hash[len(hash)-20:])
}

type Signature struct {
	R, S *big.Int
}

// verify that signature matches public key (valid)
func (signature Signature) Verify(pubKey PublicKey, data []byte) bool {
	return ecdsa.Verify(pubKey.Key, data, signature.R, signature.S)
}
