package main

import (
	"fmt"
	proto "github.com/golang/protobuf/proto"
	"os"

	"pb"
)

func main() {
	msg := &pb.Helloworld{
		Id:  proto.Int32(100),
		Str: proto.String("中文测试"),
	}

	f, err := os.Create(`./log.txt`)
	CheckError("create file", err)
	defer f.Close()

	buf, err := proto.Marshal(msg)
	CheckError("marshal", err)

	_, err = f.Write(buf)
	CheckError("write", err)

}

func CheckError(msg string, err error) {
	if err != nil {
		fmt.Printf("%s failed, error is : %s \n", msg, err.Error())
	}
}
