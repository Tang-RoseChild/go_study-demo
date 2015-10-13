package main

import proto "github.com/golang/protobuf/proto"

import (
	"fmt"
	"io"
	"os"
	"pb"
)

func main() {
	f, err := os.Open(`./log.txt`)
	CheckError("Open file", err)
	defer f.Close()

	fi, err := f.Stat()
	CheckError("file stat", err)

	buf := make([]byte, fi.Size())

	n, err := io.ReadFull(f, buf)
	if int64(n) != fi.Size() || err != nil {
		CheckError("ReadFull ", err)
	}
	msg := &pb.Helloworld{}
	err = proto.Unmarshal(buf, msg)

	CheckError("Unmarshal", err)

	fmt.Println("read from file : ", msg.GetId(), msg.GetStr())

}

func CheckError(msg string, err error) {
	if err != nil {
		fmt.Printf("%s failed, error is : %s \n", msg, err.Error())
	}
}
