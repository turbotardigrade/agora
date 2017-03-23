package main

import (
	"fmt"
	"testing"

	"github.com/fatih/structs"

	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/core/coreunix"
	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/repo/config"
)

const testNodePath = "./data/TestNode/"

var (
	testNode *Node
	testUser *User
)

func init() {
	fmt.Println("------------------------------------------------------------")
	fmt.Println("Initialize tests")
	fmt.Println("------------------------------------------------------------")

	// Remove testNode if it exists
	err := RemoveContents(testNodePath)
	if err != nil {
		Warning.Println(err)
	}

	// Initialize Curation module
	MyCurator.Init()

	// Create testNode
	// Need to change Addresses in order to avoid clashes with MyNode
	addr := &config.Addresses{
		Swarm: []string{
			"/ip4/0.0.0.0/tcp/4002",
			"/ip6/::/tcp/4002",
		},
		API:     "/ip4/127.0.0.1/tcp/5002",
		Gateway: "/ip4/127.0.0.1/tcp/8081",
	}

	err = NewNodeRepo(testNodePath, addr)
	if err != nil {
		panic(err)
	}

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
	// fmt.Println("Wait 5 sec to seed node information to network...")
	// time.Sleep(5 * time.Second)
	fmt.Println("\n------------------------------------------------------------")
	fmt.Println("Start tests")
	fmt.Println("------------------------------------------------------------")
}

func TestPostCommentCreationAndRetrival(t *testing.T) {
	fmt.Println("\n=== Try NewPost and GetPost")
	fmt.Println("Create new Post")

	postContent := "PostContent"
	postTitle := "PostTitle"

	obj, err := MyNode.NewPost(testUser, postTitle, postContent)
	if err != nil {
		panic(err)
	}

	fmt.Println("Retrieve Post", obj.Hash)
	post, err := MyNode.GetPost(obj.Hash)
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

	obj, err = MyNode.NewComment(testUser, post.Hash, post.Hash, commentContent)
	if err != nil {
		panic(err)
	}

	fmt.Println("Retrieve Comment", obj.Hash)
	comment, err := MyNode.GetComment(obj.Hash)
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

	comments, err := Client{testNode}.GetComments(MyNode.ID, post.Hash)
	if err != nil {
		panic(err)
	}

	if !StringInSlice(comment.Hash, comments) {
		t.Errorf(`Expected to retrieve comment %s via /comments`, comment.Hash)
	}

	fmt.Println("\n/comments resp: ", comments)
}

func TestSignatureVerification(t *testing.T) {
	fmt.Println("\n=== Try Post rigged content")
	var err error

	data := Post{
		Alias:     MyUser.Alias,
		Title:     "TrueTitle",
		Content:   "TrueContent",
		Timestamp: Now(),
	}

	// Inserted tempered data
	obj := &IPFSObj{Key: MyUser.PubKeyRaw}
	obj.Data = structs.New(data).Map()

	obj.Signature, err = Sign(MyUser, obj.Data)
	if err != nil {
		panic(err)
	}

	// Temper with data
	obj.Data["Content"] = "riggedContent"

	// Add to IPFS Node Repository
	hash, err := coreunix.Add(MyNode.IpfsNode, ToJSONReader(obj))
	if err != nil {
		panic(err)
	}

	obj, err = GetIPFSObj(hash)
	if err != RiggedError {
		t.Errorf("Expected to detect tampered data and throw RiggedError")
	}
}

func TestGetPostsAPI(t *testing.T) {
	fmt.Println("\n=== Try pullPost")
	MyNode.pullPostFrom(testNode.ID)

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
