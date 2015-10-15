package main

import (
	"chat/server"

	"chat/msgproto"
	proto "github.com/golang/protobuf/proto"
	"golang.org/x/net/websocket"

	"log"

	"bufio"
	"io"
	"math/rand"
	"os"
	"time"
)

var msgCodec websocket.Codec

var (
	userInputChan   chan []byte
	serverInputChan chan *msgproto.Msg
)

func init() {
	msgCodec = websocket.Codec{
		Marshal:   server.PbMarshal,
		Unmarshal: server.PbUnmarshal,
	}

	userInputChan = make(chan []byte)
	serverInputChan = make(chan *msgproto.Msg)

}
func main() {
	url := "ws://localhost:23456/sock"
	origin := "http://localhost"
	ws, err := websocket.Dial(url, "", origin)
	checkErr("Dial", err)

	rand.Seed(int64(time.Now().Nanosecond()))
	id := rand.Int31n(1001)
	pbMsg := &msgproto.Msg{
		Id:      proto.Int32(id),
		Topic:   proto.String(""),
		Content: proto.String("hello"),
		Type:    proto.Int32(int32(server.CONNECT)),
	}
	err = msgCodec.Send(ws, pbMsg)
	checkErr("send", err)
	for {

		// fromServerMsg := &msgproto.Msg{}
		// err = msgCodec.Receive(ws, proto.Message(fromServerMsg))
		// checkErr("receive", err)

		// log.Printf("received : %s,  \n", fromServerMsg.GetContent())
	}

}

func userInput(userinput chan []byte) {

	for {
		bufReader := bufio.NewReader(os.Stdin)
		msg, err := bufReader.ReadSlice('\n')
		checkErr("readslice in handle Input", err)
		userinput <- msg
	}
}

func serverInput(ws *websockte.Conn, serverinput chan *msgproto.Msg) {
	msg := &msgproto.Msg{}
	for {
		err = msgCodec.Receive(ws, msg)
		if err == nil {
			break
		}
		serverinput <- msgCodec
	}
}

// checkErr: key string and err
func checkErr(key string, err error) {
	if err != nil {
		log.Printf("%s err occur : %s ", key, err.Error())
	}
}
