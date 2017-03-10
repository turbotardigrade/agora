package main

import (
	"context"
	"errors"

	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/core/coreunix"

	"github.com/fatih/structs"
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
	obj := &IPFSObj{Key: user.PubKeyRaw}
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
