package p2p

import (
	
	"fmt"
	"net"
)

//TCPeer represents the remote node over a TCP established connection
type TCPPeer struct {
	// conn is underlying connection of the peer
	conn net.Conn

	// if we dial and retrive a connection => outbound == true
	//if we accept and retrive a connection => outbound == false
	outbound bool

}

//Close implements the Peer interface
func (p *TCPPeer) Close () error{
	return p.conn.Close()
} 

func NewTCPPeer(conn net.Conn , outbound bool) *TCPPeer {
	return &TCPPeer{
		conn: conn,
		outbound: outbound,
	}
}

type TCPTransportOpts struct {
	ListenAddr string
	HandshakeFunc HandshakeFunc
	Decoder Decoder
	OnPeer func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener      net.Listener
	rpcch chan RPC

	
}

// Consume implements the Transport interface, which will return a read-only channel
// for reading the incoming messages received from another peer in the network.
func (t *TCPTransport) Consume () <-chan RPC{
	return t.rpcch
}

func NewTCPTransport (opts TCPTransportOpts) *TCPTransport {
return &TCPTransport{
	TCPTransportOpts: opts,
	rpcch: make(chan RPC),
}

}

func (t *TCPTransport) ListenAndAccept() error {
var err error

t.listener,err = net.Listen("tcp",t.ListenAddr)
if err != nil {
	return err
}
  go t.startAcceptLoop()

  return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn,err :=t.listener.Accept()
		if err!=nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}

		fmt.Printf("New incoming connection %+v\n",conn )

		go t.handleConn(conn)
	}
}


func (t *TCPTransport) handleConn(conn net.Conn) {
	
var err error
	defer func(){
		fmt.Printf("Dropping peer connection : %s",err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn,true)

	if err:= t.HandshakeFunc(peer); err!=nil{
		conn.Close()
		fmt.Printf("TCP handshake error: %s\n",err)
		return
	}

	if t.OnPeer !=nil {
		if err = t.OnPeer(peer); err != nil {
			return 
		}
	}

	//Read loop
	rpc:=RPC{}
	for{
	
         // To loop the errors for debug
		// if err := t.Decoder.Decode(conn,&rpc); err!=nil {
		// 	fmt.Printf("TCP error: %s\n",err)
		// 	continue

		// Dropping the connection as soon as error occur 
		// We should find a way to assert specific error to drop connection	
		if err := t.Decoder.Decode(conn,&rpc); err!=nil {
		fmt.Printf("TCP error: %s\n",err)
		return
 }
 rpc.From=conn.RemoteAddr()
//  fmt.Printf("message: %+v\n",rpc)
t.rpcch <- rpc
}
}
    