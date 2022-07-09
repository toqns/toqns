// Package p2p provides p2p networking functionality.
package p2p

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/toqns/toqns/foundation/address"
)

// ErrUnsupportedProtocol is returned when an unsupported protocol is provided.
var ErrUnsupportedProtocol = errors.New("unsupported protocol")

// Node represents a node on the p2p network.
type Node struct {
	Address address.Address
	conn    net.Conn
	Encoder RequestEncoder
	Decoder RequestDecoder
	Handler Handler
}

// handleRequest handles incoming requests.
func (n *Node) handleRequest(from net.IP, b []byte) error {
	r := Request{}
	if err := n.Decoder.Unmarshal(b, &r); err != nil {
		return fmt.Errorf("decoding request: %w", err)
	}

	if from.IsPrivate() || from.IsUnspecified() {
		r.From.LocIP = &from
	} else {
		r.From.ExtIP = &from
	}

	r.From.ExtIP = &from
	r.Response = &Response{From: r.To, To: r.From}
	if err := n.Handler.Serve(r.Response, &r); err != nil {
		r.Response.StatusCode = StatusInternalServerError
		r.Response.Status = err.Error()
	}

	// TODO: Send response.

	return nil
}

// handleUDPConnections is the connection handler for UDP connections.
func (n *Node) handleUDPConnections() error {
	conn, ok := n.conn.(*net.UDPConn)
	if !ok {
		return fmt.Errorf("connection should be *net.UDPConn, but got %T", conn)
	}

	fmt.Println("Waiting for message...")
	for {
		buf := make([]byte, 2048)
		s, addr, err := conn.ReadFrom(buf)
		if err != nil {
			// TODO: Handle error?
			fmt.Println("Error:" + err.Error())
			continue
		}
		fromIP := net.ParseIP(addr.String())
		if err := n.handleRequest(fromIP, []byte(buf[:s])); err != nil {
			// TODO: Handle error?
			fmt.Println("Error:" + err.Error())
			continue
		}
	}
}

// listenAndServeUDP sets up a UDP listener and starts the
// connection handler.
func (n *Node) listenAndServeUDP() error {
	s, err := net.ResolveUDPAddr("udp", n.Address.Addr())
	if err != nil {
		return fmt.Errorf("resolving udp address: %w", err)
	}

	conn, err := net.ListenUDP("udp", s)
	if err != nil {
		return fmt.Errorf("creating listener: %w", err)
	}
	n.conn = conn

	return n.handleUDPConnections()
}

// ListenAndServe is a blocking function that will start the network listener for
// the selected procotol.
func (n *Node) ListenAndServe() error {
	if n.Handler == nil {
		return fmt.Errorf("nil handler")
	}

	if n.Encoder == nil {
		return fmt.Errorf("nil encoder")
	}

	if n.Decoder == nil {
		return fmt.Errorf("nil decoder")
	}

	switch strings.ToLower(n.Address.Proto) {
	case "udp":
		return n.listenAndServeUDP()
	default:
		return ErrUnsupportedProtocol
	}
}

// Shutdown gracefully stops the node.
func (n *Node) Shutdown(ctx context.Context) error {
	return n.conn.Close()
}
