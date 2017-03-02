package main

import (
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

var MyCurator = DefaultCurator{}

type Curator interface {
	// Init will be called on initialization, use this function to
	// initialize your curation module
	Init() error

	// OnPostAdded will be called when new posts are retrieved
	// from other peers, if this functions returns false, the
	// content will be rejected (e.g. in the case of spam) and not
	// stored by our node
	OnPostAdded(obj *Post) bool

	// OnCommentAdded will be called when new comments are
	// retrieved from other peers, if this functions returns
	// false, the content will be rejected (e.g. in the case of
	// spam) and not stored by our node
	OnCommentAdded(obj *Comment) bool

	// GetContent gives back an ordered array of post hashes of
	// suggested content by the curation module
	GetContent(params map[string]interface{})

	// FlagContent marks or unmarks hashes as spam. True means
	// content is flagged as spam, false means content is not
	// flagged as spam
	FlagContent(hash string, isFlagged bool)

	// UpvoteContent is called when user upvotes a content
	UpvoteContent(hash string)

	// DownvoteContent is called when user downvotes a content
	DownvoteContent(hash string)
}

// DefaultCurator is a placeholder curation module. Currently it
// registers all content, but does no filtering nor ranking
type DefaultCurator struct{}

var curationDB *bolt.DB

const curationDBPath = "./data/curation.db"
const curationBucket = "exampleCuration"

// Init initializes boltdb which simply keeps track of saved hashes
// and their arrivaltime
func (c *DefaultCurator) Init() error {
	Info.Println("Init Curation DB")
	config := &bolt.Options{Timeout: 2 * time.Second}

	var err error
	curationDB, err = bolt.Open(curationDBPath, 0600, config)
	if err != nil {
		Error.Println("FATAL", err)
		log.Fatal(err)
	}

	// Create Buckets if they don't exists
	return curationDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(curationBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func (c *DefaultCurator) OnPostAdded(obj *Post) bool {
	err := curationDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(B(curationBucket))
		timestamp := string(time.Now().Unix())

		log.Println(timestamp)

		return bucket.Put(B(obj.Hash), B(timestamp))
	})
	if err != nil {
		Error.Println("Error on adding content to curation", err)
		return false
	}

	return true
}

func (c *DefaultCurator) OnCommentAdded(obj *Post) bool {
	return true
}

func (c *DefaultCurator) GetContent(params map[string]interface{}) []string {
	hashes := []string{}
	err := curationDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(curationBucket))
		b.ForEach(func(k, _ []byte) error {
			hashes = append(hashes, string(k))
			return nil
		})
		return nil
	})
	if err != nil {
		return []string{}
	}

	return hashes

}

func (c *DefaultCurator) FlagContent(hash string, isFlagged bool) {

}

func (c *DefaultCurator) UpvoteContent(hash string) {

}

func (c *DefaultCurator) DownvoteContent(hash string) {

}
