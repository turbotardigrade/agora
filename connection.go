package main

// DiscoverPeers gets peers from all existing peers and adds them to the DB
func (n *Node) DiscoverPeers() error {
	Info.Println("Start discovery...")

	myPeers, err := n.GetSomePeers()
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
			// ignore self
			if newPeerID == n.ID {
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

// @TODO need to refactor
func (n *Node) pullPostFrom(target string) {
	Info.Println("Request Posts from Peer", target)
	postHashes, err := Client{n}.GetPosts(target)
	if err != nil {
		Warning.Println(err)
		return
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
			Info.Println("Reject", postObj.Hash, "from", target)
			n.onSpam(target, postObj.Hash)
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
	unique := make(map[string]string)
	for _, host := range hosters {
		cmtHashes, err := Client{n}.GetComments(host, postHash)
		if err != nil {
			Warning.Println("Could not obtain comments for post", postHash, "from", host, err)
			continue
		}

		for _, h := range cmtHashes {
			unique[h] = host
		}
	}

	// Retrieve content and pass to curation
	for cmtHash, host := range unique {
		cmt, err := n.GetComment(cmtHash)
		if err != nil {
			Warning.Println("GetComment failed for", cmtHash)
			continue
		}

		isAccepted := MyCurator.OnCommentAdded(cmt, false)
		if !isAccepted {
			Info.Println("Comment", cmtHash, "classified as spam")
			n.onSpam(host, cmtHash)
			continue
		}

		n.AssociateCommentWithPost(cmtHash, postHash)
	}

	return nil
}

func (n *Node) onSpam(peer, contentHash string) {
	// Increment and get current spam counter
	spamCount, err := n.TrackSpam(peer, contentHash)
	if err != nil {
		Error.Println("onSpam failed due to", err)
		return
	}

	postCount, err := n.GetNumberOfPostsReceivedFromPeer(peer)
	if err != nil {
		Error.Println("onSpam failed due to", err)
		return
	}

	////////////////////////////////////////
	// Blacklist Conditions

	// Minimum spam threshold (need at least that many to be
	// considered for blacklist)
	if spamCount < 20 {
		return
	}

	// Ratio threshold until spammer gets dumped
	spamRatio := float32(spamCount) / float32(postCount)
	if spamRatio > 0.5 { // equivalent to if 33 % of all posts are spam
		Info.Println("Blacklist peer", peer, "due to bad spam ratio", spamRatio)
		n.AddBlacklist(peer)
	}

	peers, err := n.GetPeers()
	if err != nil {
		Error.Println("Unable to get peers:", err)
		return
	}

	if len(peers) >= 3 {
		return
	}

	Info.Println("Number of peers low, try to optimistically unchoke")

	candidate, err := n.GetLeastSpamBlacklisted()
	if err != nil {
		Error.Println("Unable to get candidate to unchoke:", err)
		return
	}

	err = n.RemoveBlacklist(candidate)
	if err != nil {
		Error.Println("Cannot remove blacklist")
	}

	err = n.AddPeer(candidate)
	if err != nil {
		Error.Println("Cannot add peer to peer list")
	}

	Info.Println("Optimistically unchoke previously blacklisted node", candidate)
}
