package main

func main() {
	OpenDb()
	defer CloseDb()

	// Starts PeerServer (non-blocking)
	StartPeerAPI(MyNode)

	/*
		Info.Println("Request Posts")
			posts, err := Client{MyNode}.GetPosts("QmbygukvwkdAsjLt2xwricf8BiMqTifyE1h95bwgAJ7ApL")
			if err != nil {
				Warning.Println(err)
			} else {
				fmt.Println(posts)
			}
	*/

	// Starts communication pipeline for GUI
	StartGUIPipe()
}
