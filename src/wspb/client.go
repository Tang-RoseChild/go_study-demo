package main

import (
	proto "github.com/golang/protobuf/proto"
	"golang.org/x/net/websocket"
	"wspb/msgproto"

	"flag"
	"fmt"
)

var addr string

func init() {
	flag.StringVar(&addr, "addr", ":12345", "-addr=:12345")
	flag.Parse()

}

func main() {
	origin := "http://localhost/"
	url := fmt.Sprintf("ws://localhost%s/sock", addr)

	ws, err := websocket.Dial(url, "", origin)
	checkErr("Dial", err)
	defer ws.Close()
	// pbCodec := websocket.Codec{

	// 	Unmarshal: pbUnmarshal,
	// }
	pbCodec := websocket.Codec{
		Marshal:   pbMarshal,
		Unmarshal: pbUnmarshal,
	}
	// send single one
	msg := &msgproto.Msg{
		Id:      proto.Int32(101),
		Content: proto.String("testing protobuf 中文，å¡£ £ldsaj"),
		Topic:   proto.String("Testing"),
	}

	 buf, err := proto.Marshal(msg)
	checkErr("Marshal", err)
	 fmt.Println("buf is : ", buf)
	// err = websocket.Message.Send(ws, buf)
	// err = pbCodec.Send(ws, buf)
	err = pbCodec.Send(ws, msg)

	checkErr("Send", err)

	// receive one
	backMsg := &msgproto.Msg{}

	/* read a buffer,then use the buffer to Unmarshal again
	var recBuf []byte
	recBuf = make([]byte, 1<<10)

	err = websocket.Message.Receive(ws, recBuf) // here error will occur,server uses proto,but here use message.

	checkErr("Receive", err)
	err = proto.Unmarshal(recBuf, backMsg)
	checkErr("Unmarshal", err)

	// show the result
	fmt.Printf("result is : %s \n", backMsg.GetContent())
	*/

	err = pbCodec.Receive(ws, proto.Message(backMsg))
	checkErr("Receive", err)

	fmt.Printf("received : %s ", backMsg.GetContent())

}

func checkErr(key string, err error) {
	if err != nil {
		fmt.Printf("%s err occur : %s ", key, err.Error())
	}
}

/*
// Codec defined in websocket
type Codec struct {
    Marshal   func(v interface{}) (data []byte, payloadType byte, err error)
    Unmarshal func(data []byte, payloadType byte, v interface{}) (err error)
}

// Marshal in protobuf
func Marshal(pb Message) ([]byte, error)

// Unmarshal in protobuf
func Unmarshal(buf []byte, pb Message) error

*/

// pbMarshal : wrap protobuf marshal as websocket marshal
func pbMarshal(v interface{}) (data []byte, payloadType byte, err error) {
	// type assertion
	pb, ok := v.(proto.Message)

	if !ok {
		data, payloadType, err = make([]byte, 0), websocket.BinaryFrame, fmt.Errorf("type wrong, need proto.Message but getting %T ", v)
		return
	}

	data, err = proto.Marshal(pb)
	payloadType = websocket.BinaryFrame
	return
}

// pbUnmarshal : wrap protobuf Unmarshal as websocket marshal
func pbUnmarshal(data []byte, payloadType byte, v interface{}) (err error) {
	// type assertion
	pb, ok := v.(proto.Message)

	if !ok {
		data, payloadType, err = make([]byte, 0), websocket.BinaryFrame, fmt.Errorf("type wrong, need proto.Message but getting %T ", v)
		return
	}

	err = proto.Unmarshal(data, pb)
	return
}
