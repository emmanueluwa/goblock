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
	privKey *ecdsa.PrivateKey
}

func (privKey PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, privKey.privKey, data)
	if err != nil {
		return nil, err
	}

	return &Signature{
		r: r,
		s: s,
	}, nil
}

func GeneratePrivateKey() PrivateKey {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	//cannot continue if this fails so we panic
	if err != nil {
		panic(err)
	}

	return PrivateKey{
		privKey: privKey,
	}
}

func (privKey PrivateKey) PublicKey() PublicKey {
	return PublicKey{
		pubKey: &privKey.privKey.PublicKey,
	}
}

type PublicKey struct {
	pubKey *ecdsa.PublicKey
}

// access bytes from public key (curve, (x,y)BigInt)
func (pubKey PublicKey) ToSlice() []byte {
	return elliptic.MarshalCompressed(pubKey.pubKey, pubKey.pubKey.X, pubKey.pubKey.Y)
}

// create address using public key
func (pubKey PublicKey) Address() types.Address {
	hash := sha256.Sum256(pubKey.ToSlice())
	//using the last 20 bytes
	return types.AddressFromBytes(hash[len(hash)-20:])
}

type Signature struct {
	r, s *big.Int
}

// verify that signature matches public key (valid)
func (signature Signature) Verify(pubKey PublicKey, data []byte) bool {
	return ecdsa.Verify(pubKey.pubKey, data, signature.r, signature.s)
}
