package main

import (
	// "fmt"
	"io"
	"log"
	"net/http"

	"html/template"

	"golang.org/x/net/websocket"
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
	var msg string
	for {
		err := websocket.Message.Receive(ws, &msg)
		CheckErr("Msg Receive", err)
		log.Printf("received : %s", msg)

		err = websocket.Message.Send(ws, "sendback "+msg)
		CheckErr("Msg SendBack", err)
		if err == io.EOF {
			log.Println("EOF, close connection")
			break
		}
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
