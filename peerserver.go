package main

import (
	"gx/ipfs/QmRuZnMorqodado1yeTQiv1i9rmtKj29CjPSsBKM7DFXV4/go-libp2p-net"

	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/core/corenet"
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
func (p *PeerServer) HandleFunc(endpoint string, handler func(*Node, net.Stream)) error {
	list, err := corenet.Listen(p.IpfsNode, endpoint)
	if err != nil {
		return err
	}

	Info.Printf("I am peer: %s and listening at %s \n", p.ID, endpoint)

	// listen asynchronously
	go func() {
		for {
			// @TODO need to handle error gracefully,
			// e.g. tell client that an error occurred
			// @TODO we might have to do this
			// asynchronously
			stream, err := list.Accept()
			if err != nil {
				Error.Println(err)
				continue
			}

			Info.Printf("Connection from: %s\n", stream.Conn().RemotePeer().Pretty())

			isBlacklisted, err := p.IsBlacklisted(p.ID)
			if err != nil {
				Error.Println(err)
				continue
			}

			if isBlacklisted {
				Info.Println("Node is blacklisted, connection will be aborted")
			} else {
				err := p.AddKnown(p.ID)
				if err != nil {
					Error.Println(err)
					continue
				}

				handler(p.Node, stream)
			}

			stream.Close()
		}
	}()

	return nil
}
