package main

// GetNodesHostingPost gives a list of known nodes who are known to
// seed the postID given
func GetNodesHostingPost(postID string) ([]string, error) {
	// @TODO check if postID is a valid hash
	return GetHostingNodes(postID)
}

// AddNodeHostingPost adds the nodeID to the database for the given postID
// @TODO Validate postID and nodeID
func AddNodeHostingPost(postID string, nodeID string) error {
	return AddHostingNode(postID, nodeID)
}

// @TODO need to refactor
func pullPostFrom(target string) {
	Info.Println("Request Posts from Peer", target)
	postHashes, err := Client{MyNode}.GetPosts(target)
	if err != nil {
		Warning.Println(err)
	} else {
		Info.Println("Received post Hashes:", postHashes)
	}

	for _, hash := range postHashes {
		postObj, err := GetPost(hash)
		if err != nil {
			Warning.Println("PullPosts", err)
			continue
		}

		isAccepted := MyCurator.OnPostAdded(postObj)
		if !isAccepted {
			Info.Println("Content got rejected. Hash:", postObj.Hash)
			continue
		}

		AddHostingNode(postObj.Hash, target)

		// Get Comments from node
		commentHashes, err := Client{MyNode}.GetComments(target, postObj.Hash)

		for _, hash := range commentHashes {
			_, err := GetComments(hash)
			if err != nil {
				Warning.Println("PullPosts", err)
				continue
			} else {
				AssociateCommentWithPost(hash, postObj.Hash)
			}

		}
	}
}
