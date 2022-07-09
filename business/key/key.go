// Package key provides functionality for working with EDCSA keys.
package key

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	// NodeAddress is a designation for node addresses.
	NodeAddress = "nd"

	// AccountAddress is a designation for (wallet) account addresses.
	AccountAddress = "ac"
)

// Key represents a key for use on a Toqns network.
type Key struct {
	privateKey *ecdsa.PrivateKey
}

// New returns a newly initialized key.
func New() (Key, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return Key{}, fmt.Errorf("generating ecdsa key: %w", err)
	}

	return Key{privateKey: privateKey}, nil
}

// Restore parses the provided private key string and returns a Key.
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

// RestoreFromFile reads the provided file and parses the content into a Key.
func RestoreFromFile(name string) (Key, error) {
	file, err := os.OpenFile(name, os.O_RDONLY, 0600)
	if err != nil {
		return Key{}, fmt.Errorf("file %s: %w", name, err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return Key{}, fmt.Errorf("reading file: %w", err)
	}

	block, _ := pem.Decode(b)
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return Key{}, fmt.Errorf("parsing key: %w", err)
	}

	return Key{privateKey: privateKey}, nil
}

// PrivateKeyString returns the key's private key string.
func (k Key) PrivateKeyString() (string, error) {
	b, err := x509.MarshalECPrivateKey(k.privateKey)
	if err != nil {
		return "", fmt.Errorf("marshalling key: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// PublicKey returns the public key as ecdsa.
func (k Key) PublicKey() ecdsa.PublicKey {
	return k.privateKey.PublicKey
}

// PublicKeyString returns the public key string.
func (k Key) PublicKeyString() string {
	b := elliptic.Marshal(elliptic.P256(), k.privateKey.PublicKey.X, k.privateKey.PublicKey.Y)
	return hex.EncodeToString(b)
}

// Address returns the address of this key.
//
// Requires a designation, such as NodeAddress or AccountAddress.
func (k Key) Address(d string) (Address, error) {
	a, err := AddressFromKey(d, k.privateKey)
	if err != nil {
		return "", err
	}

	return a, nil
}

// Save stores the private key as a pem file.
//
// Use RestoreFromFile to restore the key from the file.
func (k Key) Save(name string) error {
	dir, _ := filepath.Split(name)
	os.MkdirAll(dir, 0700)

	b, err := x509.MarshalECPrivateKey(k.privateKey)
	if err != nil {
		return fmt.Errorf("marshalling key: %w", err)
	}

	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer file.Close()

	if err := pem.Encode(file, &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}); err != nil {
		return fmt.Errorf("encoding pem: %w", err)
	}

	return nil
}
