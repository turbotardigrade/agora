package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"log"
)

const (
	// MyNodePath specifies the path where the node repository
	// (containing data and configuration) of the default node
	// used by the application is stored
	MyUserConfPath = "./data/me.json"

	// nBitsForKeypair sets the strength of keypair
	nBitsForUserKeypair = 2048
)

var MyUser *User

func init() {
	var err error
	if !Exists(MyUserConfPath) {
		MyUser, err = NewUser("long")
		if err != nil {
			log.Fatal("Cannot generate new User identity", err)
		}

		userConfJSON, _ := json.Marshal(MyUser)
		err = ioutil.WriteFile(MyUserConfPath, userConfJSON, 0644)

		if err != nil {
			log.Fatal("Cannot write new UserConf to disk", err)
		}
	}

	file, err := ioutil.ReadFile(MyUserConfPath)
	if err != nil {
		log.Fatal("Cannot read user config", err)
	}

	var user User
	json.Unmarshal(file, &user)
	MyUser = &user
}

type User struct {
	PubKeyRaw  string
	PrivKeyRaw string
	Alias      string

	PublicKey  *rsa.PublicKey  `json:"-"`
	PrivateKey *rsa.PrivateKey `json:"-"`
}

func (u *User) GetPublicKey() (*rsa.PublicKey, error) {
	if u.PublicKey == nil {
		if err := u.LoadKeys(); err != nil {
			return nil, err
		}
	}

	return u.PublicKey, nil
}

func (u *User) GetPrivateKey() (*rsa.PrivateKey, error) {
	if u.PrivateKey == nil {
		if err := u.LoadKeys(); err != nil {
			return nil, err
		}
	}

	return u.PrivateKey, nil
}

func (u *User) LoadKeys() error {
	block, _ := pem.Decode([]byte(u.PrivKeyRaw))
	der := block.Bytes

	key, err := x509.ParsePKCS1PrivateKey(der)
	if err != nil {
		return err
	}

	u.PrivateKey = key
	u.PublicKey = &key.PublicKey

	return nil
}

func NewUser(alias string) (*User, error) {
	key, err := rsa.GenerateKey(rand.Reader, nBitsForUserKeypair)
	if err != nil {
		return nil, err
	}

	privKeyDer := x509.MarshalPKCS1PrivateKey(key)
	privKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privKeyDer,
	}
	privKeyPem := string(pem.EncodeToMemory(&privKeyBlock))

	publicKey := key.PublicKey
	publicKeyDer, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return nil, err
	}

	publicKeyBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   publicKeyDer,
	}
	publicKeyPem := string(pem.EncodeToMemory(&publicKeyBlock))

	return &User{
		PrivKeyRaw: privKeyPem,
		PubKeyRaw:  publicKeyPem,
		Alias:      alias,
	}, nil
}
