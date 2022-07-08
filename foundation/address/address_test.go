package address_test

import (
	"net"
	"testing"

	"github.com/toqns/toqns/foundation/address"
)

// Success and failure markers.
const (
	success = "\u2713"
	failed  = "\u2717"
)

func TestAddress(t *testing.T) {
	t.Log("Given the need to work with addresses.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen parsing an address to a string.", testID)
		{
			tt := []struct {
				name      string
				id        string
				ip        net.IP
				port      uint
				proto     string
				expString string
				expAddr   string
			}{
				{"publicip4", "1234", net.ParseIP("8.8.8.8"), 3000, "udp", "1234@8.8.8.8/3000/udp", "8.8.8.8:3000"},
				{"localip4", "1234", net.ParseIP("192.168.0.100"), 3000, "udp", "1234@192.168.0.100/3000/udp", "192.168.0.100:3000"},
				{"undefip4", "1234", net.ParseIP("0.0.0.0"), 3000, "udp", "1234@0.0.0.0/3000/udp", "0.0.0.0:3000"},
				{"publicip6", "1234", net.ParseIP("2001:4860:4802:32::a"), 3000, "udp", "1234@2001:4860:4802:32::a/3000/udp", "[2001:4860:4802:32::a]:3000"},
				{"localip6", "1234", net.ParseIP("fd4d:779d:5bd2:68ba:1234:abcd:4321:dcba"), 3000, "udp", "1234@fd4d:779d:5bd2:68ba:1234:abcd:4321:dcba/3000/udp", "[fd4d:779d:5bd2:68ba:1234:abcd:4321:dcba]:3000"},
				{"undefip6", "1234", net.ParseIP("::"), 3000, "udp", "1234@::/3000/udp", "[::]:3000"},
			}

			for _, tc := range tt {
				t.Run(tc.name, func(t *testing.T) {
					addr := address.Address{ID: tc.id, ExtIP: &tc.ip, Port: tc.port, Proto: tc.proto}
					got := addr.String()
					if got != tc.expString {
						t.Fatalf("\t%s\tTest %d:\tShould get address string %q, but got: %q", failed, testID, tc.expString, got)
					}
					t.Logf("\t%s\tTest %d:\tShould get address string %q.", success, testID, tc.expString)

					got = addr.Addr()
					if got != tc.expAddr {
						t.Fatalf("\t%s\tTest %d:\tShould get address %q, but got: %q", failed, testID, tc.expAddr, got)
					}
					t.Logf("\t%s\tTest %d:\tShould get address %q.", success, testID, tc.expAddr)
				})
			}

		}

		testID = 1
		t.Logf("\tTest %d:\tWhen parsing an address from a string.", testID)
		{
			tt := []struct {
				name string
				val  string
				err  error
			}{
				{"noPortv4", "1234@8.8.8.8/udp", address.ErrMalformedAddressString},
				{"noPortv6", "1234@2001:4860:4802:32::a/udp", address.ErrMalformedAddressString},
				{"noProtov4", "1234@8.8.8.8/proto", address.ErrMalformedAddressString},
				{"noProtov6", "1234@2001:4860:4802:32::a/proto", address.ErrMalformedAddressString},
				{"negativePortv4", "1234@8.8.8.8/-5/udp", address.ErrMalformedAddressString},
				{"negativePortv6", "1234@2001:4860:4802:32::a/-5/udp", address.ErrMalformedAddressString},
				{"wrongIPv4", "1234@8.8.8.-8/3000/udp", address.ErrInvalidIPAddr},
				{"wrongIPv6", "1234@2001:4860:4802:32::q/3000/udp", address.ErrInvalidIPAddr},
				{"noID4", "8.8.8.8/3000/udp", address.ErrMalformedAddressString},
				{"noIDv6", "2001:4860:4802:32::a/3000/udp", address.ErrMalformedAddressString},
				{"emptyID4", "@8.8.8.8/3000/udp", address.ErrMalformedAddressString},
				{"emptyIDv6", "@2001:4860:4802:32::a/3000/udp", address.ErrMalformedAddressString},
				{"validPubv4", "1234@8.8.8.8/3000/udp", nil},
				{"validPubv6", "1234@2001:4860:4802:32::a/3000/udp", nil},
				{"validLocv4", "1234@192.168.0.100/3000/udp", nil},
				{"validLocv6", "1234@fd4d:779d:5bd2:68ba:1234:abcd:4321:dcba/3000/udp", nil},
				{"validUnspecv4", "1234@0.0.0.0/3000/udp", nil},
				{"validUnspecv6", "1234@::/3000/udp", nil},
			}

			for _, tc := range tt {
				t.Run(tc.name, func(t *testing.T) {
					addr, err := address.Parse(tc.val)
					if err != tc.err {
						t.Fatalf("\t%s\tTest %d:\tShould get error %q, but got \"%v\".", failed, testID, tc.err, err)
					}
					t.Logf("\t%s\tTest %d:\tShould get error \"%v\".", success, testID, tc.err)

					if tc.err == nil {
						if addr.String() != tc.val {
							t.Fatalf("\t%s\tTest %d:\tShould get string value %q, but got %q.", failed, testID, tc.val, addr.String())
						}
						t.Logf("\t%s\tTest %d:\tShould get string value %q.", success, testID, tc.val)
					}
				})
			}
		}
	}
}
