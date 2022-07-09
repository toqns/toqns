package key

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"strings"
)

// Address is a custom type to support key derived addresses.
type Address string

// AddressFromKey returns an address from a private key.
func AddressFromKey(d string, k *ecdsa.PrivateKey) (Address, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&k.PublicKey)
	if err != nil {
		return "", fmt.Errorf("generating public key der: %w", err)
	}
	addressBytes := publicKeyBytes[len(publicKeyBytes)-20:]
	address := d + hex.EncodeToString(addressBytes)
	return Address(address), nil
}

// AddressFromString parses a string to an address.
// Returns an error if the address doesn't validate.
func AddressFromString(v string) (Address, error) {
	a := Address(v)
	if err := a.Validate(); err != nil {
		return "", err
	}

	return a, nil
}

// String implements the stringer interface.
func (a *Address) String() string {
	return string(*a)
}

// Validate validates an address and returns an error if validation fails.
func (a Address) Validate() error {
	str := string(a)
	if len(str) != 42 {
		return fmt.Errorf("invalid string length")
	}

	// ac: Account
	// nd: Node
	if !strings.HasPrefix(str, "ac") || !strings.HasPrefix(str, "nd") {
		return fmt.Errorf("invalid prefix")
	}

	b, err := hex.DecodeString(str[2:])
	if err != nil {
		return fmt.Errorf("invalid hex value")
	}

	if len(b) != 20 {
		return fmt.Errorf("invalid byte length")
	}

	return nil
}

// IsValid is a wrapper around Validate.
// Returns true or false respective of whether validation passed.
func (a Address) IsValid() bool {
	if err := a.Validate(); err != nil {
		return false
	}
	return true
}
