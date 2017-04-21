package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/fatih/structs"

	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/core/coreunix"
	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/repo/config"
)

var (
	testNode1 *Node
	testUser1 *User
	testNode2 *Node
	testUser2 *User
)

func init() {
	fmt.Println("------------------------------------------------------------")
	fmt.Println("Initialize tests")
	fmt.Println("------------------------------------------------------------")

	// Init paths
	basePath := ExecutionPath + "/data/"
	testNode1Path := basePath + "TestNode1/"
	testNode2Path := basePath + "TestNode2/"

	// Remove testNodes they exists
	err := RemoveContents(testNode1Path)
	if err != nil {
		Warning.Println(err)
	}
	err = RemoveContents(testNode2Path)
	if err != nil {
		Warning.Println(err)
	}

	err = os.MkdirAll(basePath, 0755)
	if err != nil {
		Warning.Println(err)
	}

	// Initialize Curation module
	MyCurator.Init()

	// Create testNode1
	err = NewNodeRepo(testNode1Path, nil)
	if err != nil {
		panic(err)
	}

	testNode1, err = NewNode(testNode1Path)
	if err != nil {
		panic(err)
	}

	// This ensures that internals like GetIPFSObj uses the test
	// node to retrieve objects
	MyNode = testNode1

	testUser1, err = NewUser("tester1")
	if err != nil {
		panic(err)
	}

	// Create testNode2
	// Need to change Addresses in order to avoid clashes
	addr := &config.Addresses{
		Swarm: []string{
			"/ip4/0.0.0.0/tcp/4003",
			"/ip6/::/tcp/4003",
		},
		API:     "/ip4/127.0.0.1/tcp/5003",
		Gateway: "/ip4/127.0.0.1/tcp/8082",
	}

	err = NewNodeRepo(testNode2Path, addr)
	if err != nil {
		panic(err)
	}

	testNode2, err = NewNode(testNode2Path)
	if err != nil {
		panic(err)
	}

	testUser2, err = NewUser("tester2")
	if err != nil {
		panic(err)
	}

	// Start PeerAPIs
	StartPeerAPI(testNode2)

	// Might need to give some time for peerAPI info to propagate
	// through IPFS network
	// fmt.Println("Wait 5 sec to seed node information to network...")
	// time.Sleep(5 * time.Second)
	fmt.Println("\n------------------------------------------------------------")
	fmt.Println("Start tests")
	fmt.Println("------------------------------------------------------------")
}

func TestBlacklistThroughCuration(t *testing.T) {

	// Disable optimistically unchoke, otherwise test will fail
	DisableOptimisticallyUnchoke = true

	// Prepare fake post and fake peer
	peerID := "RANDOMPEERID"
	postID := "RANDOMPOSTID"
	err := testNode1.AddPeer(peerID)
	if err != nil {
		t.Error("Should be able to add new Peer")
	}

	err = testNode1.AddHostingNode(postID, peerID)
	if err != nil {
		t.Error("Should be able to add Hosting Node")
	}

	// Remove from blacklist just in case it's in there already
	err = testNode1.RemoveBlacklist(peerID)
	if err != nil {
		fmt.Println(err)
	}

	// Report the same spam 20 times (should only report it as one)
	for i := 0; i < 30; i++ {
		testNode1.onSpam(peerID, postID)
	}

	isBlacklist, err := testNode1.IsBlacklisted(peerID)
	if err != nil {
		t.Error("Should be able to check Blacklist")
	}

	if isBlacklist {
		t.Error("Peer should not be blacklisted, because only one posts has been reported ( but repeatedly)")
	}

	// Add unique spam elements, which should be still under the
	// blacklist threshold
	for i := 0; i < 3; i++ {
		testNode1.onSpam(peerID, string(i))
	}

	isBlacklist, err = testNode1.IsBlacklisted(peerID)
	if err != nil {
		t.Error("Should be able to check Blacklist")
	}

	if isBlacklist {
		t.Error("Peer should not be blacklisted, because only one posts has been reported ( but repeatedly)")
	}

	// Add more, now it should be blacklisted
	for i := 0; i < 55; i++ {
		testNode1.onSpam(peerID, string(i+10))
	}

	isBlacklist, err = testNode1.IsBlacklisted(peerID)
	if err != nil {
		t.Error("Should be able to check Blacklist")
	}

	if !isBlacklist {
		t.Error("Peer should be blacklisted by now")
	}

	// Check if peer was removed by knownHosts
	peers, err := testNode1.GetPeers()
	if err != nil {
		t.Error("GetPeers should not error")
	}

	if StringInSlice(peerID, peers) {
		t.Error("Should be removed")
	}

	// Check if peer connection gets rejected
	_, err = Client{testNode1}.CheckHealth(peerID)
	if err != ErrSkipBlacklisted {
		t.Error("Should skip this request, since node is blacklisted")
	}
}

func TestPostCommentCreationAndRetrival(t *testing.T) {
	fmt.Println("\n=== Try NewPost and GetPost")
	fmt.Println("Create new Post")

	postContent := "PostContent"
	postTitle := "PostTitle"

	obj, err := testNode2.NewPost(testUser1, postTitle, postContent)
	if err != nil {
		panic(err)
	}

	fmt.Println("Retrieve Post", obj.Hash)
	post, err := testNode2.GetPost(obj.Hash)
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

	obj, err = testNode2.NewComment(testUser1, post.Hash, post.Hash, commentContent)
	if err != nil {
		panic(err)
	}

	fmt.Println("Retrieve Comment", obj.Hash)
	comment, err := testNode2.GetComment(obj.Hash)
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

	comments, err := Client{testNode1}.GetComments(testNode2.ID, post.Hash)
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
		Alias:     testUser1.Alias,
		Title:     "TrueTitle",
		Content:   "TrueContent",
		Timestamp: Now(),
	}

	// Inserted tempered data
	obj := &IPFSObj{Key: testUser1.PubKeyRaw}
	obj.Data = structs.New(data).Map()

	obj.Signature, err = Sign(testUser1, obj.Data)
	if err != nil {
		panic(err)
	}

	// Temper with data
	obj.Data["Content"] = "riggedContent"

	// Add to IPFS Node Repository
	hash, err := coreunix.Add(testNode2.IpfsNode, ToJSONReader(obj))
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
	testNode2.pullPostFrom(testNode1.ID)

	params := make(map[string]interface{})
	fmt.Println("Curation suggested comments:")
	fmt.Println(MyCurator.GetContent(params))
}

func TestHealthAPI(t *testing.T) {
	fmt.Println("\n=== Try /health")

	isHealthy, err := Client{testNode1}.CheckHealth(testNode2.ID)
	if err != nil {
		panic(err)
	}

	if !isHealthy {
		t.Errorf(`Expected Health Status to be true`)
	}
}
