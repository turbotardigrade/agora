package main

import (
	"errors"
	"time"

	"github.com/mitchellh/mapstructure"
)

/** @TODOs:
- Create cryptographic signature with author's private key
- Verify signature with author key
- Limit content size to e.g. 1kb
*/

// IPFSData is be embedded to Post and Comment and only filled when
// loaded from IPFS. It extends Post and Comment with information like
// hash, which is not known beforehand. Do not mistake with IPFSObj,
// which is the general abstraction for blobs in IPFS.
type IPFSData struct {
	Hash string
	Key  string
}

// Post defines the data structure used by our application to handle
// posts and also provides the model for database
type Post struct {
	// Alias is authors display name
	Alias     string
	Content   string
	Timestamp int64

	IPFSData
}

// Comment defines the data structure used by our application to
// handle comments and also provides the model for database
type Comment struct {
	// Post refers to the posts this comment is submitted to
	Post string
	// Parent refers to the item (can be post or comment) to which
	// this comment is replying to
	Parent string

	// Alias is authors display name
	Alias     string
	Content   string
	Timestamp int64

	IPFSData
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

	err = AddHostingNode(obj.Hash, MyNode.Identity.Pretty())
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// NewComment constructs a new comment and adds it to the IPFS
// network. Note that parent == post for comments replying to posts
func NewComment(user *User, postID, parent, content string) (*IPFSObj, error) {
	if parent == "" || postID == "" {
		// @TODO check if postID and parent are valid
		return nil, errors.New("Parent and/or Post not defined")
	}

	if content == "" {
		return nil, errors.New("Content cannot be empty")
	}

	data := Comment{
		Post:      postID,
		Parent:    parent,
		Alias:     user.Alias,
		Content:   content,
		Timestamp: time.Now().Unix(),
	}

	obj, err := NewIPFSObj(MyNode, user, data)
	if err != nil {
		return nil, err
	}

	err = AssociateCommentWithPost(obj.Hash, postID)
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

	post.Hash = obj.Hash
	post.Key = obj.Key

	return post, nil
}

func GetComment(commentID string) (*Comment, error) {
	obj, err := GetIPFSObj(commentID)
	if err != nil {
		return nil, err
	}

	comment := &Comment{}
	mapstructure.Decode(obj.Data, comment)

	comment.Hash = obj.Hash
	comment.Key = obj.Key

	return comment, nil
}

// @TODO for now just return the hash of the string but later should return []Comments
func GetComments(postID string) ([]Comment, error) {
	commentHashes, err := GetPostComments(postID)
	if err != nil {
		return nil, err
	}

	// @TODO add comments from other nodes as well

	var comments []Comment
	for _, h := range commentHashes {
		comment, err := GetComment(h)
		if err != nil {
			Warning.Println("Could not retrieve comment with id", h)
			continue
		}

		comments = append(comments, *comment)
	}

	return comments, nil
}
