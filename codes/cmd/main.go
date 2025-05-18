package main

import (
	"fmt"

	"github.com/XRS0/ToTalkB/codes"
)

func main() {
	codes.Generate("hello", "./hello.jpeg")

	if in, err := codes.ScanQRCode("./hello.jpeg"); err == nil {
		fmt.Println(in)
	} else {
		fmt.Println(err)
	}
}
