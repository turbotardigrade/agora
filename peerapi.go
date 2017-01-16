package main

import "gx/ipfs/QmdXimY9QHaasZmw6hWojWnCJvfgxETjZQfg9g6ZrA9wMX/go-libp2p-net"

// StartPeerAPI starts PeerServer and register PeerAPI handlers which
// will asynchronously listen for incoming requests
func StartPeerAPI(node *Node) {
	peerAPI := NewPeerServer(node)
	peerAPI.handleFunc("/comments", GetCommentsHandler)
	peerAPI.handleFunc("/health", GetHealthHandler)
}

type Client struct {
	Node *Node
}

type GetCommentsReq struct {
	Post string
}

type GetCommentsResp struct {
	Comments []string
}

func GetCommentsHandler(stream net.Stream) {
	req := GetCommentsReq{}
	ReadJSON(stream, &req)

	// @TODO lookup comments for given req.Post
	comments := GetCommentsResp{[]string{"hash1", "hash2", req.Post}}

	WriteJSON(stream, comments)
}

func (c Client) GetComments(target, postID string) ([]string, error) {
	var req = GetCommentsReq{postID}
	var resp GetCommentsResp

	err := c.Node.Request(target, "/comments", req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Comments, nil
}

type GetHealthResp struct {
	Status string
}

func GetHealthHandler(stream net.Stream) {
	WriteJSON(stream, GetHealthResp{"OK"})
}

func (c Client) CheckHealth(target string) (bool, error) {
	var resp GetHealthResp
	err := c.Node.Request(target, "/health", nil, &resp)
	if err != nil {
		return false, err
	}

	return resp.Status == "OK", nil
}
