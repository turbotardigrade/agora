package main

import (
	"os"
)

// HandleCmdIfCLI if one or more commands specified, run agora in CLI
// mode (means that agora terminates immidiately after execution
func HandleCmdIfCLI() {
	isCLI := false

	// Add Known peers
	if len(opts.AddPeers) != 0 {
		isCLI = true

		for _, p := range opts.AddPeers {
			Info.Println("Add peer with id", p)
			err := MyNode.AddPeer(p)
			if err != nil {
				panic(err)
			}
		}
	}

	// Delete all known peers
	if opts.DeletePeers {
		isCLI = true

		Info.Println("Remove all known peers")
		err := MyNode.RemoveAllPeers()
		if err != nil {
			panic(err)
		}
	}

	// Pull all posts from a given node
	if opts.PullPosts != "" {
		isCLI = true

		Info.Println("Pull all posts form peer", opts.PullPosts)
		MyNode.pullPostFrom(opts.PullPosts)
		Info.Println("Done pulling")
	}

	// Just terminate the program as we use it in CLI mode
	if isCLI {
		os.Exit(0)
	}
}
