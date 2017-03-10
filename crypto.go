package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
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

	data, err := json.Marshal(obj.Data)
	if err != nil {
		return false, err
	}

	pubKey, err := cryptopasta.DecodePublicKey([]byte(obj.Key))
	if err != nil {
		return false, err
	}

	return cryptopasta.Verify(data, signature, pubKey), nil
}

// Sign creates cryptographic signature to let other nodes verify if
// given User has indeed posted this
func Sign(user *User, data map[string]interface{}) (string, error) {

	dataBuf, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	privKey, err := user.GetPrivateKey()
	if err != nil {
		return "", err
	}

	signature, err := cryptopasta.Sign(dataBuf, privKey)
	if err != nil {
		return "", err
	}

	return cryptopasta.EncodeSignatureJWT(signature), nil
}
