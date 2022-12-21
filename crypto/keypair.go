package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"math/big"

	"github.com/emmanueluwa/goblock/types"
)

type PrivateKey struct {
	key *ecdsa.PrivateKey
}

func (k PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.key, data)
	if err != nil {
		return nil, err
	}

	return &Signature{
		R: r,
		S: s,
	}, nil
}

func NewPrivateKeyFromReader(reader io.Reader) PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), reader)
	//cannot continue if this fails so we panic
	if err != nil {
		panic(err)
	}

	return PrivateKey{
		key: key,
	}
}

func GeneratePrivateKey() PrivateKey {
	return NewPrivateKeyFromReader(rand.Reader)
}

func (k PrivateKey) PublicKey() PublicKey {
	return elliptic.MarshalCompressed(k.key.PublicKey, k.key.PublicKey.X, k.key.PublicKey.Y)
}

type PublicKey []byte

func (k PublicKey) String() string {
	return hex.EncodeToString(k)
}

func (k PublicKey) Address() types.Address {
	hash := sha256.Sum256(k)

	return types.AddressFromBytes(hash[len(hash)-20:])
}

type Signature struct {
	S *big.Int
	R *big.Int
}

func (signature Signature) String() string {
	b := append(signature.S.Bytes(), signature.R.Bytes()...)
	return hex.EncodeToString(b)
}

// verify that signature matches public key (valid)
func (signature Signature) Verify(pubKey PublicKey, data []byte) bool {
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), pubKey)
	key := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	return ecdsa.Verify(key, data, signature.R, signature.S)
}
