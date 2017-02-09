package main

func main() {

	// Starts PeerServer (non-blocking)
	StartPeerAPI(MyNode)

	// Starts communication pipeline for GUI
	StartGUIPipe()
}
