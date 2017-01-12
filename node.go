package main

import (
	"context"
	"fmt"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
)

var MyNode *core.IpfsNode

func init() {
	fmt.Println("INIT")

	// Basic ipfsnode setup
	// @TODO take this from Config files
	r, err := fsrepo.Open("~/.ipfs2")
	if err != nil {
		panic(err)
	}

	ctx, _ := context.WithCancel(context.Background())
	//defer cancel()

	cfg := &core.BuildCfg{
		Repo:   r,
		Online: true,
	}

	MyNode, err = core.NewNode(ctx, cfg)
	if err != nil {
		panic(err) // @TODO handle error more gracefully
	}
}
