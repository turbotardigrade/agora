package main

import (
	"time"

	"github.com/mitchellh/mapstructure"
)

/** @TODOs:
- Create cryptographic signature with author's private key
- Verify signature with author key
- Limit content size to e.g. 1kb
*/

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

// NewPost constructs a new posts and adds it to the IPFS network
func NewPost(user *User, content string) (*IPFSObj, error) {
	data := Post{
		Alias:     user.Alias,
		Content:   content,
		Timestamp: time.Now().Unix(),
	}

	obj, err := NewIPFSObj(MyNode, user, data)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// NewComment constructs a new comment and adds it to the IPFS
// network. Note that a valid hash of the post this comment is
// submitted to needs to be provided. If this comment is a reply to
// any other comment, include the parent comments in ancestors.
func NewComment(user *User, postID, content string, ancestors []string) (*IPFSObj, error) {
	// @TODO check if postID and parent are valid

	data := Comment{
		Post:      postID,
		Ancestors: ancestors,
		Alias:     user.Alias,
		Content:   content,
		Timestamp: time.Now().Unix(),
	}

	obj, err := NewIPFSObj(MyNode, user, data)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func GetPost(postID string) (*Post, error) {
	obj, err := GetIPFSObj(postID)
	if err != nil {
		return nil, err
	}

	post := &Post{}
	mapstructure.Decode(obj.Data, post)

	return post, nil
}

func GetComment(commentID string) (*Comment, error) {
	obj, err := GetIPFSObj(commentID)
	if err != nil {
		return nil, err
	}

	comment := &Comment{}
	mapstructure.Decode(obj.Data, comment)

	return comment, nil
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
			Error.Println(err)
			continue
		}

		// @TODO Add comment hash to database if not added yet
		result = append(result, comments...)
	}

	// @TODO retrieve comments from IPFS
	return result, nil
}
