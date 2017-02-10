package main

func main() {
	OpenDb()
	defer CloseDb()

	// Starts PeerServer (non-blocking)
	StartPeerAPI(MyNode)

	// Starts communication pipeline for GUI
	StartGUIPipe()
}
