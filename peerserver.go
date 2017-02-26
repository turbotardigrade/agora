package main

import (
	"gx/ipfs/QmQx1dHDDYENugYgqA22BaBrRfuv1coSsuPiM7rYh1wwGH/go-libp2p-net"

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

	Info.Printf("I am peer: %s and listening at %s \n", p.Identity.Pretty(), endpoint)

	// listen asynchronously
	go func() {
		for {
			stream, err := list.Accept()
			if err != nil {
				Error.Println(err)
				continue
			}

			Info.Printf("Connection from: %s\n", stream.Conn().RemotePeer().Pretty())

			var blacklisted bool

			BoltGet(blacklistBucket, p.Identity.Pretty(), blacklisted)

			if blacklisted {
				Info.Println("Node is blacklisted, connection will be aborted")
				continue
			}

			handler(stream)
			stream.Close()
		}
	}()

	return nil
}

func AddNodeToBlacklist(identity string) error {
	return BoltSet(blacklistBucket, identity, true)
}
