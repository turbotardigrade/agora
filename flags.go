package main

import (
	"log"
	"os"

	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	Silent     bool `short:"s" long:"silent" description:"Supresses outputs except for stderr"`
	NoPeer     bool `long:"noPeer" description:"Disable peerAPI server"`
	NoPull     bool `long:"noPull" description:"Do not pull content from other peers"`
	NoDiscover bool `long:"noDiscover" description:"Disable discovery"`
	MonPeers   bool `long:"monPeers" description:"Monitor list of peers"`
	NoComments bool `long:"noComments" description:"Disable automatic content pull of comments"`

	// CLIs
	AddPeers    []string `short:"a" long:"addPeer" description:"Add peer to known nodes"`
	PullPosts   string   `long:"pullPosts" description:"Pulls all posts from remote node"`
	DeletePeers bool     `long:"deleteAllPeers" description:"Deletes all known peers"`
	Initt       bool     `long:"init" description:"Creates ipfs repo if not exists"`

	Curator string `short:"c" long:"curator" description:"Specify the curation module used. Use 'none' to load dummy curator"`
}

func InitFlags() {
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatalln("Failed to parse args", err)
	}
}
