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
		// @TODO check if already exists

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
	}
}

func (n *Node) pullPostComments(postHash string) error {
	// Get all hosters
	hosters, err := n.GetHostingNodes(postHash)
	if err != nil {
		return err
	}

	// Get Comment hashes from hosting nodes
	unique := make(map[string]bool)
	for _, host := range hosters {
		cmtHashes, err := Client{n}.GetComments(host, postHash)
		if err != nil {
			Warning.Println("Could not obtain comments for post", postHash, "from", host, err)
			continue
		}

		for _, h := range cmtHashes {
			unique[h] = true
		}
	}

	// Retrieve content and pass to curation
	for cmtHash := range unique {
		cmt, err := n.GetComment(cmtHash)
		if err != nil {
			Warning.Println("GetComment failed for", cmtHash)
			continue
		}

		isAccepted := MyCurator.OnCommentAdded(cmt, false)
		if !isAccepted {
			Info.Println("Comment", cmtHash, "classified as spam")
			continue
		}

		n.AssociateCommentWithPost(cmtHash, postHash)
	}

	return nil
}

// DiscoverPeers gets peers from all existing peers and adds them to the DB
func (n *Node) DiscoverPeers() error {
	Info.Println("Start discovery...")

	myPeers, err := n.GetPeers()
	if err != nil {
		Warning.Println("discoverPeers", err)
		return err
	}

	for _, peerID := range myPeers {
		newPeers, err := Client{n}.GetPeers(peerID)
		if err != nil {
			Warning.Println("Error getting peers from", peerID, err)
			continue
		}

		Info.Println("Received", len(newPeers), "peers from", peerID)

		for _, newPeerID := range newPeers {
			if newPeerID == n.ID {
				// ignore self
				continue
			}

			err := n.AddPeer(newPeerID)
			if err != nil {
				Warning.Println("Error adding peer", newPeerID, "to DB", err)
			}
		}

	}

	return nil
}

func (n *Node) PullPostsFromPeers() error {
	Info.Println("Start content retrieval...")

	myPeers, err := n.GetPeers()
	if err != nil {
		Warning.Println("discoverPeers", err)
		return err
	}

	for _, peerID := range myPeers {
		n.pullPostFrom(peerID)
	}

	return nil
}
