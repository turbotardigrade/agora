package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"

	peer "gx/ipfs/QmfMmLGoKzCHDN7cGgk64PJr4iipzidDRME8HABSJqvmhC/go-libp2p-peer"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/corenet"
	"github.com/ipfs/go-ipfs/namesys"
	"github.com/ipfs/go-ipfs/repo/config"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
)

const (
	MyNodePath      = "./data/MyNode"
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
		panic(err) // @TODO handle this gracefully
	}
}

type Node struct {
	*core.IpfsNode
	cancel context.CancelFunc
}

func NewNode(path string) (*Node, error) {
	r, err := fsrepo.Open(path)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	cfg := &core.BuildCfg{
		Repo:   r,
		Online: true,
	}

	node, err := core.NewNode(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Node{
		IpfsNode: node,
		cancel:   cancel,
	}, nil
}

func (n *Node) Request(targetPeer string, path string, body interface{}, resp interface{}) error {
	target, err := peer.IDB58Decode(targetPeer)
	if err != nil {
		return err
	}

	stream, err := corenet.Dial(n.IpfsNode, target, path)
	if err != nil {
		return err
	}

	if reflect.ValueOf(resp).Kind() != reflect.Ptr {
		fmt.Println("WARNING: You must pass resp by &reference and not by value. This is not done for a request to", targetPeer, path)
	}

	WriteJSON(stream, &body)
	ReadJSON(stream, &resp)

	return nil
}

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

// Taken from github.com/ipfs/go-ipfs/blob/master/cmd/ipfs/init.go
func initializeIpnsKeyspace(repoRoot string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r, err := fsrepo.Open(repoRoot)
	if err != nil { // NB: repo is owned by the node
		return err
	}

	nd, err := core.NewNode(ctx, &core.BuildCfg{Repo: r})
	if err != nil {
		return err
	}
	defer nd.Close()

	err = nd.SetupOfflineRouting()
	if err != nil {
		return err
	}

	return namesys.InitializeKeyspace(ctx, nd.DAG, nd.Namesys, nd.Pinning, nd.PrivateKey)
}
