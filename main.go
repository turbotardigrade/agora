package main

import (
	"time"
)

func main() {
	var err error
	MyNode, err = CreateNodeIfNotExists(MyNodePath)
	if err != nil {
		panic(err)
	}

	MyUser, err = CreateUserIfNotExists(MyUserConfPath, "DefaultBob")
	if err != nil {
		panic(err)
	}

	// Checks if agora is running in CLI mode and executes cmds
	// accordingly
	HandleCmdIfCLI()

	// Initialize Curation module
	err = MyCurator.Init()
	if err != nil {
		panic(err)
	}
	defer MyCurator.Close()

	// Starts PeerServer (non-blocking)
	if opts.NoPeer {
		Info.Println("PeerServer that provides PeerAPI disabled")
	} else {
		StartPeerAPI(MyNode)
	}

	// Starts peer list monitor (non-blocking)
	if opts.MonPeers {
		Info.Println("Peer monitor enabled")
		go monitorPeers(MyNode)
	}

	// Discover new peers periodically
	if opts.NoDiscover {
		Info.Println("Discovery disabled")
	} else {
		go func() {
			for {
				MyNode.DiscoverPeers()
				time.Sleep(20 * time.Second)
			}
		}()
	}

	if opts.NoPull {
		Info.Println("Content pull disabled")
	} else {
		time.Sleep(5 * time.Second)
		go func() {
			for {
				time.Sleep(3 * time.Second)
				Info.Println("Pull new content from known network")
				peers, err := MyNode.GetSomePeers()
				if err != nil {
					Error.Println("Unable to get list of known peers", err)
					continue
				}

				for _, p := range peers {
					Info.Println("Pull content from", p)
					MyNode.pullPostFrom(p)
				}

				if !opts.NoComments {
					// @TODO better to not check all at once
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
			}
		}()
	}

	// Starts communication pipeline for GUI
	StartGUIPipe(MyNode)
}
