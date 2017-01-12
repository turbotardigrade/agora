package main

import (
	"fmt"
	"io"
	"os"

	peer "gx/ipfs/QmfMmLGoKzCHDN7cGgk64PJr4iipzidDRME8HABSJqvmhC/go-libp2p-peer"

	corenet "github.com/ipfs/go-ipfs/core/corenet"
)

func get(fromPeerID string, path string) {
	target, err := peer.IDB58Decode(fromPeerID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("I am peer %s dialing %s\n", MyNode.Identity, target)

	con, err := corenet.Dial(MyNode, target, path)
	if err != nil {
		fmt.Println(err)
		return
	}

	io.Copy(os.Stdout, con)
}
