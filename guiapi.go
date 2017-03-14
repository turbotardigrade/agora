package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/mitchellh/mapstructure"
)

const GUIAPITimeout = 60 * time.Second

// Command defines the general structure how the GUI sends commands to the PeerBackend
// In JSON it looks like this:
// { "command": "someCommand", "arguments": {...} }
type Command struct {
	Command   string
	Arguments map[string]interface{}
}

// GUIAPI is simply used as namespace
type GUIAPI struct{}

var gAPI = GUIAPI{}

// cmd2func maps the command with its respective handler function
var cmd2func = map[string]func(args map[string]interface{}) string{
	"getPost":             gAPI.getPost,
	"getComment":          gAPI.getComment,
	"getPosts":            gAPI.getPosts,
	"postPost":            gAPI.postPost,
	"getCommentsFromPost": gAPI.getCommentsFromPost,
	"postComment":         gAPI.postComment,
	"setPostUserData":     gAPI.setPostUserData,
	"setCommentUserData":  gAPI.setCommentUserData,

	"upvote":   gAPI.upvote,
	"downvote": gAPI.downvote,
	"flag":     gAPI.flag,
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

		resp := GUIHandle(cmd)

		// This ensures that we always return something as response
		if resp == "" {
			fmt.Println(`{"status": "ok"}`)
		} else {
			fmt.Println(resp)
		}
	}
}

func GUIHandle(cmd Command) string {
	handler, ok := cmd2func[cmd.Command]
	if !ok {
		return `{"error": "No such command."}`
	}

	ch := make(chan string, 1)
	go func() {
		ch <- handler(cmd.Arguments)
	}()

	select {
	case resp := <-ch:
		return resp
	case <-time.After(GUIAPITimeout):
		return `{"error": "Timelimit exceeded."}`
	}

}

//////////////////////////////////////////////////////////////////////
// handler functions

func (*GUIAPI) getPost(args map[string]interface{}) string {
	hash, ok := args["hash"].(string)
	if !ok {
		return `{"error": "Argument not well formatted."}`
	}

	post, err := GetPost(hash)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}

	res, _ := json.Marshal(post)
	return string(res)
}

func (*GUIAPI) getComment(args map[string]interface{}) string {
	hash, ok := args["hash"].(string)
	if !ok {
		return `{"error": "Argument not well formatted."}`
	}

	comment, err := GetComment(hash)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}

	res, _ := json.Marshal(comment)
	return string(res)
}

func (*GUIAPI) postPost(args map[string]interface{}) string {
	content, ok := args["content"].(string)
	if !ok {
		return `{"error": "Argument not well formatted."}`
	}

	title, ok := args["title"].(string)
	if !ok {
		return `{"error": "Argument not well formatted."}`
	}

	obj, err := NewPost(MyUser, title, content)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}

	return `{"hash": "` + obj.Hash + `"}`
}

func (*GUIAPI) getPosts(args map[string]interface{}) string {
	posts, err := GetContentPosts()
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}

	js, _ := json.Marshal(posts)
	return string(js)
}

func (*GUIAPI) getCommentsFromPost(args map[string]interface{}) string {
	hash, ok := args["hash"].(string)
	if !ok {
		return `{"error": "Argument not well formatted."}`
	}

	comments, err := GetComments(hash)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}

	js, _ := json.Marshal(comments)
	return string(js)
}

type postCommentArgs struct {
	Post    string
	Content string
	Parent  string
}

func (*GUIAPI) postComment(args map[string]interface{}) string {
	pArgs := postCommentArgs{}
	mapstructure.Decode(args, &pArgs)

	obj, err := NewComment(MyUser, pArgs.Post, pArgs.Parent, pArgs.Content)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}

	return `{"hash": "` + obj.Hash + `"}`
}

type setPostUserDataArgs struct {
	Hash     string
	UserData PostUserData
}

func (*GUIAPI) setPostUserData(args map[string]interface{}) string {
	pArgs := setPostUserDataArgs{}
	mapstructure.Decode(args, &pArgs)

	err := db.SetPostUserData(pArgs.Hash, pArgs.UserData)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}

	return `{"status": "success"}`
}

type setCommentUserDataArgs struct {
	Hash     string
	UserData CommentUserData
}

func (*GUIAPI) setCommentUserData(args map[string]interface{}) string {
	pArgs := setCommentUserDataArgs{}
	mapstructure.Decode(args, &pArgs)

	err := db.SetCommentUserData(pArgs.Hash, pArgs.UserData)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}

	return `{"status": "success"}`
}

func (*GUIAPI) upvote(args map[string]interface{}) string {
	hash, ok := args["hash"].(string)
	if !ok {
		return `{"error": "Argument not well formatted."}`
	}

	MyCurator.UpvoteContent(hash)

	return `{"status": "success"}`
}

func (*GUIAPI) downvote(args map[string]interface{}) string {
	hash, ok := args["hash"].(string)
	if !ok {
		return `{"error": "Argument not well formatted."}`
	}

	MyCurator.DownvoteContent(hash)

	return `{"status": "success"}`
}

func (*GUIAPI) flag(args map[string]interface{}) string {
	hash, ok := args["hash"].(string)
	if !ok {
		return `{"error": "Argument not well formatted."}`
	}

	isFlagged, ok := args["isFlagged"].(bool)
	if !ok {
		return `{"error": "Argument not well formatted."}`
	}

	MyCurator.FlagContent(hash, isFlagged)

	return `{"status": "success"}`
}
