package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/ipfs/go-ipfs/repo/config"
)

const testNodePath = "./data/TestNode"

var testNode *Node

func init() {
	fmt.Println("------------------------------------------------------------")
	fmt.Println("Initialize tests")
	fmt.Println("------------------------------------------------------------")

	// Create testNode
	if !Exists(testNodePath) {
		// Need to change Addresses in order to avoid clashes with MyNode
		addr := &config.Addresses{
			Swarm: []string{
				"/ip4/0.0.0.0/tcp/4002",
				"/ip6/::/tcp/4002",
			},
			API:     "/ip4/127.0.0.1/tcp/5002",
			Gateway: "/ip4/127.0.0.1/tcp/8081",
		}

		err := NewNodeRepo(testNodePath, addr)
		if err != nil {
			panic(err)
		}
	}

	var err error
	testNode, err = NewNode(testNodePath)
	if err != nil {
		panic(err)
	}

	// Start PeerAPIs
	StartPeerAPI(testNode)

	// Might need to give some time for peerAPI info to propagate
	// through IPFS network
	fmt.Println("Wait 5 sec to seed node information to network...")
	time.Sleep(5 * time.Second)
	fmt.Println("\n------------------------------------------------------------")
	fmt.Println("Start tests")
	fmt.Println("------------------------------------------------------------")
}

func TestGetComments(t *testing.T) {
	// TODO this is just for the time being
	// knownNodes = []string{testNode.Identity.Pretty()}
	fmt.Println(GetComments("1"))
}

func TestCommentsAPI(t *testing.T) {
	fmt.Println("\nTry /comments")

	targetPeer := testNode.Identity.Pretty()
	comments, err := Client{MyNode}.GetComments(targetPeer, "1")
	if err != nil {
		panic(err)
	}

	fmt.Println("resp: ", comments)
}

func TestHealthAPI(t *testing.T) {
	fmt.Println("\nTry /health")

	targetPeer := testNode.Identity.Pretty()
	isHealthy, err := Client{MyNode}.CheckHealth(targetPeer)
	if err != nil {
		panic(err)
	}

	if !isHealthy {
		t.Errorf(`Expected Health Status to be true`)
	}
}
