package main

import (
	"fmt"
	"foreverstore/p2p"
	"log"
)

func OnPeer (peer p2p.Peer) error { 

	//case 1
	// fmt.Println("doing some logic with the peer outside of TCPTransport")
	

	peer.Close()
	return nil
}

func main(){
	tcpOpts:=p2p.TCPTransportOpts{
		ListenAddr: ":3000",
		Decoder: p2p.DefaultDecoder{} ,
		HandshakeFunc:p2p.NOPHandshakeFunc,
		OnPeer: OnPeer,
	}
	tr:=p2p.NewTCPTransport(tcpOpts)

	go func(){
		for {
			msg := <-tr.Consume()
			// fmt.Println("Its from main")
			fmt.Printf("%+v\n",msg)
		}
	}()

	if err:=tr.ListenAndAccept(); err!=nil{
		log.Fatal(err)
	}
	select{}
// fmt.Println("I am working")
}