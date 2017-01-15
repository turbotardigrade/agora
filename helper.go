package main

import (
	"encoding/json"
	"gx/ipfs/QmdXimY9QHaasZmw6hWojWnCJvfgxETjZQfg9g6ZrA9wMX/go-libp2p-net"
)

func ReadJSON(stream net.Stream, ptr interface{}) {
	json.NewDecoder(stream).Decode(ptr)
}

func WriteJSON(stream net.Stream, obj interface{}) {
	res, _ := json.Marshal(&obj)
	stream.Write(res)
}
