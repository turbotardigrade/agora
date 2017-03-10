package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"gx/ipfs/QmaPHkZLbQQbvcyavn8q1GFHg6o6yeceyHFSJ3Pjf3p3TQ/go-crypto/ssh"
	"strings"

	"github.com/gtank/cryptopasta"
)

// Notes of why we moved from RSA to ECDSA
// http://crypto.stackexchange.com/questions/3216/signatures-rsa-compared-to-ecdsa

func FingerprintSHA256(key ssh.PublicKey) string {
	hash := sha256.Sum256(key.Marshal())
	b64hash := base64.StdEncoding.EncodeToString(hash[:])
	return strings.TrimRight(b64hash, "=")
}

// Verify checks if the data of IPFSObj is valid by checking the
// signature with provided PublicKey of the author
func Verify(obj *IPFSObj) (bool, error) {
	signature, err := cryptopasta.DecodeSignatureJWT(obj.Signature)
	if err != nil {
		return false, err
	}

	var dataBuf bytes.Buffer
	gob.NewEncoder(&dataBuf).Encode(obj.Data)

	pubKey, err := cryptopasta.DecodePublicKey([]byte(obj.Key))
	return cryptopasta.Verify(dataBuf.Bytes(), signature, pubKey), nil
}

// Sign creates cryptographic signature to let other nodes verify if
// given User has indeed posted this
func Sign(user *User, data interface{}) (string, error) {
	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(data)

	signature, err := cryptopasta.Sign(buf.Bytes(), user.PrivateKey)
	if err != nil {
		return "", err
	}

	return cryptopasta.EncodeSignatureJWT(signature), nil
}
