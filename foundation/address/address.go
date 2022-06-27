// Package address provides functiona;ity for global network addresses.
package address

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Address repersents a general address with and ID, IP and Port.
type Address struct {
	ID    string
	ExtIP *net.IP
	LocIP *net.IP
	Port  uint
}

var (
	// ErrMalformedAddressString is  anerror to indicate that an address string isn't properly formatted.
	ErrMalformedAddressString = errors.New("address is not in format id@ip/port")

	// ErrInvalidPortNumber is an error to indicate that the port of an address is incorrect.
	ErrInvalidPortNumber = errors.New("invalid port number in address")

	// ErrInvalidIPAddr is an error to indicate that the IP of and address is incorrect.
	ErrInvalidIPAddr = errors.New("invalid ip in address")

	// ErrInvalidIPAddr is an error to indicate that the ID of an address is incorrect.
	ErrInvalidID = errors.New("invalid id")
)

// Parse parses an address string to an Address.
//
// The address string should be in the format: id@ip/port.
// IP can be IPv4 or IPv6.
func Parse(v string) (Address, error) {
	parts := strings.Split(v, "/")
	if len(parts) != 2 {
		// Port missing.
		return Address{}, ErrMalformedAddressString
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return Address{}, ErrInvalidPortNumber
	}
	if port <= 0 {
		return Address{}, ErrInvalidPortNumber
	}

	idhost := strings.Split(parts[0], "@")
	if len(idhost) != 2 {
		// ID or IP missing.
		return Address{}, ErrMalformedAddressString
	}

	id := idhost[0]
	if id == "" {
		return Address{}, ErrInvalidID
	}

	ip := net.ParseIP(idhost[1])
	if ip == nil {
		return Address{}, ErrInvalidIPAddr
	}

	addr := Address{
		ID:   id,
		Port: uint(port),
	}

	if ip.IsPrivate() || ip.IsUnspecified() {
		addr.LocIP = &ip
	} else {
		addr.ExtIP = &ip
	}

	return addr, nil
}

// IP returns the external IP address if set, or else the local IP address.
//
// Panics is no external or local IP address is set.
func (a Address) IP() net.IP {
	if a.ExtIP != nil {
		return *a.ExtIP
	}

	if a.LocIP != nil {
		return *a.LocIP
	}

	panic("nil ips")
}

// String implements the stringer interface and returns the address string.
func (a Address) String() string {
	return fmt.Sprintf("%s@%s/%d", a.ID, a.IP().String(), a.Port)
}

// Addr returns a host address string as in the format ip:port.
func (a Address) Addr() string {
	ip := a.IP()

	if v4 := ip.To4(); v4 != nil {
		return fmt.Sprintf("%s:%d", a.IP().String(), a.Port)
	}
	return fmt.Sprintf("[%s]:%d", a.IP().String(), a.Port)
}
