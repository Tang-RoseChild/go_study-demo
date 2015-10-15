package main

import (
	"chat/server"

	"chat/msgproto"
	// proto "github.com/golang/protobuf/proto"
	"golang.org/x/net/websocket"

	"log"

	"bufio"
	"fmt"
	// "io"
	"math/rand"
	"os"
	// "strconv"
	"strings"
	"time"
)

var msgCodec websocket.Codec

var (
	userInputChan   chan []byte
	serverInputChan chan *msgproto.Msg
	id              int32
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
	id = rand.Int31n(1001)
	// pbMsg := &msgproto.Msg{
	// 	Id:      proto.Int32(id),
	// 	Topic:   proto.String(""),
	// 	Content: proto.String("hello"),
	// 	Type:    proto.Int32(int32(server.CONNECT)),
	// }
	// err = msgCodec.Send(ws, pbMsg)
	// checkErr("send", err)
	go userInput(userInputChan)
	go serverInput(ws, serverInputChan)
	for {

		// fromServerMsg := &msgproto.Msg{}
		// err = msgCodec.Receive(ws, proto.Message(fromServerMsg))
		// checkErr("receive", err)

		// log.Printf("received : %s,  \n", fromServerMsg.GetContent())
		select {
		case usercontent := <-userInputChan:
			// handle use input msg
			handleUserInputMsg(ws, usercontent)
		case servercontent := <-serverInputChan:
			handlerServerInputMsg(servercontent)
		default:
			time.Sleep(2 * time.Second)

		}
	}

}

func userInput(userinputchan chan []byte) {

	for {
		bufReader := bufio.NewReader(os.Stdin)
		msg, err := bufReader.ReadSlice('\n')
		checkErr("readslice in handle Input", err)
		userinputchan <- msg
	}
}

func serverInput(ws *websocket.Conn, serverinputchan chan *msgproto.Msg) {
	msg := &msgproto.Msg{}
	for {
		err := msgCodec.Receive(ws, msg)
		checkErr("Receive", err)
		if err != nil {
			break
		}
		serverinputchan <- msg
		// fmt.Printf("recevied msg from server : %s \n", msg.GetContent())
	}
}

// checkErr: key string and err
func checkErr(key string, err error) {
	if err != nil {
		log.Printf("%s err occur : %s ", key, err.Error())
	}
}

func handleUserInputMsg(ws *websocket.Conn, msg []byte) {
	// login
	if string(msg[:len("login")]) == "login" {

		pbMsg := server.PbMsgFactory(id, "", "", int32(server.LOGIN))
		err := msgCodec.Send(ws, pbMsg)
		checkErr("send in handleUserInputMsg", err)

	}
	// if connect XXX
	if string(msg[:len("connect")]) == "connect" {
		topic := strings.TrimSpace(string(msg[len("connect"):]))
		pbMsg := server.PbMsgFactory(id, topic, "", int32(server.CONNECT))
		err := msgCodec.Send(ws, pbMsg)
		checkErr("send in handleUserInputMsg", err)
	}
	// test connect first
}

func handlerServerInputMsg(msg *msgproto.Msg) {

	switch server.MsgType(msg.GetType()) {
	case server.LOGIN:
		fmt.Print(msg.GetContent())
	case server.CONNECT:
		fmt.Println(msg.GetContent())
	default:
		log.Println("error, msg type wrong in handlerServerInputMsg ")
	}
}
