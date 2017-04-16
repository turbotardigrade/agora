package main

import (
	"os"
	"path"
	"strings"
)

var (
	// ExecutionPath will be set on runtime to determine the location of the binary
	ExecutionPath string

	// LogDirPath is the path of the log directory
	LogDirPath string

	// MyNodePath specifies the path where the node repository
	// (containing data and configuration) of the default node
	// used by the application is stored
	MyNodePath string

	// MyNodePath specifies the path where the node repository
	// (containing data and configuration) of the default node
	// used by the application is stored
	MyUserConfPath string

	// MyNode provides a global Node instance of user's own node
	MyNode *Node

	// MyUser provides a global instance of user running this agora node
	MyUser *User
)

func InitPaths() error {
	ex, err := os.Executable()
	if err != nil {
		return err
	}

	ExecutionPath = path.Dir(ex)

	LogDirPath = ExecutionPath + "/data"
	MyNodePath = ExecutionPath + "/data/MyNode"
	MyUserConfPath = ExecutionPath + "/data/me.json"

	return nil
}

func init() {
	err := InitPaths()
	if err != nil {
		panic(err)
	}

	// Initializes flags and make flags available globally via
	// opts variable
	InitFlags()

	InitLogger()

	// Need to increse limit for number of filedescriptors to
	// avoid running out of those due to a lot of sockets
	err = checkAndSetUlimit()
	if err != nil {
		panic(err)
	}

	// Set Curator module
	if strings.ToLower(opts.Curator) == "none" {
		Info.Println("Using DummyCurator")
		MyCurator = &DummyCurator{}
	} else {
		Info.Println("Using MLCurator")
		MyCurator = &MLCurator{}
	}
}
