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
