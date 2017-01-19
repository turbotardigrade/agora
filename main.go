package main

func main() {

	// Starts PeerServer (non-blocking)
	StartPeerAPI(MyNode)

	// @TODO API server for GUI should be running and blocking
	// here instead of the endless loop
	for {
		// Endless loop
	}
}
