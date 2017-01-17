package main

var KnownNodes = []string{
	"QmPXamy2Qe3AwhgWtfm6G7TaSiwf1WS9Zg1UQKgsq71ug2",
}

func GetNodesHostingPost(postID string) ([]string, error) {
	// @TODO check if postID is a valid hash
	// @TODO query database
	return KnownNodes, nil
}
