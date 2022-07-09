// Package membership provides functionality for p2p membership management.
package membership

type Manager interface {
	Broadcast([]byte) error
}
