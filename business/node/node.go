// Package node provides functionality for Toqns nodes.
package node

import (
	"encoding/json"
	"fmt"

	"github.com/toqns/toqns/business/key"
	"github.com/toqns/toqns/foundation/address"
	"github.com/toqns/toqns/foundation/p2p"
	"go.uber.org/zap"
)

// NodeConfig contains configuration details for nodes.
type NodeConfig struct {
	Address     string
	Port        int
	Protocol    string
	NodeKeyFile string
}

// Node repersents a node on the Toqns network.
type Node struct {
	*p2p.Node
	log *zap.SugaredLogger
}

// New returns an initialized Node based on the provided configuration.
func New(log *zap.SugaredLogger, cfg NodeConfig) (*Node, error) {
	k, err := key.RestoreFromFile(cfg.NodeKeyFile)
	if err != nil {
		// TODO: Create key if not exists.
		return nil, fmt.Errorf("reading node key: %w", err)
	}

	// Node IDs are the key addresses.
	id, err := k.Address(key.NodeAddress)
	if err != nil {
		return nil, fmt.Errorf("getting node ID: %w", err)
	}

	addr, err := address.Parse(fmt.Sprintf("%s@%s/%d/%s", id, cfg.Address, cfg.Port, cfg.Protocol))
	if err != nil {
		return nil, fmt.Errorf("parsing address: %w", err)
	}

	return &Node{
		Node: &p2p.Node{
			Address: addr,
			Decoder: p2p.RequestDecoderFunc(json.Unmarshal),
			Encoder: p2p.RequestEncoderFunc(json.Marshal),
			Handler: p2p.HandleFunc(tmpHandler),
		},
		log: log,
	}, nil
}

// TODO: Handlers should probably not return an error.
func tmpHandler(w p2p.ResponseWriter, r *p2p.Request) error {
	fmt.Println("request received!")
	return nil
}
