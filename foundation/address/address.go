// Package address provides functiona;ity for global network addresses.
package address

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Address struct {
	ID    string
	ExtIP *net.IP
	LocIP *net.IP
	Port  uint
}

var (
	ErrMalformedAddressString = errors.New("address is not in format id@ip/port")
	ErrInvalidPortNumber      = errors.New("invalid port number in address")
	ErrInvalidIPAddr          = errors.New("invalid ip in address")
	ErrInvalidID              = errors.New("invalid id")
)

func New(v string) (Address, error) {
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

func (a Address) IP() net.IP {
	if a.ExtIP != nil {
		return *a.ExtIP
	}

	if a.LocIP != nil {
		return *a.LocIP
	}

	panic("nil ips")
}

func (a Address) String() string {
	return fmt.Sprintf("%s@%s/%d", a.ID, a.IP().String(), a.Port)
}

func (a Address) Addr() string {
	ip := a.IP()

	if v4 := ip.To4(); v4 != nil {
		return fmt.Sprintf("%s:%d", a.IP().String(), a.Port)
	}
	return fmt.Sprintf("[%s]:%d", a.IP().String(), a.Port)
}
