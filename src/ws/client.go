package main

import (
	"fmt"
	"golang.org/x/net/websocket"
)

func main() {
	origin := "http://localhost/"
	url := "ws://localhost:12345/sock"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		fmt.Println("err when dial : ", err.Error())
	}
	fmt.Println("ws : ", ws.RemoteAddr())
}
