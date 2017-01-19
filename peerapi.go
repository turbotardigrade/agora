package main

import "gx/ipfs/QmdXimY9QHaasZmw6hWojWnCJvfgxETjZQfg9g6ZrA9wMX/go-libp2p-net"

// StartPeerAPI starts PeerServer and register PeerAPI handlers which
// will asynchronously listen for incoming requests
func StartPeerAPI(node *Node) {
	peerAPI := NewPeerServer(node)
	peerAPI.HandleFunc("/comments", GetCommentsHandler)
	peerAPI.HandleFunc("/health", GetHealthHandler)
}

// Client should be used to send requests to the PeerAPI
// Example usage:
// 	isHealthy, err := Client{MyNode}.CheckHealth(targetPeer)
//	if err != nil {
//	     panic(err)
//	}
type Client struct {
	Node *Node
}

// ----------------------------------------
// /comments API

// GetCommentsReq defines request body
type GetCommentsReq struct {
	Post string
}

// GetCommentsResp defines response
type GetCommentsResp struct {
	Comments []string
}

// GetCommentsHandler provides stream handler
func GetCommentsHandler(stream net.Stream) {
	req := GetCommentsReq{}
	ReadJSON(stream, &req)

	// @TODO lookup comments for given req.Post
	comments := GetCommentsResp{[]string{"hash1", "hash2", req.Post}}

	WriteJSON(stream, comments)
}

// GetComments provides helper function for the client to query the
// Comments API
func (c Client) GetComments(target, postID string) ([]string, error) {
	var req = GetCommentsReq{postID}
	var resp GetCommentsResp

	err := c.Node.Request(target, "/comments", req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Comments, nil
}

// ----------------------------------------
// /health API

// GetHealthResp defines response
type GetHealthResp struct {
	Status string
}

// GetHealthHandler provides stream handler
func GetHealthHandler(stream net.Stream) {
	WriteJSON(stream, GetHealthResp{"OK"})
}

// CheckHealth provides helper function for the client to query the
// Health API
func (c Client) CheckHealth(target string) (bool, error) {
	var resp GetHealthResp
	err := c.Node.Request(target, "/health", nil, &resp)
	if err != nil {
		return false, err
	}

	return resp.Status == "OK", nil
}
