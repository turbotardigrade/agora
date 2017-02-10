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

// AddNodeHostingPost adds the nodeID to the database for the given postID
// @TODO Validate postID and nodeID
func AddNodeHostingPost(postID string, nodeID string) error {
	var knownNodes []string
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("knownNodes"))
		// If the postID already has an entry, unmarshal and append the nodeID to it
		// else create a new array containing the nodeID
		if data := bucket.Get([]byte(postID)); data != nil {
			if err := json.Unmarshal(data, &knownNodes); err != nil {
				return err
			}
			knownNodes = append(knownNodes, nodeID)
		} else {
			knownNodes = []string{nodeID}
		}
		data, _ := json.Marshal(knownNodes)
		err := bucket.Put([]byte(postID), data)
		return err
	})
	if err != nil {
		fmt.Printf("Could not add Node ID %s to Post ID %s\n", nodeID, postID)
	}
	return err
}
