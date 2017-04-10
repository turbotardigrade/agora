package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var FlagPullPostsFrom string
var FlagNoPeerServer bool
var FlagAddPeer string
var FlagMonitorPeers bool
var FlagCurator string
var FlagNoDiscovery bool
var FlagNoContentPull bool

func init() {
	// Parse flags
	silent := flag.Bool("silent", false, "Supresses all output except for stderr")
	pullPosts := flag.String("pullPost", "", "Pull all posts from given peer")
	noPeer := flag.Bool("noPeer", false, "Disable peerAPI server")
	addPeer := flag.String("addPeer", "", "Add a peer to the list of known nodes")
	monPeers := flag.Bool("monPeers", false, "Monitor list of peers")
	curator := flag.String("curator", "", "Specify the curation module used. Use 'none' to load dummy curator")
	noDiscovery := flag.Bool("noDiscovery", false, "Disable discovery")
	noContentPull := flag.Bool("noPull", false, "Do not pull content from other peers")

	flag.Parse()

	// Set Flags (important, need to be after flag.Parse()
	FlagNoPeerServer = *noPeer
	FlagPullPostsFrom = *pullPosts
	FlagAddPeer = *addPeer
	FlagMonitorPeers = *monPeers
	FlagNoDiscovery = *noDiscovery
	FlagNoContentPull = *noContentPull

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

	// Set Curator module
	if strings.ToLower(*curator) == "none" {
		Info.Println("Using DummyCurator")
		MyCurator = &DummyCurator{}
	} else {
		Info.Println("Using MLCurator")
		MyCurator = &MLCurator{}
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

	// Starts peer list monitor (non-blocking)
	if FlagMonitorPeers {
		go monitorPeers(MyNode)
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

	// Discover new peers periodically
	if !FlagNoDiscovery {
		ticker := time.NewTicker(5 * time.Second)
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

	if !FlagNoContentPull {
		go func() {
			for {
				time.Sleep(1 * time.Second)
				Info.Println("Pull new content from known network")
				peers, err := MyNode.GetSomePeers()
				if err != nil {
					Error.Println("Unable to get list of known peers", err)
					continue
				}

				// @TODO select random n nodes to connect to

				for _, p := range peers {
					Info.Println("Pull content from", p)
					MyNode.pullPostFrom(p)
				}

				posts, err := MyNode.GetPosts()
				if err != nil {
					Error.Println("Unable to get list of posts", err)
					continue
				}

				for _, p := range posts {
					err := MyNode.pullPostComments(p)
					if err != nil {
						Warning.Println("Error getting comments for post", p, err)
						continue
					}
				}

			}
		}()
	}

	// Starts communication pipeline for GUI
	StartGUIPipe(MyNode)
}
