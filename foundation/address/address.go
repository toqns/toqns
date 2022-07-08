// Package address provides functiona;ity for global network addresses.
package address

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
)

// Address repersents a general address with and ID, IP and Port.
type Address struct {
	ID    string
	ExtIP *net.IP
	LocIP *net.IP
	Port  uint
	Proto string
}

var (
	// ErrMalformedAddressString is  anerror to indicate that an address string isn't properly formatted.
	ErrMalformedAddressString = errors.New("address is not in format id@ip/port/protocol")

	// ErrInvalidPortNumber is an error to indicate that the port of an address is incorrect.
	ErrInvalidPortNumber = errors.New("invalid port number in address")

	// ErrInvalidIPAddr is an error to indicate that the IP of and address is incorrect.
	ErrInvalidIPAddr = errors.New("invalid ip in address")

	// ErrInvalidIPAddr is an error to indicate that the ID of an address is incorrect.
	// ErrInvalidID = errors.New("invalid id")
)

// Parse parses an address string to an Address.
//
// The address string should be in the format: id@ip/port.
// IP can be IPv4 or IPv6.
func Parse(v string) (Address, error) {
	re, err := regexp.Compile(`(?P<id>.+)@(?P<ip>\S+)/(?P<port>\d+)/(?P<proto>.+)`)
	if err != nil {
		return Address{}, fmt.Errorf("regex compilation failed: %v", err)
	}

	res := re.FindStringSubmatch(v)
	if len(res) < 5 {
		return Address{}, ErrMalformedAddressString
	}
	id := res[1]
	ipStr := res[2]
	portStr := res[3]
	proto := res[4]

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return Address{}, ErrInvalidPortNumber
	}
	if port <= 0 {
		return Address{}, ErrInvalidPortNumber
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return Address{}, ErrInvalidIPAddr
	}

	addr := Address{
		ID:    id,
		Port:  uint(port),
		Proto: proto,
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
	return fmt.Sprintf("%s@%s/%d/%s", a.ID, a.IP().String(), a.Port, a.Proto)
}

// Addr returns a host address string as in the format ip:port.
func (a Address) Addr() string {
	ip := a.IP()

	if v4 := ip.To4(); v4 != nil {
		return fmt.Sprintf("%s:%d", a.IP().String(), a.Port)
	}
	return fmt.Sprintf("[%s]:%d", a.IP().String(), a.Port)
}
