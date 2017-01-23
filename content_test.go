package main

import (
	"fmt"
	"testing"
)

var testUser = User{
	Alias: "longh",
	Key:   "test",
}

func TestNewPostThenGetPost(t *testing.T) {
	fmt.Println("\nTry NewPost and GetPost")
	fmt.Println("Create new Post:")
	obj, err := NewPost(testUser, "Hello World")
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
	ancestors := []string{"Test Ancestor Hash 1", "Test Ancestor Hash 2"}
	obj, err := NewComment(testUser, postID, "Hello World", ancestors)
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