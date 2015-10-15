package main

import (
	"chat/server"
	"fmt"
)

func main() {
	err := server.Run(":23456")
	if err != nil {
		fmt.Println("err : ", err.Error())
	}
}
