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

var cmd2func = map[string]func(args map[string]interface{}){
	"postPost":    postPost,
	"postComment": postComment,
	"postContent": postContent,
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

		handler, ok := cmd2func[cmd.Command]
		if !ok {
			fmt.Println("{\"error\": \"No such command.\"}")
		} else {
			handler(cmd.Arguments)
		}
	}
}

func postPost(args map[string]interface{}) {
	fmt.Println("{\"res\": \"ok Post\"}")
}

func postComment(args map[string]interface{}) {
	fmt.Println("{\"res\": \"ok Comment\"}")
}

func postContent(args map[string]interface{}) {
	fmt.Println("{\"res\": \"ok Test\"}")
}
