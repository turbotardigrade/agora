package main

import (
	"fmt"
	"gx/ipfs/QmdXimY9QHaasZmw6hWojWnCJvfgxETjZQfg9g6ZrA9wMX/go-libp2p-net"

	"github.com/ipfs/go-ipfs/core/corenet"
)

type PeerServer struct {
	*Node
}

func NewPeerServer(node *Node) PeerServer {
	return PeerServer{Node: node}
}

func (p *PeerServer) handleFunc(pattern string, handler func(net.Stream)) error {
	list, err := corenet.Listen(p.IpfsNode, pattern)
	if err != nil {
		return err
	}

	fmt.Printf("I am peer: %s and listening at %s \n", p.Identity.Pretty(), pattern)

	// listen asynchronously
	go func() {
		for {
			stream, err := list.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Printf("Connection from: %s\n", stream.Conn().RemotePeer().Pretty())

			handler(stream)
			stream.Close()
		}
	}()

	return nil
}
