package main

import proto "github.com/golang/protobuf/proto"
import (
	"fmt"
	// "io"
	"log"
	"net/http"

	"html/template"

	"golang.org/x/net/websocket"

	"pb"
)

func main() {
	// web server
	// http.Handle("/sock", websocket.Handler(WSServer))
	http.HandleFunc("/sock", WS)
	http.HandleFunc("/", index)
	err := http.ListenAndServe(":12345", nil)
	CheckErr("ListenAndServe", err)
}

func index(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, "welcome")

	t, err := template.ParseFiles(`index.html`)
	CheckErr("Parse", err)
	err = t.Execute(w, nil)
	CheckErr("template execute", err)
}

func WSServer(ws *websocket.Conn) {
	defer ws.Close()
	log.Printf("connected : %s \n", ws.Request().RemoteAddr)

	// var msg []byte
	pbCodec := websocket.Codec{
		pbMarshal,
		pbUnmarshal,
	}
	pbMsg := &pb.Helloworld{}
	// msg = make([]byte, 1<<10)
	for {
		// err := websocket.Message.Receive(ws, &msg)
		// instead of pbCodec
		err := pbCodec.Receive(ws, proto.Message(pbMsg))

		CheckErr("Msg Receive", err)
		log.Printf("received : %v", pbMsg)
		if err != nil {
			log.Printf("error %s , close connection\n", err.Error())
			break
		}

		// //Unmarshal through protobuf
		// err = proto.Unmarshal(msg, pbMsg)
		// CheckErr("proto Unmarshal ", err)
		// log.Printf("after Unmarshal : %v data : %d,%s\n", pbMsg, pbMsg.GetId(), pbMsg.GetStr())

		// // err = websocket.Message.Send(ws, "sendback "+msg)
		// // CheckErr("Msg SendBack", err)
		// if err != nil {
		// 	log.Printf("error %s , close connection\n", err.Error())
		// 	break
		// }
	}
}
func WS(w http.ResponseWriter, r *http.Request) {
	websocket.Handler(WSServer).ServeHTTP(w, r)
	// s := websocket.Server{Handler: websocket.Handler(WSServer)}
	// s.ServeHTTP(w, r)
}

// func WS(w http.ResponseWriter, r *http.Request) {
// 	websocket.Handler(WSServer).ServeHTTP(w, r)
// 	// s := websocket.Server{Handler: websocket.Handler(WSServer)}
// 	// s.ServeHTTP(w, r)
// }

func CheckErr(t string, err error) {
	if err != nil {
		log.Printf("%s failed, err is : %s\n", t, err.Error())
	}
}

// func decorate new marshal and unmarshal through  proto.Marshal and proto.Unmarshal
/*
// websocket codec
type Codec struct {
    Marshal   func(v interface{}) (data []byte, payloadType byte, err error)
    Unmarshal func(data []byte, payloadType byte, v interface{}) (err error)
}

// protobuf marshal
func Marshal(pb Message) ([]byte, error)

// protobuf unmarshal
func Unmarshal(buf []byte, pb Message) error
*/
func pbMarshal(v interface{}) (data []byte, payloadType byte, err error) {
	pb, ok := v.(proto.Message)
	if !ok {
		return make([]byte, 0), websocket.BinaryFrame, fmt.Errorf("need proto.Message type, in fact is %T ", v)
	}

	data, err = proto.Marshal(pb)
	payloadType = websocket.BinaryFrame
	return
}

func pbUnmarshal(data []byte, payloadType byte, v interface{}) (err error) {
	// proto.Unmarshal(data, pb)
	pb, ok := v.(proto.Message)
	if !ok {
		err = fmt.Errorf("need proto.Message type, in fact is %T ", v)
		return
	}

	err = proto.Unmarshal(data, pb)
	v = pb
	fmt.Println("data", data, "in pbUnmarshal : ", v)
	return
}
