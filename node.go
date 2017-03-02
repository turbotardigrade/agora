package main

import (
	"context"
	"errors"
	"os"
	"reflect"

	peer "gx/ipfs/QmfMmLGoKzCHDN7cGgk64PJr4iipzidDRME8HABSJqvmhC/go-libp2p-peer"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/corenet"
	"github.com/ipfs/go-ipfs/repo/config"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
)

const (
	// MyNodePath specifies the path where the node repository
	// (containing data and configuration) of the default node
	// used by the application is stored
	MyNodePath = "./data/MyNode"

	// nBitsForKeypair sets the strength of keypair
	nBitsForKeypair = 2048
)

// MyNode provides a global Node instance of user's own node
var MyNode *Node

func init() {
	var err error

	if !Exists(MyNodePath) {
		err = NewNodeRepo(MyNodePath, nil)
		if err != nil {
			panic(err)
		}
	}

	MyNode, err = NewNode(MyNodePath)
	if err != nil {
		panic(err)
	}
}

// Node provides an abstraction for IpfsNode and is the prefered way
// of accessing Nodes in our application. Note that IpfsNode is an
// embedded type.
type Node struct {
	*core.IpfsNode
	cancel context.CancelFunc
}

// NewNode creates a new Node from an existing node repository
func NewNode(path string) (*Node, error) {

	// Need to increse limit for number of filedescriptors to
	// avoid running out of those due to a lot of sockets
	err := checkAndSetUlimit()
	if err != nil {
		return nil, err
	}

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
		return nil, err
	}

	return &Node{
		IpfsNode: node,
		cancel:   cancel,
	}, nil
}

// Request is the generalized method to connect to another peer and
// send requests and receive responses. This is used by Client defined
// in peerapi.go and should not be used directly.
func (n *Node) Request(targetPeer string, path string, body interface{}, resp interface{}) error {

	// Check if Node hash is valid
	target, err := peer.IDB58Decode(targetPeer)
	if err != nil {
		return err
	}

	// Connect to targetPeer
	stream, err := corenet.Dial(n.IpfsNode, target, path)
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
	os.MkdirAll(repoRoot, 0755)

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

	fsrepo.Init(repoRoot, conf)
	if err != nil {
		return err
	}

	return initializeIpnsKeyspace(repoRoot)
}
