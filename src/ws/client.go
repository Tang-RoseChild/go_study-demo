package main

import proto "github.com/golang/protobuf/proto"

import (
	"fmt"
	"golang.org/x/net/websocket"
	"os"

	"pb"
)

func main() {
	origin := "http://localhost/"
	url := "ws://localhost:12345/sock"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		fmt.Println("err when dial : ", err.Error())
	}
	defer ws.Close()
	fmt.Println("ws : ", ws.RemoteAddr())
	msg := &pb.Helloworld{
		Id:  proto.Int32(100),
		Str: proto.String("中文测试碎å"),
	}
	buf, err := proto.Marshal(msg)
	fmt.Printf("buf is : %v \n", buf)

	CheckErr("marshal", err)
	// nfw, err := ws.NewFrameWriter(websocket.BinaryFrame)
	// CheckErr("NewFrameWriter", err)

	// n, err := nfw.Write(buf)
	// if n != len(buf) {
	// 	CheckErr("write", fmt.Errorf("not write all buf"))
	// }
	// CheckErr("write", err)

	/* use message send to send msg */
	err = websocket.Message.Send(ws, buf)
	CheckErr("Msg send ", err)
}

func CheckErr(msg string, err error) {
	if err != nil {
		fmt.Printf("%s failed, error is : %s \n", msg, err.Error())
		os.Exit(1)
	}
}
