package main

import (
	"flag"
	"io/ioutil"
	"log"
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
	OpenDb()
	defer CloseDb()

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
		pullPostFrom(target)
		Info.Println("Done pulling")

	}

	// Starts communication pipeline for GUI
	StartGUIPipe()
}
