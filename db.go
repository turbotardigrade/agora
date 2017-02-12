package main

import (
	"encoding/json"
	"fmt"
	"log"
	"path"
	"runtime"
	"time"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

const (
	dbPath             = "data/data.db"
	postCommentsBucket = "posts2comment"
	postHostersBucket  = "post2hoster"
)

// Used to iterate through all buckets e.g. on initialization
var bucketNames = []string{
	postCommentsBucket,
	postHostersBucket,
}

func GetHostingNodes(postID string) (nodes []string, err error) {
	return BoltGetList(postHostersBucket, postID)
}

func GetPostComments(postID string) (nodes []string, err error) {
	return BoltGetList(postCommentsBucket, postID)
}

func BoltGetList(bucketName, key string) ([]string, error) {
	var list []string
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(B(bucketName))
		data := bucket.Get(B(key))

		if len(data) == 0 {
			// results in empty list
			return nil
		}

		if err := json.Unmarshal(data, &list); err != nil {
			return err
		}

		return nil
	})

	return list, err
}

func BoltAppendList(bucketName, key, elem string) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(B(bucketName))

		// If the postID already has an entry, unmarshal and append the nodeID to it
		// else create a new array containing the nodeID
		var list []string
		if data := bucket.Get(B(key)); data != nil {
			if err := json.Unmarshal(data, &list); err != nil {
				return err
			}
			list = append(list, elem)
		} else {
			list = []string{elem}
		}

		data, _ := json.Marshal(list)
		return bucket.Put(B(key), data)
	})
}

func AddHostingNode(postID, nodeID string) error {
	return BoltAppendList(postHostersBucket, postID, nodeID)
}

func AssociateCommentWithPost(comment, post string) error {
	return BoltAppendList(postCommentsBucket, post, comment)
}

// OpenDb opens bolt database
func OpenDb() error {
	Info.Println("Init DB")

	_, filename, _, _ := runtime.Caller(0) // get full path of this file
	dbfile := path.Join(path.Dir(filename), dbPath)
	config := &bolt.Options{Timeout: 1 * time.Second}

	var err error
	db, err = bolt.Open(dbfile, 0644, config)
	if err != nil {
		log.Fatal(err)
	}

	// Create Buckets if they don't exists
	return db.Update(func(tx *bolt.Tx) error {
		for _, bucketName := range bucketNames {
			_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
		}
		return nil
	})
}

// CloseDb closes bolt database
func CloseDb() {
	db.Close()
}
