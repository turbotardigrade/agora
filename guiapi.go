package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
)

// Command defines the general structure how the GUI sends commands to the PeerBackend
// In JSON it looks like this:
// { "command": "someCommand", "arguments": {...} }
type Command struct {
	Command   string
	Arguments map[string]interface{}
}

// cmd2func maps the command with its respective handler function
var cmd2func = map[string]func(args map[string]interface{}){
	"getPost":             getPost,
	"getComment":          getComment,
	"getPosts":            getPosts,
	"postPost":            postPost,
	"getCommentsFromPost": getCommentsFromPost,
	"postComment":         postComment,
	"postContent":         postContent,
	"setPostUserData":     setPostUserData,
	"setCommentUserData":  setCommentUserData,
}

// STartGUIPipe is a blocking loop providing a pipe for the GUI to
// interact with the PeerBackend. Use EOF (CTRL+D) to gracefully close
// the pipe
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

//////////////////////////////////////////////////////////////////////
// handler functions

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

func getComment(args map[string]interface{}) {
	hash, ok := args["hash"].(string)
	if !ok {
		fmt.Println(`{"error": "Argument not well formatted."}`)
		return
	}

	comment, err := GetComment(hash)
	if err != nil {
		fmt.Println(`{"error": "`, err, `"}`)
		return
	}

	res, _ := json.Marshal(comment)
	fmt.Println(string(res))
}

func postPost(args map[string]interface{}) {
	content, ok := args["content"].(string)
	if !ok {
		fmt.Println(`{"error": "Argument not well formatted."}`)
		return
	}

	title, ok := args["title"].(string)
	if !ok {
		fmt.Println(`{"error": "Argument not well formatted."}`)
		return
	}

	obj, err := NewPost(MyUser, title, content)
	if err != nil {
		fmt.Println(`{"error": "`, err, `"}`)
		return
	}

	fmt.Println(`{"hash": "` + obj.Hash + `"}`)
}

func getPosts(args map[string]interface{}) {
	posts, err := GetAllPosts()
	if err != nil {
		fmt.Println(`{"error": "`, err, `"}`)
		return
	}

	js, _ := json.Marshal(posts)
	fmt.Println(string(js))
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

	js, _ := json.Marshal(comments)
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

type setPostUserDataArgs struct {
	Hash     string
	UserData PostUserData
}

func setPostUserData(args map[string]interface{}) {
	pArgs := setPostUserDataArgs{}
	mapstructure.Decode(args, &pArgs)

	err := SetPostUserData(pArgs.Hash, pArgs.UserData)
	if err != nil {
		fmt.Println(`{"error": "`, err, `"}`)
		return
	}

	fmt.Println(`{"status": "success"}`)
}

type setCommentUserDataArgs struct {
	Hash     string
	UserData CommentUserData
}

func setCommentUserData(args map[string]interface{}) {
	pArgs := setCommentUserDataArgs{}
	mapstructure.Decode(args, &pArgs)

	err := SetCommentUserData(pArgs.Hash, pArgs.UserData)
	if err != nil {
		fmt.Println(`{"error": "`, err, `"}`)
		return
	}

	fmt.Println(`{"status": "success"}`)
}

func postContent(args map[string]interface{}) {
	fmt.Println(`{"res": "DEPRECATED"}`)
}
