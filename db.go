package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

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
	spamBucket         = "spamcount"
)

// bucketNames is used to iterate through all buckets e.g. on
// initialization
var bucketNames = []string{
	postCommentsBucket,
	postHostersBucket,
	postBucket,
	blacklistBucket,
	knownNodesBucket,
	spamBucket,
}

//////////////////////////////////////////////////////////////////////
/// Open and Close

// OpenDB opens bolt database and provides a db Model instance
// globally
func OpenDB(path string) (*Model, error) {
	Info.Println("Init DB")

	config := &bolt.Options{Timeout: 2 * time.Second}
	dbInstance, err := bolt.Open(path, 0600, config)
	if err != nil {
		// No point of running without DB, just kill the
		// application
		Error.Println("FATAL", err)
		log.Fatal(err)
	}

	db := &Model{dbInstance}

	// Create Bucket if they don't exists
	err = dbInstance.Update(func(tx *bolt.Tx) error {
		for _, bucketName := range bucketNames {
			_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	err = db.AddPeer(BootstrapPeerID)
	if err != nil {
		return nil, err
	}

	return db, nil
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

func (m *Model) GetSomePeers() ([]string, error) {
	peers, err := BoltGetKeys(m.DB, knownNodesBucket)
	if err != nil {
		return nil, err
	}

	somePeers, err := RandomStringsFromArray(peers, 1)
	if err != nil {
		return nil, err
	}

	return somePeers, nil
}

func (m *Model) AddBlacklist(identity string) error {
	err := m.RemovePeer(identity)
	if err != nil {
		return err
	}

	return BoltSet(m.DB, blacklistBucket, identity, true)
}

func (m *Model) RemoveBlacklist(identity string) error {
	return BoltDelete(m.DB, blacklistBucket, identity)
}

func (m *Model) IsBlacklisted(identity string) (bool, error) {
	var isBlacklist bool
	err := BoltGet(m.DB, blacklistBucket, identity, &isBlacklist)
	return isBlacklist, err
}

func (m *Model) AddPeer(identity string) error {
	isBlacklist, err := m.IsBlacklisted(identity)
	if err != nil {
		return err
	}

	if isBlacklist {
		Info.Println("AddPeer: Skip blacklisted identity")
		return nil
	}

	return BoltSetIfNil(m.DB, knownNodesBucket, identity, time.Now().UnixNano())
}

func (m *Model) RemovePeer(identity string) error {
	return BoltDelete(m.DB, knownNodesBucket, identity)
}

func (m *Model) RemoveAllPeers() error {
	return m.DB.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(B(knownNodesBucket))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucket(B(knownNodesBucket))
		return err
	})
}

func (m *Model) GetNumberOfPostsReceivedFromPeer(identity string) (int, error) {
	counter := 0

	err := m.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(postHostersBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			list := []string{}
			err := json.Unmarshal(v, &list)
			if err != nil {
				return err
			}

			if StringInSlice(identity, list) {
				counter += 1
			}
		}

		return nil
	})

	return counter, err
}

// TrackSpam associate hash of spam content with peer and returns the
// current spam count
func (m *Model) TrackSpam(identity, contentHash string) (int, error) {
	spamset := make(map[string]struct{})
	err := BoltGet(m.DB, spamBucket, identity, &spamset)
	if err != nil {
		return 0, err
	}

	spamset[contentHash] = struct{}{}
	return len(spamset), BoltSet(m.DB, spamBucket, identity, spamset)
}

func (m *Model) GetSpamCounts() (map[string]int, error) {
	counts := make(map[string]int)

	err := m.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(spamBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			list := []string{}
			err := json.Unmarshal(v, &list)
			if err != nil {
				return err
			}

			counts[string(k)] = len(list)
		}

		return nil
	})

	return counts, err
}

func (m *Model) GetBlacklist() (map[string]int, error) {
	counts := make(map[string]int)

	err := m.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blacklistBucket))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			identity := string(k)

			spamset := make(map[string]struct{})
			err := BoltGet(m.DB, spamBucket, identity, &spamset)
			if err != nil {
				return err
			}

			counts[string(k)] = len(spamset)
			Info.Println("Spamset count of", identity, "is", len(spamset))
		}

		return nil
	})

	return counts, err
}

func (m *Model) GetLeastSpamBlacklisted() (string, error) {
	counts, err := m.GetBlacklist()
	if err != nil {
		return "", err
	}

	if len(counts) == 0 {
		return "", errors.New("Blacklist is empty")
	}

	min := float32(9999999.9)
	peer := ""

	for k, v := range counts {
		postCounts, err := m.GetNumberOfPostsReceivedFromPeer(string(k))
		if err != nil {
			Error.Println("Failed to get number of posts received from peer", err)
			continue
		}

		spamRatio := float32(v) / float32(postCounts)

		if spamRatio < min {
			min = spamRatio
			peer = string(k)
		}
	}

	return peer, nil
}
