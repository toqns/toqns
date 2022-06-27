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
				expString string
				expAddr   string
			}{
				{name: "publicip4", id: "1234", ip: net.ParseIP("8.8.8.8"), port: 3000, expString: "1234@8.8.8.8/3000", expAddr: "8.8.8.8:3000"},
				{name: "localip4", id: "1234", ip: net.ParseIP("192.168.0.100"), port: 3000, expString: "1234@192.168.0.100/3000", expAddr: "192.168.0.100:3000"},
				{name: "undefip4", id: "1234", ip: net.ParseIP("0.0.0.0"), port: 3000, expString: "1234@0.0.0.0/3000", expAddr: "0.0.0.0:3000"},
				{name: "publicip6", id: "1234", ip: net.ParseIP("2001:4860:4802:32::a"), port: 3000, expString: "1234@2001:4860:4802:32::a/3000", expAddr: "[2001:4860:4802:32::a]:3000"},
				{name: "localip6", id: "1234", ip: net.ParseIP("fd4d:779d:5bd2:68ba:1234:abcd:4321:dcba"), port: 3000, expString: "1234@fd4d:779d:5bd2:68ba:1234:abcd:4321:dcba/3000", expAddr: "[fd4d:779d:5bd2:68ba:1234:abcd:4321:dcba]:3000"},
				{name: "undefip6", id: "1234", ip: net.ParseIP("::"), port: 3000, expString: "1234@::/3000", expAddr: "[::]:3000"},
			}

			for _, tc := range tt {
				t.Run(tc.name, func(t *testing.T) {
					addr := address.Address{ID: tc.id, ExtIP: &tc.ip, Port: tc.port}
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
				{"noPortv4", "1234@8.8.8.8", address.ErrMalformedAddressString},
				{"noPortv6", "1234@2001:4860:4802:32::a", address.ErrMalformedAddressString},
				{"negativePortv4", "1234@8.8.8.8/-5", address.ErrInvalidPortNumber},
				{"negativePortv6", "1234@2001:4860:4802:32::a/-5", address.ErrInvalidPortNumber},
				{"wrongIPv4", "1234@8.8.8.-8/3000", address.ErrInvalidIPAddr},
				{"wrongIPv6", "1234@2001:4860:4802:32::q/3000", address.ErrInvalidIPAddr},
				{"noID4", "8.8.8.8/3000", address.ErrMalformedAddressString},
				{"noIDv6", "2001:4860:4802:32::a/3000", address.ErrMalformedAddressString},
				{"emptyID4", "@8.8.8.8/3000", address.ErrInvalidID},
				{"emptyIDv6", "@2001:4860:4802:32::a/3000", address.ErrInvalidID},
				{"validPubv4", "1234@8.8.8.8/3000", nil},
				{"validPubv6", "1234@2001:4860:4802:32::a/3000", nil},
				{"validLocv4", "1234@192.168.0.100/3000", nil},
				{"validLocv6", "1234@fd4d:779d:5bd2:68ba:1234:abcd:4321:dcba/3000", nil},
				{"validUnspecv4", "1234@0.0.0.0/3000", nil},
				{"validUnspecv6", "1234@::/3000", nil},
			}

			for _, tc := range tt {
				t.Run(tc.name, func(t *testing.T) {
					addr, err := address.New(tc.val)
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
