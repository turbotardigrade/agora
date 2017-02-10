package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"time"

	"gx/ipfs/QmQx1dHDDYENugYgqA22BaBrRfuv1coSsuPiM7rYh1wwGH/go-libp2p-net"

	"github.com/boltdb/bolt"
)

var db *bolt.DB
var dbOpened bool

// ReadJSON reads stream content in JSON format into provided the
// struct variable. Important: Struct variable must be provided as a
// pointer
func ReadJSON(stream net.Stream, ptr interface{}) {
	json.NewDecoder(stream).Decode(ptr)
}

// WriteJSON writes provided struct as JSON into stream.
func WriteJSON(stream net.Stream, obj interface{}) {
	res, _ := json.Marshal(&obj)
	stream.Write(res)
}

// ToJSONReader convert a struct to io.Reader
func ToJSONReader(obj interface{}) io.Reader {
	byteData, _ := json.Marshal(obj)
	return bytes.NewReader(byteData)
}

// FromJSONReader convert a struct to io.Reader
func FromJSONReader(r io.Reader, ptr interface{}) error {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, ptr)
	if err != nil {
		return err
	}

	return nil
}

// Exists check if path exists
func Exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// CheckWriteable checks if directory is writable by the
// application. It tests this by creating a temporary file on that
// directory.
// Taken from github.com/ipfs/go-ipfs/blob/master/cmd/ipfs/init.go
func CheckWriteable(dir string) error {
	_, err := os.Stat(dir)
	if err == nil {
		// dir exists, make sure we can write to it
		testfile := path.Join(dir, "test")
		fi, err := os.Create(testfile)
		if err != nil {
			if os.IsPermission(err) {
				return fmt.Errorf("%s is not writeable by the current user", dir)
			}
			return fmt.Errorf("unexpected error while checking writeablility of repo root: %s", err)
		}
		fi.Close()
		return os.Remove(testfile)
	}

	if os.IsNotExist(err) {
		// dir doesnt exist, check that we can create it
		return os.Mkdir(dir, 0775)
	}

	if os.IsPermission(err) {
		return fmt.Errorf("cannot write to %s, incorrect permissions", err)
	}

	return err
}

// OpenDb opens bolt database
func OpenDb() error {
	fmt.Println("Init DB ----------------------------------------------------")
	var err error
	_, filename, _, _ := runtime.Caller(0) // get full path of this file
	dbfile := path.Join(path.Dir(filename), "data/data.db")
	config := &bolt.Options{Timeout: 1 * time.Second}
	db, err = bolt.Open(dbfile, 0644, config)
	if err != nil {
		log.Fatal(err)
	}
	dbOpened = true

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("knownNodes"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return nil
}

// CloseDb closes bolt database
func CloseDb() {
	dbOpened = false
	db.Close()
}
