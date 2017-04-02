package main

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

// BoltGetKeys returns an array of all the keys from a given bucket
func BoltGetKeys(db *bolt.DB, bucketName string) ([]string, error) {
	var keys []string
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		b.ForEach(func(k, _ []byte) error {
			keys = append(keys, string(k))
			return nil
		})
		return nil
	})

	return keys, err
}

// BoltGetList will get a specific list of string (as array) from a
// boltdb bucket
func BoltGetList(db *bolt.DB, bucketName, key string) ([]string, error) {
	var list []string
	err := BoltGet(db, bucketName, key, &list)
	return list, err
}

func BoltGet(db *bolt.DB, bucketName, key string, ptr interface{}) error {
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

func BoltSet(db *bolt.DB, bucketName, key string, obj interface{}) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(B(bucketName))
		data, _ := json.Marshal(obj)
		return bucket.Put(B(key), data)
	})
}

func BoltSetIfNil(db *bolt.DB, bucketName, key string, obj interface{}) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(B(bucketName))
		orig := bucket.Get(B(key))

		if orig == nil {
			data, _ := json.Marshal(obj)
			return bucket.Put(B(key), data)
		}

		return nil
	})
}

func BoltDelete(db *bolt.DB, bucketName, key string) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(B(bucketName))
		return bucket.Delete(B(key))
	})
}

// BoltAppendList appends a string to a list
func BoltAppendList(db *bolt.DB, bucketName, key, elem string) error {
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
