package key_test

import (
	"testing"

	"github.com/toqns/toqns/business/key"
)

// Success and failure markers.
const (
	success = "\u2713"
	failed  = "\u2717"
)

func TestKey(t *testing.T) {
	t.Log("Given the need to work with keys.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen working with a new key.", testID)
		{
			k, err := key.New()
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create new key: %v.", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create new key.", success, testID)

			privKeyStr, _ := k.PrivateKeyString()
			k2, err := key.Restore(privKeyStr)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to restore key: %v.", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to restore key.", success, testID)

			k2PrivKeyStr, _ := k2.PrivateKeyString()
			if privKeyStr != k2PrivKeyStr {
				t.Fatalf("\t%s\tTest %d:\tShould have matching private keys.", failed, testID)
			}
			t.Logf("\t%s\tTest %d:\tShould have matching private keys.", success, testID)
		}
	}
}
