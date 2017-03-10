package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"gx/ipfs/QmaPHkZLbQQbvcyavn8q1GFHg6o6yeceyHFSJ3Pjf3p3TQ/go-crypto/ssh"
	"strings"
)

func FingerprintSHA256(key ssh.PublicKey) string {
	hash := sha256.Sum256(key.Marshal())
	b64hash := base64.StdEncoding.EncodeToString(hash[:])
	return strings.TrimRight(b64hash, "=")
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
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(data)
	if err != nil {
		return "", err
	}

	privKey, err := user.GetPrivateKey()
	if err != nil {
		return "", err
	}

	hashFunc := crypto.SHA1
	hasher := hashFunc.New()
	hasher.Write(buf.Bytes())

	signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, hashFunc, hasher.Sum(nil))
	if err != nil {
		return "", err
	}

	Info.Println("Signature:", base64.StdEncoding.EncodeToString(signature))

	return string(signature), nil
}
