package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var FlagPullPostsFrom string
var FlagNoPeerServer bool
var FlagAddPeer string

func init() {
	// Parse flags
	silent := flag.Bool("silent", false, "Supresses all output except for stderr")
	pullPosts := flag.String("pullPost", "", "Supresses all output except for stderr")
	noPeer := flag.Bool("noPeer", false, "Supresses all output except for stderr")
	addPeer := flag.String("addPeer", "", "Add a peer to the list of known nodes")

	flag.Parse()

	// Set Flags (important, need to be after flag.Parse()
	FlagNoPeerServer = *noPeer
	FlagPullPostsFrom = *pullPosts
	FlagAddPeer = *addPeer

	// Initialize Logger
	if *silent {
		// If --silent is set log to a file instead
		file, err := os.OpenFile("data/log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("Failed to open log file")
		}
		LoggerInit(file, file, file, file)
	} else {
		LoggerInit(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	}
}

func main() {
	// Initialize Curation module
	err := MyCurator.Init()
	if err != nil {
		panic(err)
	}
	defer MyCurator.Close()

	// Starts PeerServer (non-blocking)
	if !FlagNoPeerServer {
		StartPeerAPI(MyNode)
	}

	if FlagPullPostsFrom != "" {
		target := FlagPullPostsFrom
		MyNode.pullPostFrom(target)
		Info.Println("Done pulling")

	}

	if FlagAddPeer != "" {
		peerID := FlagAddPeer
		MyNode.AddPeer(peerID)
	}

	// Starts communication pipeline for GUI
	StartGUIPipe(MyNode)

	// Discover new peers periodically
	ticker := time.NewTicker(30 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				MyNode.DiscoverPeers()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
