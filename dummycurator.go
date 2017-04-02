package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

// DummyCurator is a placeholder curation module. Currently it
// registers all content, but does no filtering nor ranking
type DummyCurator struct{}

var curationDB *bolt.DB

const curationDBPath = "./data/curation.db"
const curPostBucket = "exampleCuration_post"
const curFlagBucket = "exampleCuration_flags"

// Init initializes boltdb which simply keeps track of saved hashes
// and their arrivaltime
func (c *DummyCurator) Init() error {
	Info.Println("Init Curation DB")
	config := &bolt.Options{Timeout: 2 * time.Second}

	var err error
	curationDB, err = bolt.Open(curationDBPath, 0600, config)
	if err != nil {
		Error.Println("FATAL", err)
		log.Fatal(err)
	}

	// Create Buckets if they don't exists
	err = curationDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(curPostBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return curationDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(curFlagBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func (c *DummyCurator) OnPostAdded(obj *Post, isWhitelabeled bool) bool {
	err := curationDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(B(curPostBucket))
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)

		return bucket.Put(B(obj.Hash), B(timestamp))
	})
	if err != nil {
		Error.Println("Error on adding content to curation", err)
		return false
	}

	return true
}

func (c *DummyCurator) OnCommentAdded(obj *Comment, isWhitelabeled bool) bool {
	return true
}

func (c *DummyCurator) GetContent(params map[string]interface{}) []string {
	hashes := []string{}
	err := curationDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(B(curPostBucket))
		b.ForEach(func(k, _ []byte) error {

			// Check if content is flagged
			fB := tx.Bucket(B(curFlagBucket))

			isFlagged := fB.Get(k)
			if isFlagged != nil && string(isFlagged) == "true" {
				// Skip flagged elemens
				return nil
			}

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

func (c *DummyCurator) FlagContent(hash string, isFlagged bool) {
	Info.Println(hash, "flagged", isFlagged)

	err := curationDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(B(curFlagBucket))

		if isFlagged {
			return bucket.Put(B(hash), B("true"))
		} else {
			return bucket.Put(B(hash), B("false"))
		}

	})
	if err != nil {
		Warning.Println("Error on adding content to curation", err)
	}
}

func (c *DummyCurator) UpvoteContent(hash string) {
	Info.Println(hash, "upvoted")
}

func (c *DummyCurator) DownvoteContent(hash string) {
	Info.Println(hash, "donwvoted")
}

func (c *DummyCurator) Close() error {
	Info.Println("Destroy curation module")
	return nil
}
