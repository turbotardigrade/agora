package main

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
)

// GetNodesHostingPost gives a list of known nodes who are known to
// seed the postID given
func GetNodesHostingPost(postID string) ([]string, error) {
	var knownNodes []string
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("knownNodes"))
		data := bucket.Get([]byte(postID))
		if err := json.Unmarshal(data, &knownNodes); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Could not get Post ID %s\n", postID)
		return nil, err
	}
	// @TODO check if postID is a valid hash
	return knownNodes, nil
}
