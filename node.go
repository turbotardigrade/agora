package main

import (
	"context"
	"fmt"
	"reflect"

	peer "gx/ipfs/QmfMmLGoKzCHDN7cGgk64PJr4iipzidDRME8HABSJqvmhC/go-libp2p-peer"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/corenet"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
)

// MyNode provides a global Node instance of user's own node
var MyNode *Node

// @TODO create Node locally in execution folder, do not use the one
// in the home folder
func init() {
	var err error
	MyNode, err = NewNode("~/.ipfs2")
	if err != nil {
		panic(err) // @TODO handle this gracefully
	}
}

type Node struct {
	ipfsNode *core.IpfsNode
	cancel   context.CancelFunc
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
		ipfsNode: node,
		cancel:   cancel,
	}, nil
}

func (n *Node) Request(targetPeer string, path string, body interface{}, resp interface{}) error {
	target, err := peer.IDB58Decode(targetPeer)
	if err != nil {
		return err
	}

	stream, err := corenet.Dial(n.ipfsNode, target, path)
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
