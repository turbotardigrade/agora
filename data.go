package main

import (
	"context"
	"errors"
	"time"

	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/core/coreunix"

	"github.com/fatih/structs"
)

const IPFSTimeoutDuration = 30 * time.Second

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
	obj := &IPFSObj{Key: user.PubKeyRaw}
	obj.Data = structs.New(data).Map()

	var err error
	obj.Signature, err = Sign(user, obj.Data)
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

// RiggedError is thrown if the signature of a post does not match
// with public key identity, means someone tries to impersonate
// someone else
var RiggedError = errors.New("This object got rigged")

// GetIPFSObj retrieves data from IPFS and verifies the signature. If
// the signature does not match, it means that someone pretends to
// send content under someone else's identity, which will throw a
// RiggedError
//
// @TODO this smells like refactoring since it uses global variable
// MyNode
func GetIPFSObj(hash string) (*IPFSObj, error) {
	ctx, cancel := context.WithTimeout(context.Background(), IPFSTimeoutDuration)
	defer cancel()

	r, err := coreunix.Cat(ctx, MyNode.IpfsNode, hash)
	if err != nil {
		return nil, err
	}

	obj := &IPFSObj{}
	err = FromJSONReader(r, obj)
	if err != nil {
		return nil, err
	}

	ok, err := Verify(obj)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, RiggedError
	}

	// Set this as the hash is not stored inside IPFS blob
	obj.Hash = hash
	return obj, nil
}
