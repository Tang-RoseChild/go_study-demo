package main

import (
	proto "github.com/golang/protobuf/proto"
	"golang.org/x/net/websocket"
	"wspb/msgproto"

	"fmt"
	"log"
	"net/http"

	"flag"
)

var addr string

func init() {
	flag.StringVar(&addr, "addr", ":12345", "-addr=:12345")
	flag.Parse()

}

func main() {
	/* add handler */
	// welcome and chat web
	http.HandleFunc("/", index)

	// websocket
	http.HandleFunc("/sock", wsserver)

	// start server
	err := http.ListenAndServe(addr, nil)
	checkErr("ListenAndServe", err)
}

// checkErr: key string and err
func checkErr(key string, err error) {
	if err != nil {
		fmt.Printf("%s err occur : %s ", key, err.Error())
	}
}

// index:show a chat web
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "welcome, because i don't know how to use js to marshal protobuf,there is no chat room")
}

// wsserver:handle websocket connect，just wrap the websocket's conn
func wsserver(w http.ResponseWriter, r *http.Request) {
	fmt.Println("connected : ", r.RemoteAddr)
	websocket.Handler(echoServer).ServeHTTP(w, r)
}

func echoServer(ws *websocket.Conn) {
	defer ws.Close()

	msg := &msgproto.Msg{}

	// 'cause should use new marshal and unmarshal to get the message, so need to create new codec
	pbCodec := websocket.Codec{
		Marshal:   pbMarshal,
		Unmarshal: pbUnmarshal,
	}

	var (
		err     error
		backMsg *msgproto.Msg
	)

	// loop to get and send msg
	for {
		err = pbCodec.Receive(ws, proto.Message(msg))
		checkErr("Receive", err)
		if err != nil {
			break
		}
		backMsg = &msgproto.Msg{
			Id:      proto.Int32(msg.GetId()),
			Topic:   proto.String(msg.GetTopic()),
			Content: proto.String("back msg " + msg.GetContent()),
		}
		err = pbCodec.Send(ws, backMsg)
		checkErr("send ", err)
		if err != nil {
			break
		}
		log.Printf("received : %s \t send : %s", msg.GetContent(), backMsg.GetContent())
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
