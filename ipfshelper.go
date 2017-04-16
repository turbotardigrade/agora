package main

import (
	"context"
	"fmt"
	"syscall"

	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/core"
	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/namesys"
	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/repo/fsrepo"
)

var ipfsFileDescNum = uint64(9999)

// Taken from github.com/ipfs/go-ipfs/blob/master/cmd/ipfs/ulimit_unix.go
func checkAndSetUlimit() error {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return fmt.Errorf("Error getting rlimit: %s", err)
	}

	if rLimit.Cur < ipfsFileDescNum {
		if rLimit.Max < ipfsFileDescNum {
			Error.Println("Error: adjusting max")
			rLimit.Max = ipfsFileDescNum
		}
		Info.Println("Adjusting current ulimit to ", ipfsFileDescNum, "...")
		rLimit.Cur = ipfsFileDescNum
	}

	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return fmt.Errorf("Error setting ulimit: %s", err)
	}

	return nil
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
