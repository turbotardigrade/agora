package main

import (
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

// Globally available variables
var db *Model
var dbPath = "./data/data.db"

// Model is a wrapper for the DB connection, mainly to create a safe
// namespace
type Model struct {
	*bolt.DB
}

const (
	postCommentsBucket = "posts2comment"
	postHostersBucket  = "post2hoster"
	postBucket         = "postBucket"
	blacklistBucket    = "blacklistBucket"
	knownNodesBucket   = "knownNodesBucket"
)

// bucketNames is used to iterate through all buckets e.g. on
// initialization
var bucketNames = []string{
	postCommentsBucket,
	postHostersBucket,
	postBucket,
	blacklistBucket,
	knownNodesBucket,
}

//////////////////////////////////////////////////////////////////////
/// Open and Close

// OpenDb opens bolt database and provides a db Model instance
// globally
func OpenDb() error {
	Info.Println("Init DB")

	config := &bolt.Options{Timeout: 2 * time.Second}
	dbInstance, err := bolt.Open(dbPath, 0600, config)
	if err != nil {
		// No point of running without DB, just kill the
		// application
		Error.Println("FATAL", err)
		log.Fatal(err)
	}

	db = &Model{dbInstance}

	// Create Bucket if they don't exists
	return dbInstance.Update(func(tx *bolt.Tx) error {
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
/// Model methods to access and manipulate data

func (m *Model) GetPosts() ([]string, error) {
	return BoltGetKeys(m.DB, postHostersBucket)
}

func (m *Model) GetHostingNodes(postID string) (nodes []string, err error) {
	return BoltGetList(m.DB, postHostersBucket, postID)
}

func (m *Model) GetPostComments(postID string) (nodes []string, err error) {
	return BoltGetList(m.DB, postCommentsBucket, postID)
}

func (m *Model) AddHostingNode(postID, nodeID string) error {
	return BoltAppendList(m.DB, postHostersBucket, postID, nodeID)
}

func (m *Model) AssociateCommentWithPost(comment, post string) error {
	return BoltAppendList(m.DB, postCommentsBucket, post, comment)
}

type PostUserData struct {
	Score   int
	Flagged bool
}

func (m *Model) GetPostUserData(hash string) (res PostUserData) {
	BoltGet(m.DB, postBucket, hash, &res)
	return res
}

func (m *Model) SetPostUserData(hash string, data PostUserData) error {
	return BoltSet(m.DB, postBucket, hash, data)
}

type CommentUserData struct {
	Score   int
	Flagged bool
}

func (m *Model) GetCommentUserData(hash string) (res CommentUserData) {
	BoltGet(m.DB, postBucket, hash, &res)
	return res
}

func (m *Model) SetCommentUserData(hash string, data CommentUserData) error {
	return BoltSet(m.DB, postBucket, hash, data)
}

func (m *Model) GetPeers() ([]string, error) {
	return BoltGetKeys(m.DB, knownNodesBucket)
}
