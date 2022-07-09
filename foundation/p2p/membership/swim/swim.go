// Package swim provides functionality for the SWIM member management protocol.
package swim

import (
	"sync"

	"github.com/toqns/toqns/foundation/address"
)

type SwimManager struct {
	mu     sync.Mutex
	nodes  map[address.Address]struct{}
	joined []address.Address
	failed []address.Address
}

func New() *SwimManager {
	return &SwimManager{
		nodes: make(map[address.Address]struct{}),
	}
}
