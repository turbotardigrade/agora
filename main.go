package main

import (
	"flag"
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
		target := FlagPullPostsFrom
		Info.Println("Request Posts from Peer", target)
		postHashes, err := Client{MyNode}.GetPosts(target)
		if err != nil {
			Warning.Println(err)
		} else {
			Info.Println("Received post Hashes:", postHashes)
		}

		for _, hash := range postHashes {
			postObj, err := GetPost(hash)
			if err != nil {
				Warning.Println("PullPosts", err)
				continue
			} else {
				AddHostingNode(postObj.Hash, target)
			}

			// Get Comments from node
			commentHashes, err := Client{MyNode}.GetComments(target, postObj.Hash)

			for _, hash := range commentHashes {
				_, err := GetComments(hash)
				if err != nil {
					Warning.Println("PullPosts", err)
					continue
				} else {
					AssociateCommentWithPost(hash, postObj.Hash)
				}

			}
		}

		Info.Println("Done pulling")

	}

	// Starts communication pipeline for GUI
	StartGUIPipe()
}
