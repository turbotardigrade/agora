package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"syscall"

	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/core"
	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/core/corenet"
	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/namesys"
	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/repo/config"
	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/repo/fsrepo"
	peer "gx/ipfs/QmZcUPvPhD1Xvk6mwijYF8AfR3mG31S1YsEfHG4khrFPRr/go-libp2p-peer"
)

const nBitsForKeypair = 2048

func Request(node *core.IpfsNode, targetPeer string, path string, body string) (string, error) {
	// Check if Node hash is valid
	target, err := peer.IDB58Decode(targetPeer)
	if err != nil {
		return "", err
	}

	// Connect to targetPeer
	stream, err := corenet.Dial(node, target, path)
	if err != nil {
		return "", err
	}

	// Exchange request and response
	stream.Write([]byte(body))

	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.String(), nil
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

const MyNodePath = "./data/clientNode"

func main() {
	if !Exists(MyNodePath) {
		err := NewNodeRepo(MyNodePath, nil)
		if err != nil {
			panic(err)
		}
	}

	node, err := NewNode(MyNodePath)
	if err != nil {
		panic(err)
	}

	res, err := Request(node, "QmdtfJBMitotUWBX5YZ6rYeaYRFu6zfXXMZP6fygEWK2iu", "/health", "")
	if err != nil {
		panic(err)
	}

	fmt.Println(res)

}

// NewNode creates a new Node from an existing node repository
func NewNode(path string) (*core.IpfsNode, error) {
	// Need to increse limit for number of filedescriptors to
	// avoid running out of those due to a lot of sockets
	// @TODO maybe move this to init
	err := checkAndSetUlimit()
	if err != nil {
		return nil, err
	}

	// Open and check node repository
	r, err := fsrepo.Open(path)
	if err != nil {
		return nil, err
	}

	// Run Node
	cfg := &core.BuildCfg{
		Repo:   r,
		Online: true,
	}

	ctx, cancel := context.WithCancel(context.Background())
	node, err := core.NewNode(ctx, cfg)
	if err != nil {
		cancel()
		return nil, err
	}

	return node, nil
}

// NewNodeRepo will create a new data and configuration folder for a
// new IPFS node at the provided location
func NewNodeRepo(repoRoot string, addr *config.Addresses) error {
	err := os.MkdirAll(repoRoot, 0755)
	if err != nil {
		return err
	}

	if fsrepo.IsInitialized(repoRoot) {
		return errors.New("Repo already exists")
	}

	conf, err := config.Init(os.Stdout, nBitsForKeypair)
	if err != nil {
		return err
	}

	if addr != nil {
		conf.Addresses = *addr
	}

	fsrepo.Init(repoRoot, conf)
	if err != nil {
		return err
	}

	return initializeIpnsKeyspace(repoRoot)
}

// Taken from github.com/ipfs/go-ipfs/blob/master/cmd/ipfs/init.go
func initializeIpnsKeyspace(repoRoot string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r, err := fsrepo.Open(repoRoot)
	if err != nil { // NB: repo is owned by the node
		return err
	}

	nd, err := core.NewNode(ctx, &core.BuildCfg{Repo: r})
	if err != nil {
		return err
	}
	defer nd.Close()

	err = nd.SetupOfflineRouting()
	if err != nil {
		return err
	}

	return namesys.InitializeKeyspace(ctx, nd.DAG, nd.Namesys, nd.Pinning, nd.PrivateKey)
}

var ipfsFileDescNum = uint64(5120)

// Taken from github.com/ipfs/go-ipfs/blob/master/cmd/ipfs/ulimit_unix.go
func checkAndSetUlimit() error {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return fmt.Errorf("Error getting rlimit: %s", err)
	}

	if rLimit.Cur < ipfsFileDescNum {
		if rLimit.Max < ipfsFileDescNum {
			log.Println("Error: adjusting max")
			rLimit.Max = ipfsFileDescNum
		}
		// Info.Println("Adjusting current ulimit to ", ipfsFileDescNum, "...")
		rLimit.Cur = ipfsFileDescNum
	}

	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return fmt.Errorf("Error setting ulimit: %s", err)
	}

	return nil
}
