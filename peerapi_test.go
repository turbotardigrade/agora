package main

import (
	"fmt"
	"testing"
	"time"
)

var TestNode *Node

func init() {
	fmt.Println("------------------------------------------------------------")
	fmt.Println("Initialize tests")
	fmt.Println("------------------------------------------------------------")

	// Create TestNode
	var err error
	TestNode, err = NewNode("~/.ipfs")
	if err != nil {
		panic(err)
	}

	// Start PeerAPIs
	StartPeerAPI(TestNode)

	// Might need to give some time for peerAPI info to propagate
	// through IPFS network
	fmt.Println("Wait 10 sec to seed node information to network...")
	time.Sleep(5 * time.Second)
	fmt.Println("\n------------------------------------------------------------")
	fmt.Println("Start tests")
	fmt.Println("------------------------------------------------------------")
}

func TestGetComments(t *testing.T) {
	fmt.Println(GetComments("1"))
}

func TestCommentsAPI(t *testing.T) {
	fmt.Println("\nTry /comments")

	targetPeer := TestNode.ipfsNode.Identity.Pretty()
	comments, err := Client{MyNode}.GetComments(targetPeer, "1")
	if err != nil {
		panic(err)
	}

	fmt.Println("resp: ", comments)
}

func TestHealthAPI(t *testing.T) {
	fmt.Println("\nTry /health")

	targetPeer := TestNode.ipfsNode.Identity.Pretty()
	isHealthy, err := Client{MyNode}.CheckHealth(targetPeer)
	if err != nil {
		panic(err)
	}

	if !isHealthy {
		t.Errorf(`Expected Health Status to be true`)
	}
}
