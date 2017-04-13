package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/gtank/cryptopasta"
)

const (
	// nBitsForKeypair sets the strength of keypair
	nBitsForUserKeypair = 2048
)

type User struct {
	PubKeyRaw  string
	PrivKeyRaw string
	Alias      string

	PublicKey  *ecdsa.PublicKey  `json:"-"`
	PrivateKey *ecdsa.PrivateKey `json:"-"`
}

func (u *User) GetPublicKey() (*ecdsa.PublicKey, error) {
	if u.PublicKey == nil {
		if err := u.LoadKeys(); err != nil {
			return nil, err
		}
	}

	return u.PublicKey, nil
}

func (u *User) GetPrivateKey() (*ecdsa.PrivateKey, error) {
	if u.PrivateKey == nil {
		if err := u.LoadKeys(); err != nil {
			return nil, err
		}
	}

	return u.PrivateKey, nil
}

func (u *User) LoadKeys() (err error) {
	u.PrivateKey, err = cryptopasta.DecodePrivateKey([]byte(u.PrivKeyRaw))
	if err != nil {
		return err
	}

	u.PublicKey, err = cryptopasta.DecodePublicKey([]byte(u.PubKeyRaw))
	return err
}

func NewUser(alias string) (*User, error) {
	key, err := cryptopasta.NewSigningKey()
	if err != nil {
		return nil, err
	}

	pubKeyRaw, err := cryptopasta.EncodePublicKey(&key.PublicKey)
	privKeyRaw, err := cryptopasta.EncodePrivateKey(key)

	u := &User{
		PubKeyRaw:  string(pubKeyRaw),
		PrivKeyRaw: string(privKeyRaw),
		PrivateKey: key,
		PublicKey:  &key.PublicKey,
		Alias:      alias,
	}

	return u, nil
}

func CreateUserIfNotExists(path, alias string) (*User, error) {
	if !Exists(path) {
		user, err := NewUser("DefaultBob")
		if err != nil {
			log.Fatal("Cannot generate new User identity", err)
		}

		userConfJSON, _ := json.MarshalIndent(user, "", "   ")
		err = ioutil.WriteFile(MyUserConfPath, userConfJSON, 0644)
		if err != nil {
			return nil, err
		}
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Cannot read user config", err)
	}

	var user User
	err = json.Unmarshal(file, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
