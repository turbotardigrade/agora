package main

import (
	"fmt"
	"testing"
	"time"
)

func TestPeerAPI(t *testing.T) {

	// Create TestNode
	node, err := NewNode("~/.ipfs")
	if err != nil {
		panic(err)
	}

	// Start PeerAPI server and get it's nodeID
	StartPeerAPI(MyNode)
	targetPeer := MyNode.ipfsNode.Identity.Pretty()

	// Might need to give some time for peerAPI info to propagate
	// through IPFS network
	time.Sleep(5 * time.Second)

	//////////////////////////////
	// Testing health API
	fmt.Println("\nTry /health")

	var healthResp GetHealthResp

	err = node.Request(targetPeer, "/health", nil, &healthResp)
	if err != nil {
		panic(err)
	}
	fmt.Println("resp: ", healthResp)

	if healthResp.Status != "OK" {
		t.Errorf(`Expected Health Status to be OK and not '%s'`, healthResp.Status)
	}

	//////////////////////////////
	// Testing comments API
	fmt.Println("\nTry /comments")

	comments, err := Cl{node}.GetComments(targetPeer, "1")
	if err != nil {
		panic(err)
	}

	fmt.Println("resp: ", comments)
}
