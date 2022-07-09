package key

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
)

type Key struct {
	privateKey *ecdsa.PrivateKey
}

func New() (Key, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return Key{}, fmt.Errorf("generating ecdsa key: %w", err)
	}

	return Key{privateKey: privateKey}, nil
}

func Restore(privKey string) (Key, error) {
	h, err := hex.DecodeString(privKey)
	if err != nil {
		return Key{}, fmt.Errorf("decoding key string: %w", err)
	}

	key, err := x509.ParseECPrivateKey(h)
	if err != nil {
		return Key{}, fmt.Errorf("parsing key: %w", err)
	}

	return Key{privateKey: key}, nil
}

func (k Key) PrivateKeyString() (string, error) {
	b, err := x509.MarshalECPrivateKey(k.privateKey)
	if err != nil {
		return "", fmt.Errorf("marshalling key: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func (k Key) PublicKey() ecdsa.PublicKey {
	return k.privateKey.PublicKey
}

func (k Key) PublicKeyString() string {
	b := elliptic.Marshal(elliptic.P256(), k.privateKey.PublicKey.X, k.privateKey.PublicKey.Y)
	return hex.EncodeToString(b)
}

func (k Key) Address(d string) (Address, error) {
	a, err := AddressFromKey(d, k.privateKey)
	if err != nil {
		return "", err
	}

	return a, nil
}
