package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"

	"gx/ipfs/QmRuZnMorqodado1yeTQiv1i9rmtKj29CjPSsBKM7DFXV4/go-libp2p-net"
)

// B converts string to byte, just a helper for less typing...
func B(obj string) []byte {
	return []byte(obj)
}

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

// StringInSlice check if a string is inside given string array
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func PrettyPrint(obj interface{}) {
	b, _ := json.MarshalIndent(obj, "", "  ")
	Info.Println(string(b))
}

func Now() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}
