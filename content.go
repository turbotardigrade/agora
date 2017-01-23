package main

import (
	"fmt"
	"time"

	"github.com/ipfs/go-ipfs/core/coreunix"
)

/** @TODOs:
- Create cryptographic signature with author's private key
- Verify signature with author key
- Limit content size to e.g. 1kb
*/

// IPFSObj is an abstraction to deal with objects / blobs from IPFS;
// also does signing and verification
type IPFSObj struct {
	// Hash is the hash address given by IPFS
	Hash string
	// Key is the PublicKey of the signing user
	Key string
	// Data is an opaque field to dump the payload in
	Data interface{}
	// Signature used to verify this object was indeed sent by
	// user with Key
	Signature string
}

// Post defines the data structure used by our application to handle
// posts and also provides the model for database
type Post struct {
	// Alias is authors display name
	Alias     string
	Content   string
	Timestamp int64
}

// Comment defines the data structure used by our application to
// handle comments and also provides the model for database
type Comment struct {
	// Post refers to the posts this comment is submitted to
	Post string
	// Ancestors is the list of hashes of upper level comments,
	// which this comment replies to
	Ancestors []string

	// Alias is authors display name
	Alias     string
	Content   string
	Timestamp int64
}

// NewIPFSObj is a generalized helper function to created signed
// IPFSObj and add it to the IPFS network
func NewIPFSObj(node *Node, key string, data interface{}) (*IPFSObj, error) {
	obj := &IPFSObj{Key: key, Data: data}

	// @TODO make cryptographic signature with given key and data
	obj.Signature = "TODO"

	// Add to IPFS Node Repository
	hash, err := coreunix.Add(node.IpfsNode, ToJSONReader(obj))
	if err != nil {
		return nil, err
	}
	obj.Hash = hash

	return obj, nil
}

// Verify checks if the data of IPFSObj is valid by checking the
// signature with provided PublicKey of the author
func Verify(obj IPFSObj) (bool, error) {
	// @TODO
	return true, nil
}

// NewPost constructs a new posts and adds it to the IPFS network
func NewPost(user User, content string) (*IPFSObj, error) {
	data := Post{
		Alias:     user.Alias,
		Content:   content,
		Timestamp: time.Now().Unix(),
	}

	obj, err := NewIPFSObj(MyNode, user.Key, data)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// NewComment constructs a new comment and adds it to the IPFS
// network. Note that a valid hash of the post this comment is
// submitted to needs to be provided. If this comment is a reply to
// any other comment, include the parent comments in ancestors.
func NewComment(user User, postID, content string, ancestors []string) (*IPFSObj, error) {
	// @TODO check if postID and parent are valid

	data := Comment{
		Post:      postID,
		Ancestors: ancestors,
		Alias:     user.Alias,
		Content:   content,
		Timestamp: time.Now().Unix(),
	}

	obj, err := NewIPFSObj(MyNode, user.Key, data)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// @TODO for now just return the hash of the string but later should return []Comments
func GetComments(postID string) ([]string, error) {
	hosts, err := GetNodesHostingPost(postID)
	if err != nil {
		return nil, err
	}

	client := Client{MyNode}
	result := []string{}

	for _, target := range hosts {
		comments, err := client.GetComments(target, "1")
		if err != nil { // @TODO handle errors
			fmt.Println(err)
			continue
		}

		// @TODO Add comment hash to database if not added yet
		result = append(result, comments...)
	}

	// @TODO retrieve comments from IPFS
	return result, nil
}
