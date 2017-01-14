package main

import (
	"fmt"
	"testing"
	"time"
)

func TestPeerAPI(t *testing.T) {
	node, err := NewNode("~/.ipfs")
	if err != nil {
		panic(err)
	}

	StartPeerAPI()

	// Might need to give some time for peerAPI info to propagate
	// through IPFS network
	time.Sleep(5 * time.Second)

	fmt.Println("\nTry /comments")
	buf, err := node.Post("QmYHZAqj9Y1D2agM29BwtHPx5CvWKLwtfETkhb73ZCRoqs", GetCommentsReq{"1"}, "/comments")
	if err != nil {
		panic(err)
	}
	fmt.Println("resp: ", buf)

	fmt.Println("\nTry /health")
	buf, err = node.Get("QmYHZAqj9Y1D2agM29BwtHPx5CvWKLwtfETkhb73ZCRoqs", "/health")
	if err != nil {
		panic(err)
	}
	fmt.Println("resp: ", buf)

	if buf != `{"Status":"OK"}` {
		t.Errorf(`Expected Health to be {"Status":"OK"} and not %s`, buf)
	}
}
