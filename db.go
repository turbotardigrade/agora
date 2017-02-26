package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

var db *bolt.DB
var dbPath = "./data/data.db"

const (
	postCommentsBucket = "posts2comment"
	postHostersBucket  = "post2hoster"
	postBucket         = "postBucket"
	commentBucket      = "commentBucket"
	blacklistBucket    = "blacklistBucket"
)

// Used to iterate through all buckets e.g. on initialization
var bucketNames = []string{
	postCommentsBucket,
	postHostersBucket,
	postBucket,
	commentBucket,
	blacklistBucket,
}

// @TODO this function looks awful, should be somehow refactored
func GetPosts() ([]string, error) {
	pMap := make(map[string]struct{})

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(postHostersBucket))
		b.ForEach(func(k, _ []byte) error {
			pMap[string(k)] = struct{}{}
			return nil
		})
		return nil
	})
	if err != nil {
		return []string{}, err
	}

	posts := make([]string, len(pMap))
	i := 0
	for k, _ := range pMap {
		posts[i] = k
		i++
	}

	return posts, nil
}

func GetHostingNodes(postID string) (nodes []string, err error) {
	return BoltGetList(postHostersBucket, postID)
}

func GetPostComments(postID string) (nodes []string, err error) {
	return BoltGetList(postCommentsBucket, postID)
}

func AddHostingNode(postID, nodeID string) error {
	return BoltAppendList(postHostersBucket, postID, nodeID)
}

func AssociateCommentWithPost(comment, post string) error {
	return BoltAppendList(postCommentsBucket, post, comment)
}

type PostUserData struct {
	Score   int
	Flagged bool
}

func GetPostUserData(hash string) (res PostUserData) {
	BoltGet(postBucket, hash, &res)
	return res
}

func SetPostUserData(hash string, data PostUserData) error {
	return BoltSet(postBucket, hash, data)
}

type CommentUserData struct {
	Score   int
	Flagged bool
}

func GetCommentUserData(hash string) (res CommentUserData) {
	BoltGet(postBucket, hash, &res)
	return res
}

func SetCommentUserData(hash string, data CommentUserData) error {
	return BoltSet(postBucket, hash, data)
}

// OpenDb opens bolt database
func OpenDb() error {
	Info.Println("Init DB")

	config := &bolt.Options{Timeout: 2 * time.Second}

	var err error
	db, err = bolt.Open(dbPath, 0600, config)
	if err != nil {
		Error.Println("FATAL", err)
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

//////////////////////////////////////////////////////////////////////
// Helper functions to deal with boltdb

// BoltGetList will get a specific list of string (as array) from a
// boltdb bucket
func BoltGetList(bucketName, key string) ([]string, error) {
	var list []string
	err := BoltGet(bucketName, key, &list)
	return list, err
}

func BoltGet(bucketName, key string, ptr interface{}) error {
	return db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(B(bucketName))
		data := bucket.Get(B(key))

		if len(data) == 0 {
			// results in no changes to ptr
			return nil
		}

		if err := json.Unmarshal(data, ptr); err != nil {
			return err
		}

		return nil
	})
}

func BoltSet(bucketName, key string, obj interface{}) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(B(bucketName))
		data, _ := json.Marshal(obj)
		return bucket.Put(B(key), data)
	})
}

// BoltAppendList appends a string to a list
func BoltAppendList(bucketName, key, elem string) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(B(bucketName))

		// If already has an entry, unmarshal and append to it,
		// else create a new array containing the element
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
