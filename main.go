package main

import (
	
	"foreverstore/p2p"
	"log"
)

func main(){
	tcpOpts:=p2p.TCPTransportOpts{
		ListenAddr: ":3000",
		Decoder: p2p.GOBDecoder{} ,
		HandshakeFunc:p2p.NOPHandshakeFunc,
	}
	tr:=p2p.NewTCPTransport(tcpOpts)

	if err:=tr.ListenAndAccept(); err!=nil{
		log.Fatal(err)
	}
	select{}
// fmt.Println("I am working")
}