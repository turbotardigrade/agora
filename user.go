package main

// @TODO Use Constructor and Getter when interacting with User struct

var MyUser = &User{
	PubKey:  "TODO",
	PrivKey: "TODO",
	Alias:   "Long",
}

type User struct {
	PubKey  string
	PrivKey string
	Alias   string
}
