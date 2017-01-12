package main

import (
	"fmt"
	"gx/ipfs/QmdXimY9QHaasZmw6hWojWnCJvfgxETjZQfg9g6ZrA9wMX/go-libp2p-net"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/corenet"
)

type PeerServer struct {
	ipfsNode *core.IpfsNode
}

func NewPeerServer(node *Node) PeerServer {
	return PeerServer{
		ipfsNode: node.ipfsNode,
	}
}

func (p *PeerServer) handleFunc(pattern string, handler func(net.Stream)) error {
	list, err := corenet.Listen(p.ipfsNode, pattern)
	if err != nil {
		return err
	}

	fmt.Printf("I am peer: %s and listening at %s \n", p.ipfsNode.Identity.Pretty(), pattern)

	// listen asynchronously
	go func() {
		for {
			con, err := list.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			}

			handler(con)
			con.Close()
		}
	}()

	return nil
}
