package main

import (
	"context"
	"errors"

	"github.com/fatih/structs"
	"github.com/ipfs/go-ipfs/core/coreunix"
)

// IPFSObj is an abstraction to deal with objects / blobs from IPFS;
// also does signing and verification
type IPFSObj struct {
	// Hash is the hash address given by IPFS
	Hash string
	// Key is the PublicKey of the signing user
	Key string
	// Data is an opaque field to dump the payload in
	Data map[string]interface{}
	// Signature used to verify this object was indeed sent by
	// user with Key
	Signature string
}

// NewIPFSObj is a generalized helper function to created signed
// IPFSObj and add it to the IPFS network
func NewIPFSObj(node *Node, user *User, data interface{}) (*IPFSObj, error) {
	obj := &IPFSObj{Key: user.PubKey}
	obj.Data = structs.New(data).Map()

	// @TODO make cryptographic signature with given key and data
	var err error
	obj.Signature, err = Sign(user, data)
	if err != nil {
		return nil, err
	}

	// Add to IPFS Node Repository
	hash, err := coreunix.Add(node.IpfsNode, ToJSONReader(obj))
	if err != nil {
		return nil, err
	}
	obj.Hash = hash

	return obj, nil
}

func GetIPFSObj(hash string) (*IPFSObj, error) {
	ctx := context.Background() // Not sure what this should be used for
	r, err := coreunix.Cat(ctx, MyNode.IpfsNode, hash)
	if err != nil {
		return nil, err
	}

	obj := IPFSObj{}
	err = FromJSONReader(r, &obj)

	// @TODO verify
	ok, err := Verify(&obj)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("This got rigged")
	}

	obj.Hash = hash
	return &obj, nil
}

// Verify checks if the data of IPFSObj is valid by checking the
// signature with provided PublicKey of the author
func Verify(obj *IPFSObj) (bool, error) {
	// --------------------------------------------------
	// @TODO
	// --------------------------------------------------
	return true, nil
}

// Sign creates cryptographic signature to let other nodes verify if
// given User has indeed posted this
func Sign(user *User, data interface{}) (string, error) {
	// --------------------------------------------------
	// @TODO
	// --------------------------------------------------
	return "TODO", nil
}
