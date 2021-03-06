package main

import (
	"context"
	"errors"
	"os"
	"reflect"
	"time"

	peer "gx/ipfs/QmZcUPvPhD1Xvk6mwijYF8AfR3mG31S1YsEfHG4khrFPRr/go-libp2p-peer"

	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/core"
	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/core/corenet"
	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/repo/config"
	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/repo/fsrepo"
)

const (
	// nBitsForKeypair sets the strength of keypair
	nBitsForKeypair = 2048

	// BootstrapPeerID is the peer id of agora's bootstrap node
	BootstrapPeerID = "QmdtfJBMitotUWBX5YZ6rYeaYRFu6zfXXMZP6fygEWK2iu"
	// BootstrapMultiAddr is the ipfs address of agora's bootstrap node
	BootstrapMultiAddr = "/ip4/54.178.171.10/tcp/4001/ipfs/" + BootstrapPeerID
)

func monitorPeers(n *Node) {
	go func() {
		for {
			printPeers(n.IpfsNode)
			time.Sleep(2 * time.Second)
		}
	}()
}

func printPeers(n *core.IpfsNode) {
	conns := n.PeerHost.Network().Conns()

	Info.Println("---- PeerList")
	for _, c := range conns {
		pid := c.RemotePeer()
		addr := c.RemoteMultiaddr()

		Info.Println(pid, "\t", addr, "\t", n.Peerstore.LatencyEWMA(pid))
	}

}

// Node provides an abstraction for IpfsNode and is the prefered way
// of accessing Nodes in our application. Note that IpfsNode is an
// embedded type.
type Node struct {
	*core.IpfsNode
	*Model

	ID     string
	cancel context.CancelFunc
}

// NewNode creates a new Node from an existing node repository
func NewNode(path string) (*Node, error) {
	// Open and check node repository
	r, err := fsrepo.Open(path)
	if err != nil {
		return nil, err
	}

	// Run Node
	cfg := &core.BuildCfg{
		Repo:   r,
		Online: true,
	}

	ctx, cancel := context.WithCancel(context.Background())
	node, err := core.NewNode(ctx, cfg)
	if err != nil {
		cancel()
		return nil, err
	}

	// Open Node's DB Instance
	db, err := OpenDB(path + "/agora.db")
	if err != nil {
		cancel()
		return nil, err
	}

	return &Node{
		IpfsNode: node,
		Model:    db,
		ID:       node.Identity.Pretty(),
		cancel:   cancel,
	}, nil
}

// CreateNodeIfNotExists creates a new ipfs repo at given location if
// needed and returns an agora node istance
func CreateNodeIfNotExists(path string) (*Node, error) {
	if !Exists(path) {
		err := NewNodeRepo(path, nil)
		if err != nil {
			return nil, err
		}
	}

	return NewNode(MyNodePath)
}

var ErrSkipBlacklisted = errors.New("Peer is blacklisted, will skip this request")

// Request is the generalized method to connect to another peer and
// send requests and receive responses. This is used by Client defined
// in peerapi.go and should not be used directly.
func (n *Client) Request(targetPeer string, path string, body interface{}, resp interface{}) error {

	// @TODO not the most elegant solution due to time pressure
	isBlacklist, err := n.Node.IsBlacklisted(targetPeer)
	if err != nil {
		return err
	}

	if isBlacklist {
		return ErrSkipBlacklisted
	}

	// Check if Node hash is valid
	target, err := peer.IDB58Decode(targetPeer)
	if err != nil {
		return err
	}

	// Connect to targetPeer
	stream, err := corenet.Dial(n.Node.IpfsNode, target, path)
	if err != nil {
		return err
	}

	// This gives you a warning if you accidentially send a
	// pointer instead of the struct as body, note that the
	// warning will not stop the transaction
	if reflect.ValueOf(resp).Kind() != reflect.Ptr {
		Warning.Println("You must pass resp by &reference and not by value. This is not done for a request to", targetPeer, path)
	}

	// Exchange request and response
	WriteJSON(stream, &body)
	ReadJSON(stream, &resp)

	return nil
}

// NewNodeRepo will create a new data and configuration folder for a
// new IPFS node at the provided location
func NewNodeRepo(repoRoot string, addr *config.Addresses) error {
	err := os.MkdirAll(repoRoot, 0755)
	if err != nil {
		return err
	}

	if fsrepo.IsInitialized(repoRoot) {
		return errors.New("Repo already exists")
	}

	conf, err := config.Init(os.Stdout, nBitsForKeypair)
	if err != nil {
		return err
	}

	if addr != nil {
		conf.Addresses = *addr
	}

	own, err := config.ParseBootstrapPeer(BootstrapMultiAddr)
	if err != nil {
		return err
	}

	defaults, err := config.DefaultBootstrapPeers()
	if err != nil {
		return err
	}

	bps := []config.BootstrapPeer{own}
	bps = append(bps, defaults...)

	// Add our own bootstrap node
	conf.SetBootstrapPeers(bps)

	fsrepo.Init(repoRoot, conf)
	if err != nil {
		return err
	}

	return initializeIpnsKeyspace(repoRoot)
}
