package main

import (
	"encoding/json"
	"fmt"
	"gx/ipfs/QmdXimY9QHaasZmw6hWojWnCJvfgxETjZQfg9g6ZrA9wMX/go-libp2p-net"
)

type GetCommentsReq struct {
	Post string
}

type GetCommentsResp struct {
	Comments []string
}

func main() {
	peerAPI := NewPeerServer(MyNode)
	peerAPI.handleFunc("/comments", func(stream net.Stream) {
		var req = GetCommentsReq{}
		json.NewDecoder(stream).Decode(&req)

		stream.Write([]byte(req.Post))
		fmt.Printf("Connection from: %s\n", stream.Conn().RemotePeer())
	})

	AnotherNode, err := NewNode("~/.ipfs")
	if err != nil {
		panic(err)
	}

	peerAPI2 := NewPeerServer(AnotherNode)
	peerAPI2.handleFunc("/health", func(stream net.Stream) {
		fmt.Fprintln(stream, "Other")
		fmt.Printf("Connection from: %s\n", stream.Conn().RemotePeer())
	})

	buf, err := AnotherNode.Post("QmYHZAqj9Y1D2agM29BwtHPx5CvWKLwtfETkhb73ZCRoqs", GetCommentsReq{"1"}, "/comments")
	if err != nil {
		panic(err)
	}
	fmt.Println("Buffer Me: ", buf)

	buf2, err := MyNode.Get("QmPXamy2Qe3AwhgWtfm6G7TaSiwf1WS9Zg1UQKgsq71ug2", "/health")
	if err != nil {
		panic(err)
	}
	fmt.Println("Buffer Other: ", buf2)
}
