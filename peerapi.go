package main

import (
	"gx/ipfs/QmRuZnMorqodado1yeTQiv1i9rmtKj29CjPSsBKM7DFXV4/go-libp2p-net"
	"time"
)

// StartPeerAPI starts PeerServer and register PeerAPI handlers which
// will asynchronously listen for incoming requests
func StartPeerAPI(node *Node) {
	peerAPI := NewPeerServer(node)
	peerAPI.HandleFunc("/comments", GetCommentsHandler)
	peerAPI.HandleFunc("/posts", GetPostsHandler)
	peerAPI.HandleFunc("/health", GetHealthHandler)
	peerAPI.HandleFunc("/peers", GetPeersHandler)
	peerAPI.HandleFunc("/blacklist", GetBlacklistHandler)

	Info.Println("Seed for 5 seconds...\n")
	time.Sleep(5 * time.Second)
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
func GetCommentsHandler(n *Node, stream net.Stream) {
	req := GetCommentsReq{}
	ReadJSON(stream, &req)

	// @TODO lookup comments for given req.Post
	commentHashes, err := n.GetPostComments(req.Post)
	if err != nil {
		Warning.Println(err)
	}

	comments := GetCommentsResp{commentHashes}

	WriteJSON(stream, comments)
}

// GetComments provides helper function for the client to query the
// Comments API
func (c Client) GetComments(target, postID string) ([]string, error) {
	var req = GetCommentsReq{postID}
	var resp GetCommentsResp

	err := c.Request(target, "/comments", req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Comments, nil
}

// ----------------------------------------
// /posts API - temporary

// GetPostsResp defines response
type GetPostsResp struct {
	Posts []string
}

// GetPostsHandler provides stream handler
func GetPostsHandler(n *Node, stream net.Stream) {
	posts, err := n.GetPosts()
	if err != nil {
		Warning.Println("Error retrieving list of posts: ", err)
		return
	}

	resp := GetPostsResp{posts}

	WriteJSON(stream, resp)
}

// GetPosts provides helper function for the client to query the
// Comments API
func (c Client) GetPosts(target string) ([]string, error) {
	var resp GetPostsResp
	err := c.Request(target, "/posts", nil, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Posts, nil
}

// ----------------------------------------
// /health API

// GetHealthResp defines response
type GetHealthResp struct {
	Status string
}

// GetHealthHandler provides stream handler
func GetHealthHandler(n *Node, stream net.Stream) {
	WriteJSON(stream, GetHealthResp{"OK"})
}

// CheckHealth provides helper function for the client to query the
// Health API
func (c Client) CheckHealth(target string) (bool, error) {
	var resp GetHealthResp
	err := c.Request(target, "/health", nil, &resp)
	if err != nil {
		return false, err
	}

	return resp.Status == "OK", nil
}

// ----------------------------------------
// /peers API

// GetPeersResp defines response
type GetPeersResp struct {
	Peers []string
}

// GetPeersHandler provides stream handler
func GetPeersHandler(n *Node, stream net.Stream) {
	peers, err := n.GetPeers()
	if err != nil {
		Warning.Println("Error retrieving list of peers: ", err)
		return
	}
	WriteJSON(stream, GetPeersResp{Peers: peers})
}

// ----------------------------------------
// /blacklist API

// GetBlacklistResp defines response
type GetBlacklistResp struct {
	Peers []map[string]float32
}

// GetBlacklistHandler provides stream handler
func GetBlacklistHandler(n *Node, stream net.Stream) {
	counts, err := n.GetBlacklist()
	if err != nil {
		Warning.Println("Error retrieving list of peers: ", err)
		return
	}

	peers := make([]map[string]float32, len(counts))
	i := 0
	for k, v := range counts {
		postCount, err := n.GetNumberOfPostsReceivedFromPeer(k)
		if err != nil {
			Error.Println("Failed to get post received count", err)
			postCount = 1
		}

		peers[i] = make(map[string]float32)
		peers[i][k] = float32(v) / float32(postCount)
		i += 1
	}

	WriteJSON(stream, GetBlacklistResp{Peers: peers})
}

// GetPeers provides helper function query peers of a node
func (c Client) GetPeers(target string) ([]string, error) {
	var resp GetPeersResp
	err := c.Request(target, "/peers", nil, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Peers, nil
}
