package main

import (
	"errors"

	"github.com/mitchellh/mapstructure"
)

/** @TODOs:
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
	Alias   string
	Title   string
	Content string

	// Note that we are using string and not
	// int64, because int64 will be converted to
	// floats when marshalling from interface{}
	Timestamp string

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
	Alias   string
	Content string

	// Note that we are using string and not
	// int64, because int64 will be converted to
	// floats when marshalling from interface{}
	Timestamp string

	IPFSData
	UserData CommentUserData
}

// NewPost constructs a new posts and adds it to the IPFS network
func (n *Node) NewPost(user *User, title, content string) (*IPFSObj, error) {
	data := Post{
		Alias:     user.Alias,
		Title:     title,
		Content:   content,
		Timestamp: Now(),
	}

	obj, err := NewIPFSObj(n, user, data)
	if err != nil {
		return nil, err
	}
	data.Hash = obj.Hash

	err = n.AddHostingNode(obj.Hash, n.ID)
	if err != nil {
		return nil, err
	}

	MyCurator.OnPostAdded(&data, true)

	return obj, nil
}

// NewComment constructs a new comment and adds it to the IPFS
// network. Note that parent == post for comments replying to posts
func (n *Node) NewComment(user *User, postID, parent, content string) (*IPFSObj, error) {
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
		Timestamp: Now(),
	}

	obj, err := NewIPFSObj(MyNode, user, data)
	if err != nil {
		return nil, err
	}
	data.Hash = obj.Hash

	isAccepted := MyCurator.OnCommentAdded(&data, true)
	if !isAccepted {
		return nil, errors.New("Curation rejected the content")
	}

	err = n.AssociateCommentWithPost(obj.Hash, postID)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (n *Node) GetPost(postID string) (*Post, error) {
	// @TODO maybe only need to get POST if we haven't retrieved it yet
	obj, err := GetIPFSObj(postID)
	if err != nil {
		return nil, err
	}

	post := &Post{}
	mapstructure.Decode(obj.Data, post)

	post.Hash = obj.Hash
	post.Key = obj.Key

	userData := n.GetPostUserData(postID)
	post.UserData = userData

	return post, nil
}

func (n *Node) GetComment(commentID string) (*Comment, error) {
	obj, err := GetIPFSObj(commentID)
	if err != nil {
		return nil, err
	}

	comment := &Comment{}
	mapstructure.Decode(obj.Data, comment)

	comment.Hash = obj.Hash
	comment.Key = obj.Key

	userData := n.GetCommentUserData(commentID)
	comment.UserData = userData

	return comment, nil
}

func (n *Node) GetComments(postID string) ([]Comment, error) {
	commentHashes, err := n.GetPostComments(postID)
	if err != nil {
		return nil, err
	}

	// @TODO @PERFORMANCE can do this concurrently
	var comments []Comment
	for _, h := range commentHashes {
		comment, err := n.GetComment(h)
		if err != nil {
			Warning.Println("Could not retrieve comment with id", h)
			continue
		}

		comments = append(comments, *comment)
	}

	return comments, nil
}

func (n *Node) GetContentPosts() ([]Post, error) {
	postHashes := MyCurator.GetContent(make(map[string]interface{}))

	// @TODO @PERFORMANCE can do this concurrently
	posts := []Post{}
	for _, h := range postHashes {
		post, err := n.GetPost(h)
		if err != nil {
			Warning.Println("Could not retrieve post with id", h)
			continue
		}

		posts = append(posts, *post)
	}

	return posts, nil
}

func (n *Node) GetAllPosts() ([]Post, error) {
	postHashes, err := n.GetPosts()
	if err != nil {
		return nil, err
	}

	// @TODO @PERFORMANCE can do this concurrently
	var posts []Post
	for _, h := range postHashes {
		post, err := n.GetPost(h)
		if err != nil {
			Warning.Println("Could not retrieve post with id", h)
			continue
		}

		posts = append(posts, *post)
	}

	return posts, nil
}
