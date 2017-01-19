package main

import (
	"context"
	"fmt"
	"syscall"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/namesys"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
)

var ipfsFileDescNum = uint64(5120)

// Taken from github.com/ipfs/go-ipfs/blob/master/cmd/ipfs/ulimit_unix.go
func checkAndSetUlimit() error {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return fmt.Errorf("Error getting rlimit: %s", err)
	}

	var setting bool
	if rLimit.Cur < ipfsFileDescNum {
		if rLimit.Max < ipfsFileDescNum {
			fmt.Println("Error: adjusting max")
			rLimit.Max = ipfsFileDescNum
		}
		fmt.Printf("Adjusting current ulimit to %d...\n", ipfsFileDescNum)
		rLimit.Cur = ipfsFileDescNum
		setting = true
	}

	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return fmt.Errorf("Error setting ulimit: %s", err)
	}

	if setting {
		fmt.Printf("Successfully raised file descriptor limit to %d.\n", ipfsFileDescNum)
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
