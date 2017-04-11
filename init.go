package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	Silent     bool `short:"s" long:"silent" description:"Supresses outputs except for stderr"`
	NoPeer     bool `long:"noPeer" description:"Disable peerAPI server"`
	NoPull     bool `long:"noPull" description:"Do not pull content from other peers"`
	NoDiscover bool `long:"noDiscover" description:"Disable discovery"`
	MonPeers   bool `long:"monPeers" description:"Monitor list of peers"`

	// CLIs
	AddPeers  []string `short:"a" long:"addPeer" description:"Add peer to known nodes"`
	PullPosts string   `long:"pullPosts" description:"Pulls all posts from remote node"`

	Curator string `short:"c" long:"curator" description:"Specify the curation module used. Use 'none' to load dummy curator"`
}

func init() {
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Println(err)
		panic(err)
		os.Exit(1)
	}

	// Initialize Logger
	if opts.Silent {
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
	if strings.ToLower(opts.Curator) == "none" {
		Info.Println("Using DummyCurator")
		MyCurator = &DummyCurator{}
	} else {
		Info.Println("Using MLCurator")
		MyCurator = &MLCurator{}
	}
}
