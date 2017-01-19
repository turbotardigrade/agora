package main

// KnownNodes used temporarily to store nodes, later should use a db
var KnownNodes = []string{
	"QmPXamy2Qe3AwhgWtfm6G7TaSiwf1WS9Zg1UQKgsq71ug2",
}

// GetNodesHostingPost gives a list of known nodes who are known to
// seed the postID given
func GetNodesHostingPost(postID string) ([]string, error) {
	// @TODO check if postID is a valid hash
	// @TODO query database
	return KnownNodes, nil
}
