package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type Command struct {
	Command   string
	Arguments map[string]interface{}
}

func StartGUIPipe() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var cmd Command
		input := scanner.Text()
		err := json.Unmarshal([]byte(input), &cmd)
		if err != nil {
			fmt.Println("{\"error\": \"JSON object not well formed.\"}")
			continue
		}

		switch cmd.Command {
		case "postContent":
			postContent(cmd.Arguments)
		default:
			fmt.Println("{\"error\": \"No such command.\"}")
		}

	}
}

func postContent(args map[string]interface{}) {
	fmt.Println("{\"res\": \"ok la\"}")
}
