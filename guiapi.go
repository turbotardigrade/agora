package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
)

type Command struct {
	Command   string
	Arguments map[string]interface{}
}

var cmd2func = map[string]func(args map[string]interface{}){
	"getPost":             getPost,
	"postPost":            postPost,
	"getCommentsFromPost": getCommentsFromPost,
	"postComment":         postComment,
	"postContent":         postContent,
}

func StartGUIPipe() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var cmd Command
		input := scanner.Text()
		err := json.Unmarshal([]byte(input), &cmd)
		if err != nil {
			fmt.Println(`{"error": "JSON object not well formed."}`)
			continue
		}

		handler, ok := cmd2func[cmd.Command]
		if !ok {
			fmt.Println(`{"error": "No such command."}`)
		} else {
			handler(cmd.Arguments)
		}
	}
}

func getPost(args map[string]interface{}) {
	hash, ok := args["hash"].(string)
	if !ok {
		fmt.Println(`{"error": "Argument not well formatted."}`)
		return
	}

	post, err := GetPost(hash)
	if err != nil {
		fmt.Println(`{"error": "`, err, `"}`)
		return
	}

	res, _ := json.Marshal(post)
	fmt.Println(string(res))
}

func postPost(args map[string]interface{}) {
	content, ok := args["content"].(string)
	if !ok {
		fmt.Println(`{"error": "Argument not well formatted."}`)
		return
	}

	obj, err := NewPost(MyUser, content)
	if err != nil {
		fmt.Println(`{"error": "`, err, `"}`)
		return
	}

	fmt.Println(`{"hash": "` + obj.Hash + `"}`)
}

func getCommentsFromPost(args map[string]interface{}) {
	hash, ok := args["hash"].(string)
	if !ok {
		fmt.Println(`{"error": "Argument not well formatted."}`)
		return
	}

	comments, err := GetComments(hash)
	if err != nil {
		fmt.Println(`{"error": "`, err, `"}`)
		return
	}

	// @DANGER can easily break at this place
	// @TODO do concurrently

	var res []Comment
	for _, c := range comments {
		comment, err := GetComment(c)
		if err != nil {
			fmt.Println(`{"error": "`, err, `"}`)
			return
		}

		res = append(res, *comment)
	}

	js, _ := json.Marshal(res)
	fmt.Println(string(js))
}

type postCommentArgs struct {
	Post    string
	Content string
	Parent  string
}

func postComment(args map[string]interface{}) {
	pArgs := postCommentArgs{}
	mapstructure.Decode(args, &pArgs)

	obj, err := NewComment(MyUser, pArgs.Post, pArgs.Parent, pArgs.Content)
	if err != nil {
		fmt.Println(`{"error": "`, err, `"}`)
		return
	}

	fmt.Println(`{"hash": "` + obj.Hash + `"}`)
}

func postContent(args map[string]interface{}) {
	fmt.Println(`{"res": "DEPRECATED"}`)
}
