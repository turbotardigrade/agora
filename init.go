package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

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

const LogDirPath = "./data"

func init() {
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatalln("Failed to parse args", err)
	}

	// Initialize Logger
	if opts.Silent {
		// If --silent is set log to a file instead
		logPath := LogDirPath + "/log.txt"
		err := os.MkdirAll(LogDirPath, 0755)
		if err != nil {
			log.Fatalln("Failed to create log directory", err)
		}

		err = CreateFileIfNotExists(logPath)
		if err != nil {
			log.Fatalln("Failed to create file", err)
		}

		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("Failed to open log file", err)
		}

		file.WriteString("----------------------------------------\n")
		file.WriteString("   Log of " + time.Now().Format("Jan _2 15:04") + "\n")
		file.WriteString("----------------------------------------\n")
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
