package main

// @TODO need to refactor
func (n *Node) pullPostFrom(target string) {
	Info.Println("Request Posts from Peer", target)
	postHashes, err := Client{n}.GetPosts(target)
	if err != nil {
		Warning.Println(err)
	} else {
		Info.Println("Received post Hashes:", postHashes)
	}

	for _, hash := range postHashes {
		postObj, err := n.GetPost(hash)
		if err != nil {
			Warning.Println("PullPosts", err)
			continue
		}

		isAccepted := MyCurator.OnPostAdded(postObj, false)
		if !isAccepted {
			Info.Println("Content got rejected. Hash:", postObj.Hash)
			continue
		}

		n.AddHostingNode(postObj.Hash, target)

		// Get Comments from node
		commentHashes, err := Client{n}.GetComments(target, postObj.Hash)

		for _, hash := range commentHashes {
			_, err := n.GetComments(hash)
			if err != nil {
				Warning.Println("PullPosts", err)
				continue
			} else {
				n.AssociateCommentWithPost(hash, postObj.Hash)
			}

		}
	}
}

// DiscoverPeers gets peers from all existing peers and adds them to the DB
func (n *Node) DiscoverPeers() (err error) {
	myPeers, err := n.GetPeers()
	var allPeers []string

	if err != nil {
		Warning.Println("discoverPeers", err)
	}

	for _, peerID := range myPeers {
		newPeers, err := Client{n}.GetPeers(peerID)
		if err != nil {
			Warning.Println("Error getting peers from "+peerID, err)
		}

		Info.Println("Recieved %d peers from %s", len(newPeers), peerID)

		for _, newPeerID := range newPeers {
			err := n.AddPeer(newPeerID)

			if err != nil {
				Warning.Println("Error adding peer to DB"+newPeerID, err)
			}
		}

	}

	return
}
