package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"gx/ipfs/QmRuZnMorqodado1yeTQiv1i9rmtKj29CjPSsBKM7DFXV4/go-libp2p-net"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
	res, err := json.Marshal(&obj)
	if err != nil {
		Error.Println(err)
	}

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

// RandomStringsFromArray picks n elements form an array. Caution: This
// function changes the order of the input array. Returns all if n is
// larger than array size
func RandomStringsFromArray(list []string, n int) ([]string, error) {
	if n >= len(list) {
		return list, nil
	}

	for i := 0; i < n; i++ {
		randpos := rand.Intn(len(list))
		list[i], list[randpos] = list[randpos], list[i]
	}

	return list[0:n], nil
}

// CreateFileIfNotExists creates file iff file doesn't exists already
func CreateFileIfNotExists(path string) error {
	// detect if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
		return nil
	}

	return err
}

// RemoveContents deletes a folder at given dir path and its content
func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
