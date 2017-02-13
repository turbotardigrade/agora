package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var FlagPullPostsFrom string
var FlagNoPeerServer bool

func init() {
	// Parse flags
	silent := flag.Bool("silent", false, "Supresses all output except for stderr")
	pullPosts := flag.String("pullPost", "", "Supresses all output except for stderr")
	noPeer := flag.Bool("noPeer", false, "Supresses all output except for stderr")

	flag.Parse()

	// Set Flags (important, need to be after flag.Parse()
	FlagNoPeerServer = *noPeer
	FlagPullPostsFrom = *pullPosts

	// Initialize Logger
	if *silent {
		LoggerInit(ioutil.Discard, ioutil.Discard, ioutil.Discard, os.Stderr)
	} else {
		LoggerInit(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	}
}

func main() {
	OpenDb()
	defer CloseDb()

	// Starts PeerServer (non-blocking)
	if !FlagNoPeerServer {
		StartPeerAPI(MyNode)
	}

	if FlagPullPostsFrom != "" {
		Info.Println("Request Posts from Peer", FlagPullPostsFrom)
		posts, err := Client{MyNode}.GetPosts(FlagPullPostsFrom)
		if err != nil {
			Warning.Println(err)
		} else {
			fmt.Println(posts)
		}
	}

	// Starts communication pipeline for GUI
	StartGUIPipe()
}
