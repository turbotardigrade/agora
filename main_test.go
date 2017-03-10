package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/repo/config"
)

const testNodePath = "./data/TestNode"

var (
	testNode *Node
	testUser *User
)

func init() {
	fmt.Println("------------------------------------------------------------")
	fmt.Println("Initialize tests")
	fmt.Println("------------------------------------------------------------")

	// Overwrite dbPath to use test database instead
	dbPath = "data/testdata.db"

	// Remove existing database if it exists
	os.Remove(dbPath)

	// Open connection to database
	OpenDb()

	// Initialize Curation module
	MyCurator.Init()

	// Create testNode if not exists
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

	testUser, err = NewUser("tester")
	if err != nil {
		panic(err)
	}

	// Start PeerAPIs
	StartPeerAPI(MyNode)

	// Might need to give some time for peerAPI info to propagate
	// through IPFS network
	fmt.Println("Wait 5 sec to seed node information to network...")
	time.Sleep(5 * time.Second)
	fmt.Println("\n------------------------------------------------------------")
	fmt.Println("Start tests")
	fmt.Println("------------------------------------------------------------")
}

func TestPostCommentCreationAndRetrival(t *testing.T) {
	fmt.Println("\n=== Try NewPost and GetPost")
	fmt.Println("Create new Post")

	postContent := "PostContent"
	postTitle := "PostTitle"

	obj, err := NewPost(testUser, postTitle, postContent)
	if err != nil {
		panic(err)
	}

	fmt.Println("Retrieve Post", obj.Hash)
	post, err := GetPost(obj.Hash)
	if err != nil {
		panic(err)
	}

	fmt.Println("")
	PrettyPrint(post)

	if post.Content != postContent || post.Title != postTitle {
		t.Errorf(`Expected posted post and retrieved post to be the same`)
	}

	fmt.Println("\n=== Try NewComment and GetComment")
	fmt.Println("Create new Comment")

	commentContent := "CommentContent"

	obj, err = NewComment(testUser, post.Hash, post.Hash, commentContent)
	if err != nil {
		panic(err)
	}

	fmt.Println("Retrieve Comment", obj.Hash)
	comment, err := GetComment(obj.Hash)
	if err != nil {
		panic(err)
	}

	fmt.Println("")
	PrettyPrint(comment)

	if comment.Content != commentContent {
		t.Errorf(`Expected posted comment and retrieved comment to be the same`)
	}

	if comment.Parent != post.Hash || comment.Post != post.Hash {
		t.Errorf(`Expected posted comment parent and post to be %s`, post.Hash)
	}

	fmt.Println("\n=== Try /comments")

	comments, err := Client{testNode}.GetComments(MyNode.Identity.Pretty(), post.Hash)
	if err != nil {
		panic(err)
	}

	if !StringInSlice(comment.Hash, comments) {
		t.Errorf(`Expected to retrieve comment %s via /comments`, comment.Hash)
	}

	fmt.Println("\n/comments resp: ", comments)
}

func TestGetPostsAPI(t *testing.T) {
	fmt.Println("\n=== Try pullPost")
	pullPostFrom(testNode.Identity.Pretty())

	params := make(map[string]interface{})
	fmt.Println("Curation suggested comments:")
	fmt.Println(MyCurator.GetContent(params))
}

func TestHealthAPI(t *testing.T) {
	fmt.Println("\n=== Try /health")

	isHealthy, err := Client{testNode}.CheckHealth(MyNode.Identity.Pretty())
	if err != nil {
		panic(err)
	}

	if !isHealthy {
		t.Errorf(`Expected Health Status to be true`)
	}
}
