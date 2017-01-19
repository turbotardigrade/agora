package main

import (
	"fmt"
	"gx/ipfs/QmdXimY9QHaasZmw6hWojWnCJvfgxETjZQfg9g6ZrA9wMX/go-libp2p-net"

	"github.com/ipfs/go-ipfs/core/corenet"
)

// PeerServer asynchronous (non-blocking) server that registers
// handler functions listening to specific endpoints like /health
type PeerServer struct {
	*Node
}

// NewPeerServer constructs PeerServer given a node on which the
// server binds its services
func NewPeerServer(node *Node) PeerServer {
	return PeerServer{Node: node}
}

// HandleFunc registers a function which get called to handle incoming
// requests triggered on given endpoint
func (p *PeerServer) HandleFunc(endpoint string, handler func(net.Stream)) error {
	list, err := corenet.Listen(p.IpfsNode, endpoint)
	if err != nil {
		return err
	}

	fmt.Printf("I am peer: %s and listening at %s \n", p.Identity.Pretty(), endpoint)

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
