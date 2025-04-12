package p2p

import "net"

// RPC represents any arbitrary data that
// is bring sent over the each transport between nodes in the network
type RPC struct {
	From net.Addr
	Payload []byte
}