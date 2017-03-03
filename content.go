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
	Title     string
	Content   string
	Timestamp int64

	// @TODO these fields are still getting saved to IPFS
	IPFSData
	UserData PostUserData
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
	UserData CommentUserData
}

// NewPost constructs a new posts and adds it to the IPFS network
func NewPost(user *User, title, content string) (*IPFSObj, error) {
	data := Post{
		Alias:     user.Alias,
		Title:     title,
		Content:   content,
		Timestamp: time.Now().Unix(),
	}

	obj, err := NewIPFSObj(MyNode, user, data)
	if err != nil {
		return nil, err
	}
	data.Hash = obj.Hash

	err = AddHostingNode(obj.Hash, MyNode.Identity.Pretty())
	if err != nil {
		return nil, err
	}

	MyCurator.OnPostAdded(&data, true)

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
	data.Hash = obj.Hash

	err = AssociateCommentWithPost(obj.Hash, postID)
	if err != nil {
		return nil, err
	}

	MyCurator.OnCommentAdded(&data, true)

	return obj, nil
}

func GetPost(postID string) (*Post, error) {
	// @TODO maybe only need to get POST if we haven't retrieved it yet
	obj, err := GetIPFSObj(postID)
	if err != nil {
		return nil, err
	}

	post := &Post{}
	mapstructure.Decode(obj.Data, post)

	post.Hash = obj.Hash
	post.Key = obj.Key

	userData := GetPostUserData(postID)
	post.UserData = userData

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

	userData := GetCommentUserData(commentID)
	comment.UserData = userData

	err = AssociateCommentWithPost(obj.Hash, comment.Post)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func GetComments(postID string) ([]Comment, error) {
	commentHashes, err := GetPostComments(postID)
	if err != nil {
		return nil, err
	}

	// @TODO @PERFORMANCE can do this concurrently
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

func GetContentPosts() ([]Post, error) {
	postHashes := MyCurator.GetContent(make(map[string]interface{}))

	// @TODO @PERFORMANCE can do this concurrently
	posts := []Post{}
	for _, h := range postHashes {
		post, err := GetPost(h)
		if err != nil {
			Warning.Println("Could not retrieve post with id", h)
			continue
		}

		posts = append(posts, *post)
	}

	return posts, nil
}

func GetAllPosts() ([]Post, error) {
	postHashes, err := GetPosts()
	if err != nil {
		return nil, err
	}

	// @TODO @PERFORMANCE can do this concurrently
	var posts []Post
	for _, h := range postHashes {
		post, err := GetPost(h)
		if err != nil {
			Warning.Println("Could not retrieve post with id", h)
			continue
		}

		posts = append(posts, *post)
	}

	return posts, nil

}
