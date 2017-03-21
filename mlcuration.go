package main

import (
	"bufio"
	"encoding/json"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path"
)

type CCommand struct {
	Id        int                    `json:"id"`
	Command   string                 `json:"command"`
	Arguments map[string]interface{} `json:"arguments"`
}

type Result struct {
	Id     int         `json:"id"`
	Result interface{} `json:"result"`
	Error  string      `json:"error"`
}

type MLCurator struct {
	cmd     *exec.Cmd
	in      io.WriteCloser
	scanner *bufio.Scanner
}

func (c *MLCurator) Init() error {
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	exPath := path.Dir(ex)
	c.cmd = exec.Command(exPath + "/dist/main")
	c.in, err = c.cmd.StdinPipe()
	if err != nil {
		return err
	}

	out, err := c.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	c.scanner = bufio.NewScanner(out)
	err = c.cmd.Start()
	if err != nil {
		Error.Println(err)
	}
	return err
}

func (c *MLCurator) OnPostAdded(obj *Post, isWhitelabeled bool) bool {
	res, _ := c.sendCommand(CCommand{
		rand.Int(), "on_post_added", map[string]interface{}{
			"obj":            obj,
			"isWhitelabeled": isWhitelabeled,
		},
	})
	if res.Error != "" {
		Error.Println(res.Error)
		return false
	}
	r, ok := res.Result.(bool)
	if !ok {
		Error.Println("Result is not a bool")
		return false
	}
	return r
}

func (c *MLCurator) OnCommentAdded(obj *Comment, isWhitelabeled bool) bool {
	res, _ := c.sendCommand(CCommand{
		rand.Int(), "on_comment_added", map[string]interface{}{
			"obj":            obj,
			"isWhitelabeled": isWhitelabeled,
		},
	})
	if res.Error != "" {
		Error.Println(res.Error)
		return false
	}
	r, ok := res.Result.(bool)
	if !ok {
		Error.Println("Result is not a bool")
		return false
	}
	return r
}

func (c *MLCurator) GetContent(params map[string]interface{}) []string {
	res, _ := c.sendCommand(CCommand{
		rand.Int(), "get_content", map[string]interface{}{
			"params": params,
		},
	})
	if res.Error == "" {
		r := res.Result.([]interface{})
		if len(r) > 0 {
			retval := make([]string, len(r))
			for i, v := range r {
				retval[i] = v.(string)
			}
			return retval
		}
		return nil
	}
	return nil
}

func (c *MLCurator) FlagContent(hash string, isFlagged bool) {
	c.sendCommand(CCommand{
		rand.Int(), "flag_content", map[string]interface{}{
			"hash":      hash,
			"isFlagged": isFlagged,
		},
	})
}

func (c *MLCurator) UpvoteContent(hash string) {
	c.sendCommand(CCommand{
		rand.Int(), "upvote_content", map[string]interface{}{
			"hash": hash,
		},
	})
}

func (c *MLCurator) DownvoteContent(hash string) {
	c.sendCommand(CCommand{
		rand.Int(), "downvote_content", map[string]interface{}{
			"hash": hash,
		},
	})
}

func (c *MLCurator) Close() error {
	_, err := c.sendCommand(CCommand{
		rand.Int(), "quit", nil,
	})
	return err
}

// send command to pipe and read corresponding response
func (c *MLCurator) sendCommand(command CCommand) (*Result, error) {
	obj, err := json.Marshal(command)
	var res Result
	if err != nil {
		res = Result{
			command.Id, nil, "Marshal: " + err.Error(),
		}
		return &res, err
	}
	_, err = io.WriteString(c.in, string(obj)+"\n")
	if err != nil {
		res = Result{
			command.Id, nil, "I/O: " + err.Error(),
		}
		return &res, err
	}
	for c.scanner.Scan() {
		if len(c.scanner.Text()) > 0 {
			break
		}
	}
	err = json.Unmarshal([]byte(c.scanner.Text()), &res)
	if err != nil {
		res = Result{
			command.Id, nil, "I/O: " + err.Error(),
		}
		return &res, err
	}
	return &res, nil
}
