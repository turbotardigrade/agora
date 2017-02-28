package main

import (
	"fmt"
	"testing"
)

var testUser = &User{
	Alias:   "longh",
	PubKey:  "PubKey test",
	PrivKey: "PrivKey test",
}

func TestNewPostThenGetPost(t *testing.T) {
	fmt.Println("\nTry NewPost and GetPost")
	fmt.Println("Create new Post:")
	obj, err := NewPost(testUser, "Title", "Content")
	if err != nil {
		panic(err)
	}

	fmt.Println(obj)

	fmt.Println("Retrieve Post with from ", obj.Hash)
	post, err := GetPost(obj.Hash)
	if err != nil {
		panic(err)
	}

	fmt.Println(post)
}

func TestNewComment(t *testing.T) {
	fmt.Println("\nTry NewComment and GetComment")
	fmt.Println("Create new Comment:")

	postID := "Test Post Hash"
	obj, err := NewComment(testUser, postID, postID, "Hello World")
	if err != nil {
		panic(err)
	}

	fmt.Println(obj)

	fmt.Println("Retrieve Comment with from ", obj.Hash)
	comment, err := GetComment(obj.Hash)
	if err != nil {
		panic(err)
	}

	fmt.Println(comment)
}
