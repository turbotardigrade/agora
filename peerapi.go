package main

import (
	"fmt"
	"gx/ipfs/QmdXimY9QHaasZmw6hWojWnCJvfgxETjZQfg9g6ZrA9wMX/go-libp2p-net"
)

func StartPeerAPI() {
	peerAPI := NewPeerServer(MyNode)
	peerAPI.handleFunc("/comments", GetComments)
	peerAPI.handleFunc("/health", GetHealth)
}

type GetCommentsReq struct {
	Post string
}

type GetCommentsResp struct {
	Comments []string
}

func GetComments(stream net.Stream) {
	req := GetCommentsReq{}
	ReadJSON(stream, &req)

	// @TODO lookup comments for given req.Post
	comments := GetCommentsResp{[]string{"hash1", "hash2", req.Post}}

	WriteJSON(stream, comments)
}

type GetHealthResp struct {
	Status string
}

func GetHealth(stream net.Stream) {
	fmt.Printf("Connection from: %s\n", stream.Conn().RemotePeer())
	WriteJSON(stream, GetHealthResp{"OK"})
}
