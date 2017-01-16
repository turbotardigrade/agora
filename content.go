package main

import "fmt"

/** @TODOs:
- Constructor for Comment
- Create cryptographic signature with author's private key
- Verify signature with author key
*/

type Comment struct {
	PostID string
	Parent string

	AuthorKey   string
	AuthorAlias string
	Content     string
	Signature   string

	Timestamp int64
}

// @TODO for now just return the hash of the string but later should return []Comments
func GetComments(postID string) ([]string, error) {
	hosts, err := GetNodesHostingPost(postID)
	if err != nil {
		return nil, err
	}

	client := Client{MyNode}
	result := []string{}

	for _, target := range hosts {
		comments, err := client.GetComments(target, "1")
		if err != nil { // @TODO handle errors
			fmt.Println(err)
			continue
		}

		// @TODO Add comment hash to database if not added yet
		result = append(result, comments...)
	}

	// @TODO retrieve comments from IPFS
	return result, nil
}
